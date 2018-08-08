package bluetooth

import (
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

var serviceList = make([]*gatt.Service, 0)

func RegisterService(s *gatt.Service) error {
	serviceList = append(serviceList, s)

	return nil
}

func StartServer() error {
	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		return err
	}

	// Register optional handlers.
	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) { log.Println("Connect: ", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) { log.Println("Disconnect: ", c.ID()) }),
	)

	// A mandatory handler for monitoring device state.
	onStateChanged := func(d gatt.Device, s gatt.State) {
		log.Printf("State: %s\n", s)
		switch s {
		case gatt.StatePoweredOn:
			d.SetServices(serviceList)

			// Advertise device name and service's UUIDs.
			d.AdvertiseNameAndServices("PiNickel", []gatt.UUID{})

			// Advertise as an OpenBeacon iBeacon
			d.AdvertiseIBeacon(gatt.MustParseUUID("AA6062F098CA42118EC4193EB73CCEB6"), 1, 2, -59)

		default:
		}
	}

	d.Init(onStateChanged)

	return nil
}
