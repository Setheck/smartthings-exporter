package main

import (
	"log"

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

			//for _, component := range device.Components {
			//	if componentStatus, err := collector.client.GetDeviceComponentStatus(device.DeviceID, component.ID); err == nil {
			//		for id,cmp := range componentStatus {
			//			switch id {
			//			case "switch":
			//				if m, err := prometheus.NewConstMetric(
			//					prometheus.NewDesc("smartthings_device", "", []string{"deviceId", "deviceLabel", "name"}, nil),
			//					prometheus.GaugeValue,
			//					1,
			//					device.DeviceID, device.Label, device.Name); err == nil {
			//
			//					metrics <- m
			//				}
			//			}
			//		}
			//	}
			//}
		}
	}
}
