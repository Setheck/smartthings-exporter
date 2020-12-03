package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/setheck/smartthings-exporter/smartthings"
)

type Collector struct {
	client *smartthings.Client
}

func NewCollector(client *smartthings.Client) *Collector {
	return &Collector{client: client}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

func (collector *Collector) Collect(metrics chan<- prometheus.Metric) {
	if devices, err := collector.client.ListDevices(); err != nil {
		log.Println(err)
	} else {
		for _, device := range devices {
			// Metrics
			// 1. device -  deviceId, label
			//      smartthings_device{deviceId="",deviceLabel="", name=""} 1
			// 2. device info - deviceId, manufacturer name
			//      smartthings_device_info{deviceId="",manufacturerName="",deviceManufacturerCode="",deviceTypeId="",deviceTypeId="",deviceNetworkType=""} 1
			// 3. each capability,
			//   	smartthings_device_light{deviceId="",component="switch",timestamp=""} 0  # 0 for "off"
			//		smartthings_device_temperature{deviceId="",capability="temperatureMeasurement",unit="F",timestamp=""} 68.44
			//		smartthings_device_motion{deviceId="",capability="motionSensor",timestamp=""} 0

			if m, err := prometheus.NewConstMetric(
				prometheus.NewDesc("smartthings_device", "", []string{"deviceId", "deviceLabel", "name"}, nil),
				prometheus.GaugeValue,
				1,
				device.DeviceID, device.Label, device.Name); err == nil {

				metrics <- m
			}

			if m, err := prometheus.NewConstMetric(
				prometheus.NewDesc("smartthings_device_info", "", []string{"deviceId", "manufacturerName", "deviceManufacturerCode", "deviceTypeId", "deviceNetworkType"}, nil),
				prometheus.GaugeValue,
				1,
				device.DeviceID, device.ManufacturerName, device.DeviceManufacturerCode, device.DeviceTypeID, device.DeviceNetworkType); err == nil {

				metrics <- m
			}

			for _, component := range device.Components {
				if componentStatus, err := collector.client.GetDeviceComponentStatus(device.DeviceID, component.ID); err == nil {
					for componentId, attributes := range componentStatus {
						for attributeId, properties := range attributes {
							labels := []string{"deviceId"}
							values := []string{device.DeviceID}
							labels = append(labels, "componentId")
							values = append(values, componentId)
							//labels = append(labels, "attributeId")
							//values = append(values, attributeId)

							var extras map[string]string
							metricValue := float64(0)
							for name, value := range properties {
								switch name {
								case "value":
									extras, metricValue = parseValue(attributeId, value)
									for k, v := range extras {
										labels = append(labels, k)
										values = append(values, v)
									}
								case "data":
									for k, v := range value.(map[string]interface{}) {
										labels = append(labels, k)
										values = append(values, fmt.Sprint(v))
									}
								case "timestamp":
									if tm, err := time.Parse(time.RFC3339Nano, fmt.Sprint(value)); err == nil {
										tsMillis := tm.UnixNano() / int64(time.Millisecond)
										labels = append(labels, name)
										values = append(values, fmt.Sprint(tsMillis))
									}
								default:
									labels = append(labels, name)
									values = append(values, fmt.Sprint(value))
								}
							}

							if m, err := prometheus.NewConstMetric(
								prometheus.NewDesc("smartthings_attribute_"+attributeId, "", labels, nil),
								prometheus.GaugeValue,
								metricValue,
								values...); err == nil {

								metrics <- m
							}
						}
					}
				}
			}
		}
	}
}

func parseValue(attributeId string, value interface{}) (map[string]string, float64) {
	extras := make(map[string]string)
	resultValue := float64(0)
	switch attributeId {
	case "switch":
		if str, ok := value.(string); ok {
			if str == "on" {
				resultValue = 1
			}
		}
	case "lock":
		if str, ok := value.(string); ok {
			if str == "locked" {
				resultValue = 1
			}
			extras["state"] = str
		}
	case "motion":
		if str, ok := value.(string); ok {
			if str == "active" {
				resultValue = 1
			}
			extras["state"] = str
		}
	case "indicatorStatus":
		if str, ok := value.(string); ok {
			if str == "active" {
				resultValue = 1
			}
			extras["status"] = str
		}
	default:
		switch value.(type) {
		case float64:
			resultValue = value.(float64)
		case string:
			var err error
			strValue := value.(string)
			resultValue, err = strconv.ParseFloat(strValue, 64)
			if err != nil {
				extras["value"] = strValue
			}
		}
	}

	return extras, resultValue
}
