package main

import (
  "github.com/ufoscout/go-up"
  // "github.com/xlab/closer"
)

func init() {
  // setup any config defaults

  // setup any options for the CLI

  // add any actions
  actions["ui"] = action{
    Usage:    "",
    Purpose:  "starts the Pi Car Go user interface",
    Function: func(config go_up.GoUp, args []string) error {
      // do whatever needs doing

      return nil
    },
  }
}
