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
			m, err := prometheus.NewConstMetric(
				prometheus.NewDesc("iot", "", []string{"deviceId"}, nil),
				prometheus.GaugeValue,
				1,
				device.DeviceID)
			if err != nil {
				log.Println(err)
			} else {
				metrics <- m
			}
		}
	}
}
