package main

import (
	"fmt"
	"path/filepath"
	"os"

	goflag "flag"
	flag "github.com/spf13/pflag"
  "github.com/ufoscout/go-up"
)

// options that will be recognized by the CLI
var options = flag.CommandLine

// defaults for any configuration options
var defaults = go_up.NewGoUp()

// actions avaliable to call into
var actions = make(map[string]action)

type action struct {
	Usage    string
	Purpose  string
	Function func(go_up.GoUp, []string)
}

func main() {
	var action string
	var args []string

	// Load the Environment variables
	// The are used as they are defined, e.g. ENV_VARIABLE=XXX
	// Used just to provide values for substitution
	defaults.AddReaderWithPriority(go_up.NewEnvReader("", false, false), go_up.HighestPriority)

	// Load the Environment variables and convert their keys
	// from ENV_VARIABLE=XXX to env.variable=XXX
	defaults.AddReaderWithPriority(go_up.NewEnvReader("", true, true), go_up.HighestPriority)

	// load config file, if present
	defaults.AddFile("./picargo.conf", true)

	// parse CLI options
	flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))
    options.PrintDefaults()
		for name, act := range actions {
			fmt.Fprintf(os.Stderr, "\t%s %s\n", name, act.Usage)
			fmt.Fprintf(os.Stderr, "\t\t%s\n", act.Purpose)
		}
	}

	options.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	// pull CLI options into the configuration
	flag.Visit(func(option *flag.Flag) {
		if key, ok := option.Annotations["picargo_config_option"]; ok {
			defaults.Add(key[0], option.Value.String())
		}
	})

	// build the config
	config, err := defaults.Build()
	if err != nil {
		panic(err)
	}

	// identify selected action
	if flag.NArg() > 1 {
		action = flag.Arg(0)
		args = flag.Args()[1:]
	} else if flag.NArg() == 1 {
		action = flag.Arg(0)
		args = []string{}
	} else {
		action = ""
		args = []string{}
	}

	// call into selected action
	if _, ok := actions[action]; ok {
		actions[action].Function(config, args)
	} else {
		flag.Usage()
	}
}
