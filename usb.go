package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
)

func usbIterate() error {
	var devList []string
	ctx := gousb.NewContext()
	defer ctx.Close()

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		classify := usbid.Classify(desc)+": "

		for _, cfg := range desc.Configs {
			for _, intf := range cfg.Interfaces {
				for _, ifSetting := range intf.AltSettings {
					if !strings.Contains(classify, usbid.Classify(ifSetting)) {
						classify = strings.Join([]string{classify, usbid.Classify(ifSetting)}, " // ")
					}
				}
			}
		}

		devList = append(devList,
			fmt.Sprintf("Bus %03d Device %03d: ID %s:%s %s [%s]",
				desc.Bus,
				desc.Address,
				desc.Vendor,
				desc.Product,
				usbid.Describe(desc),
				strings.Trim(strings.Replace(classify, ":  // ", ": ", -1), ": "),
			),
		)
		return false
	})

	// All Devices returned from OpenDevices must be closed.
	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	// OpenDevices can occaionally fail, so be sure to check its return value.
	if err != nil {
		return err
	}

	log.Printf("Found %d devices\n%s", len(devList), strings.Join(devList, "\n"))

	if len(devList) < 1 {
		return errors.New("No USB devices found!")
	}

	return nil
}
