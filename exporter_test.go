package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/setheck/smartthings-exporter/smartthings"
)

func TestPlayground(t *testing.T) {
	t.SkipNow()

	client := smartthings.NewClient(os.Getenv("STE_API_TOKEN"), nil)
	//caps, err := client.ListAllCapabilities(nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(smartthings.ToString(caps))
	//
	//_, err = client.GetCapabilitiesByIDAndVersion(caps[0].ID, caps[0].Version)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//prof, err := client.ListAllDeviceProfiles(nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(smartthings.ToString(prof))

	devices, err := client.ListDevices(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	for _, device := range devices {
		for _, component := range device.Components {
			cs, err := client.GetDeviceComponentStatus(context.TODO(), device.DeviceID, component.ID)
			if err != nil {
				fmt.Println(err)
				//t.Fatal(err)
			} else {
				//fmt.Println(cs)
				for capabilityId, attributes := range cs {
					for attributeId, properties := range attributes {
						if value, ok := properties["value"]; ok {
							fmt.Println("capabilityId:", capabilityId, "attribute:", attributeId, "value:", value)
						}
					}
				}
			}
		}
	}
}

func TestParseTime(t *testing.T) {
	t.SkipNow()
	ts := "2020-12-03T06:41:54.441Z"
	tm, _ := time.Parse(time.RFC3339Nano, ts)
	fmt.Println(tm.UnixNano() / int64(time.Millisecond))
}
