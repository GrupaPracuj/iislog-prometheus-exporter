package lib

import (
	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	countTotal          *prometheus.CounterVec
	bytesSentTotal      *prometheus.SummaryVec
	bytesReceivedTotal  *prometheus.SummaryVec
	responseMiliSeconds *prometheus.HistogramVec
}

func (m *Metrics) Init(cfg *config.MetricConfig) {
	labels := make([]string, len(cfg.Labels))
	for i, label := range cfg.Labels {
		labels[i] = label
	}

	m.countTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: cfg.MetricPrefix,
		Name:      "http_response_count_total",
		Help:      "Amount of processed HTTP requests",
	}, labels)

	m.bytesSentTotal = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: cfg.MetricPrefix,
		Name:      "http_response_size_bytes",
		Help:      "Total amount of transferred bytes",
	}, labels)

	m.bytesReceivedTotal = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: cfg.MetricPrefix,
		Name:      "http_request_size_bytes",
		Help:      "Total amount of transferred bytes",
	}, labels)

	m.responseMiliSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: cfg.MetricPrefix,
		Name:      "http_response_time_miliseconds",
		Help:      "Time needed by IIS to handle requests",
		Buckets:   []float64{1000.0, 2000.0, 3000.0, 4000.0, 5000.0, 6000.0, 7000.0},
	}, labels)

	prometheus.MustRegister(m.countTotal)
	prometheus.MustRegister(m.bytesSentTotal)
	prometheus.MustRegister(m.bytesReceivedTotal)
	prometheus.MustRegister(m.responseMiliSeconds)
}
