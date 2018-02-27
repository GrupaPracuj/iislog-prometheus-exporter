package lib

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/logging"

	"github.com/prometheus/client_golang/prometheus"
)

func ExportLogs(cfg *config.Config, logger *log.Logger) {
	metrics := Metrics{}
	metrics.Init(&cfg.Metric)

	enableHelpPage(cfg)

	for _, site := range cfg.Sites {
		logging.Info(logger, fmt.Sprintf("Start processing site: %s", site.Name))
		go processOneSite(site, &metrics, logger, cfg.Metric)
	}

	listenAddr := fmt.Sprintf("%s:%d", cfg.Listen.Address, cfg.Listen.Port)
	logging.Info(logger, fmt.Sprintf("Running HTTP server on address %s", listenAddr))

	http.Handle("/metrics", prometheus.Handler())
	go http.ListenAndServe(listenAddr, nil)
}
