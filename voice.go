package main

import (
  "github.com/ufoscout/go-up"
  "github.com/danhunsaker/pi-car-go/voice"
)

func init() {
  // setup any config defaults

  // setup any options for the CLI

  // add any actions
  actions["voice"] = action{
    Usage:    "",
    Purpose:  "starts the PiCarGo voice interface",
    Function: func(config go_up.GoUp, args []string) {
      // do whatever needs doing
      voice.StartListening()
      select{}
    },
  }
}
