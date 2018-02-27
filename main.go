package main

import (
	"fmt"
	"os"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/logging"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/service"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/version"
	"github.com/mitchellh/cli"
)

const defaultServiceName = "IISLogExporter"

func main() {

	var serviceName string
	if len(os.Args) < 3 {
		serviceName = defaultServiceName
	} else {
		serviceName = os.Args[2]
	}

	c := cli.NewCLI(serviceName, version.Version)
	c.Args = os.Args[1:]
	c.HiddenCommands = []string{service.RunAsDaemonArg}
	c.Commands = map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			return &versionCommand{
				version: version.Version,
			}, nil
		},
		"debug": func() (cli.Command, error) {
			cfg, confErr := config.LoadConfig()
			if confErr != nil {
				panic(confErr)
			}
			logger := logging.Init(cfg, true)
			return &debugCommand{
				serviceName: serviceName,
				config:      cfg,
				logger:      logger,
			}, nil
		},
		"install": func() (cli.Command, error) {
			return &installCommand{
				serviceName: serviceName,
			}, nil
		},
		"remove": func() (cli.Command, error) {
			cfg, confErr := config.LoadConfig()
			if confErr != nil {
				panic(confErr)
			}
			logger := logging.Init(cfg, false)
			return &removeCommand{
				serviceName: serviceName,
				config:      cfg,
				logger:      logger,
			}, nil
		},
		service.RunAsDaemonArg: func() (cli.Command, error) {
			cfg, confErr := config.LoadConfig()
			if confErr != nil {
				panic(confErr)
			}
			logger := logging.Init(cfg, false)
			return &daemonCommand{
				config: cfg,
				logger: logger,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
