package main

import (
  "github.com/ufoscout/go-up"
)

func init() {
  // setup any config defaults

  // setup any options for the CLI

  // add any actions
  actions["hardware"] = action{
    Usage:    "",
    Purpose:  "runs a basic hardware check",
    Function: func(config go_up.GoUp, args []string) {
      // do whatever needs doing

      // Check audio out - beep/speaker
      // Check USB?
      // Check Bluetooth
      // Check network uplink
      // Check OBD II
      // Check GPS
      // Check network downlink
      // Check audio in
    },
  }
}
