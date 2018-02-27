package discovery

import (
	"fmt"
	"log"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/logging"
	consul "github.com/hashicorp/consul/api"
)

const checkTimeoutInSeconds int = 3
const checkIntervalInSeconds int = 15

type client struct {
	consul *consul.Client
}

func RegisterInConsul(appConfig *config.Config, logger *log.Logger) error {
	client, err := newConsulClient(appConfig.Consul.Address)
	if err != nil {
		return err
	}

	if err := client.register(appConfig); err != nil {
		return err
	}

	if err != nil {
		logging.Error(logger, "Error in Consul registration", err)
	} else {
		logging.Info(logger, "Service registered in Consul")
	}

	return nil
}

func DeregisterFromConsul(appConfig *config.Config, logger *log.Logger) error {
	client, err := newConsulClient(appConfig.Consul.Address)
	if err != nil {
		return err
	}

	if err := client.deregister(appConfig.Consul.Name); err != nil {
		return err
	}

	if err != nil {
		logging.Error(logger, "Error in Consul deregistration", err)
	} else {
		logging.Info(logger, "Service deregistered from Consul")
	}

	return nil
}

func newConsulClient(address string) (*client, error) {
	consulConfig := consul.DefaultConfig()
	consulConfig.Address = address

	c, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	return &client{consul: c}, nil
}

func (c *client) deregister(serviceID string) error {
	return c.consul.Agent().ServiceDeregister(serviceID)
}

func (c *client) register(appConfig *config.Config) error {
	agent := c.consul.Agent()

	registration := &consul.AgentServiceRegistration{
		ID:   appConfig.Consul.Name,
		Name: appConfig.Consul.Name,
		Port: appConfig.Listen.Port,
	}

	if err := agent.ServiceRegister(registration); err != nil {
		return err
	}

	serviceListenAddress := appConfig.Listen.Address

	if serviceListenAddress == "0.0.0.0" {
		serviceListenAddress = "localhost"
	}

	return agent.CheckRegister(&consul.AgentCheckRegistration{
		ID:        fmt.Sprintf("%v_IsAlive", appConfig.Consul.Name),
		Name:      "IsAlive",
		ServiceID: appConfig.Consul.Name,
		AgentServiceCheck: consul.AgentServiceCheck{
			Interval: fmt.Sprintf("%vs", checkIntervalInSeconds),
			Timeout:  fmt.Sprintf("%vs", checkTimeoutInSeconds),
			HTTP:     fmt.Sprintf("http://%v:%v/metrics", serviceListenAddress, appConfig.Listen.Port),
			Status:   "passing",
		},
	})
}
