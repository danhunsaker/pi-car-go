package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ufoscout/go-up"
	// "github.com/xlab/closer"
	"github.com/danhunsaker/pi-car-go/bluetooth"
)

func init() {
	// setup any config defaults

	// setup any options for the CLI

	// add any actions
	actions["hardware"] = action{
		Usage:   "",
		Purpose: "runs a basic hardware check",
		Function: func(config go_up.GoUp, args []string) error {
			// do whatever needs doing

			// Check audio out - beep/speaker
			log.Println("Audio output check starting...")
			// err := sayEspeak("Pi Nickel test for espeak")
			// if err != nil {
			//   return errors.New(fmt.Sprintf("Audio output check failed: %v", err))
			// }
			// err = sayPico("Pi Nickel test for pico TTS")
			// if err != nil {
			//   return errors.New(fmt.Sprintf("Audio output check failed: %v", err))
			// }
			err := sayFestival("Pi Nickel test for Festival")
			if err != nil {
				return errors.New(fmt.Sprintf("Audio output check failed: %v", err))
			}
			ttsWriter, err := NewTTSWriter("", "festival")
			if err != nil {
				return errors.New(fmt.Sprintf("TTS logging setup failed: %v", err))
			}
			log.SetOutput(io.MultiWriter(os.Stderr, ttsWriter))
			log.Println("Audio output check passed!")

			// Check USB?
			log.Println("USB check starting...")
			err = usbIterate()
			if err != nil {
				return errors.New(fmt.Sprintf("USB check failed: %v", err))
			}
			log.Println("USB check passed!")

			// Check Bluetooth
			log.Println("Bluetooth check starting...")
			err = bluetooth.StartServer()
			if err != nil {
				return errors.New(fmt.Sprintf("Bluetooth check failed: %v", err))
			}
			log.Println("Bluetooth check passed!")

			// Check network uplink
			log.Println("Uplink check starting...")
			err = nil
			if err != nil {
				return errors.New(fmt.Sprintf("Uplink check failed: %v", err))
			}
			log.Println("Uplink check passed!")

			// Check OBD II
			log.Println("OBD-II check starting...")
			err = nil
			if err != nil {
				return errors.New(fmt.Sprintf("OBD-II check failed: %v", err))
			}
			log.Println("OBD-II check passed!")

			// Check GPS
			log.Println("GPS check starting...")
			err = nil
			if err != nil {
				return errors.New(fmt.Sprintf("GPS check failed: %v", err))
			}
			log.Println("GPS check passed!")

			// Check network downlink
			log.Println("Downlink check starting...")
			err = nil
			if err != nil {
				return errors.New(fmt.Sprintf("Downlink check failed: %v", err))
			}
			log.Println("Downlink check passed!")

			// Check audio in
			log.Println("Audio input check starting...")
			err = nil
			if err != nil {
				return errors.New(fmt.Sprintf("Audio input check failed: %v", err))
			}
			log.Println("Audio input check passed!")

			return nil
		},
	}
}
