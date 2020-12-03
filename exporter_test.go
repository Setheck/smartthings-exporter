package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/setheck/smartthings-exporter/smartthings"
)

func TestPlayground(t *testing.T) {
	t.SkipNow()

	client := smartthings.NewClient(os.Getenv("STE_API_TOKEN"))
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

	devices, err := client.ListDevices()
	if err != nil {
		t.Fatal(err)
	}

	for _, device := range devices {
		for _, component := range device.Components {
			_, err := client.GetDeviceComponentStatus(device.DeviceID, component.ID)
			if err != nil {
				fmt.Println(err)
				//t.Fatal(err)
			}
		}
	}
}
