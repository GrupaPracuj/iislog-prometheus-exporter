package main

import (
	"fmt"
	"log"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/discovery"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/service"
)

type versionCommand struct {
	version string
}

type debugCommand struct {
	serviceName string
	config      *config.Config
	logger      *log.Logger
}

type installCommand struct {
	serviceName string
}

type removeCommand struct {
	serviceName string
	config      *config.Config
	logger      *log.Logger
}

type daemonCommand struct {
	config *config.Config
	logger *log.Logger
}

type checkConfigCommand struct{}

func (c *versionCommand) Run(_ []string) int {
	fmt.Printf("%v\r\n", c.version)
	return 0
}

func (c *versionCommand) Help() string {
	return "Display IISLogExporter version"
}

func (c *versionCommand) Synopsis() string {
	return "Display version"
}

func (c *debugCommand) Run(_ []string) int {
	service.Run(c.serviceName, true, c.config, c.logger)
	return 0
}

func (c *debugCommand) Help() string {
	return "Debug IISLogExporter"
}

func (c *debugCommand) Synopsis() string {
	return "Debug IISLogExporter"
}

func (c *installCommand) Run(_ []string) int {
	service.Install(c.serviceName)
	fmt.Printf("install command succeded\r\n")
	return 0
}

func (c *installCommand) Help() string {
	return "Install IISLogExporter as Windows Service"
}

func (c *installCommand) Synopsis() string {
	return "Install IISLogExporter as Windows Service"
}

func (c *removeCommand) Run(_ []string) int {
	discovery.DeregisterFromConsul(c.config, c.logger)
	service.Remove(c.serviceName)
	fmt.Printf("remove command succeded\r\n")
	return 0
}

func (c *removeCommand) Help() string {
	return "Remove IISLogExporter from Windows Services"
}

func (c *removeCommand) Synopsis() string {
	return "Remove IISLogExporter from Windows Services"
}

func (c *daemonCommand) Run(_ []string) int {
	service.Run(defaultServiceName, false, c.config, c.logger)
	return 0
}

func (c *daemonCommand) Help() string {
	return "Run IISLogExporter as daemon"
}

func (c *daemonCommand) Synopsis() string {
	return "Run IISLogExporter as daemon"
}

func (c *checkConfigCommand) Run(_ []string) int {
	fmt.Printf("Configuration file is valid")
	return 0
}

func (c *checkConfigCommand) Help() string {
	return "Check if configuration file is valid"
}

func (c *checkConfigCommand) Synopsis() string {
	return "Check if configuration file is valid"
}
