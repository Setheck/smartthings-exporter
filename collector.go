package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/setheck/smartthings-exporter/smartthings"
)

type Collector struct {
	client smartthings.Client
}

func NewCollector(client smartthings.Client) *Collector {
	return &Collector{client: client}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

func (collector *Collector) Collect(metrics chan<- prometheus.Metric) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devices, err := collector.client.ListDevices(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for _, device := range devices {
		if m, err := prometheus.NewConstMetric(
			prometheus.NewDesc("smartthings_device",
				"a registered device",
				[]string{"deviceId", "deviceLabel", "name"}, nil),
			prometheus.GaugeValue,
			1,
			device.DeviceID, device.Label, device.Name); err == nil {

			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			prometheus.NewDesc("smartthings_device_info",
				"information about the device",
				[]string{"deviceId", "manufacturerName", "deviceManufacturerCode", "deviceTypeId", "deviceNetworkType"}, nil),
			prometheus.GaugeValue,
			1,
			device.DeviceID, device.ManufacturerName, device.DeviceManufacturerCode, device.DeviceTypeID, device.DeviceNetworkType); err == nil {

			metrics <- m
		}

		for _, component := range device.Components {
			componentStatus, err := collector.client.GetDeviceComponentStatus(ctx, device.DeviceID, component.ID)
			if err != nil {
				log.Println(err)
				return
			}

			componentMetrics(device.DeviceID, componentStatus, metrics)
		}
	}
}

func componentMetrics(deviceId string, componentStatus smartthings.ComponentStatus, metrics chan<- prometheus.Metric) {
	for componentId, attributes := range componentStatus {
		for attributeId, properties := range attributes {
			labels := []string{"deviceId", "componentId"}
			values := []string{deviceId, componentId}

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
				default:
					labels = append(labels, name)
					values = append(values, fmt.Sprint(value))
				}
			}

			if m, err := prometheus.NewConstMetric(
				prometheus.NewDesc(fmt.Sprint("smartthings_attribute_", attributeId), "", labels, nil),
				prometheus.GaugeValue,
				metricValue,
				values...); err == nil {

				metrics <- m
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
	case "contactSensor":
		if str, ok := value.(string); ok {
			if str == "open" {
				resultValue = 1
			}
			extras["state"] = str
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
