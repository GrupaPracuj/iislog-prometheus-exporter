package service

import (
	"fmt"
	"log"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/discovery"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/lib"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/logging"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

var elog debug.Log

type iisLogExporterService struct {
	Logger *log.Logger
	Config *config.Config
}

func (s *iisLogExporterService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	if s.Config.Consul.Enable {
		discovery.RegisterInConsul(s.Config, s.Logger)
	}

	go lib.ExportLogs(s.Config, s.Logger)

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop:
				if s.Config.Consul.DeregisterOnServiceStop {
					discovery.DeregisterFromConsul(s.Config, s.Logger)
				}
				logging.Info(s.Logger, fmt.Sprintf("Stopping service #%d", c))
				break loop
			case svc.Shutdown:
				discovery.DeregisterFromConsul(s.Config, s.Logger)
			default:
				continue
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func Run(name string, isDebug bool, cfg *config.Config, logger *log.Logger) {
	logging.Info(logger, fmt.Sprintf("Starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err := run(name, &iisLogExporterService{Config: cfg, Logger: logger})
	if err != nil {
		logging.Error(logger, fmt.Sprintf("%s service failed", name), err)
		return
	}
	logging.Info(logger, fmt.Sprintf("%s service stopped", name))
}
