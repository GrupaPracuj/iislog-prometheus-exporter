package lib

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"
	"github.com/GrupaPracuj/iislog-prometheus-exporter/logging"

	"github.com/araddon/dateparse"
	"github.com/hpcloud/tail"
	"github.com/vjeantet/grok"
)

func processOneSite(site config.SiteConfig, metrics *Metrics, logger *log.Logger, metricCfg config.MetricConfig) {
	logging.Info(logger, fmt.Sprintf("Starting listener for site %s", site.Name))

	grokParser := newGrok(site)
	var tails = []*tail.Tail{}

	newestFileName := findNewestFile(site.LogsDir, logger)
	if newestFileName != "" {
		logging.Info(logger, fmt.Sprintf("Found newest log file: %s for site %s", newestFileName, site.Name))
		tails = append(tails, processOneFile(site.LogsDir+"\\"+newestFileName, grokParser, metrics, site, logger, metricCfg, time.Now().UTC()))
	}

	//observing directory
	newFileAlert, err := newFileCheck(site.LogsDir, logger)
	if err != nil {
		logging.Error(logger, "", err)
	}

	//when new file is created tail it
	for fileName := range newFileAlert {
		logging.Info(logger, fmt.Sprintf("New log file created: %s for site %s", fileName, site.Name))
		tails = append(tails, processOneFile(fileName, grokParser, metrics, site, logger, metricCfg, time.Time{}))

		if len(tails) > 1 {
			tails[len(tails)-2].StopAtEOF() //We tail new file so there is no need to tail the previous one
			logging.Info(logger, fmt.Sprintf("Tail for file %s cleaned up for site %s", tails[len(tails)-2].Filename, site.Name))
		}
	}
}

//Tail one file and modify metrics
func processOneFile(filename string, grokParser *grok.Grok, metrics *Metrics, siteCfg config.SiteConfig, logger *log.Logger, metricCfg config.MetricConfig, firstExecutionTime time.Time) (t *tail.Tail) {
	//Tail file
	t, err := tail.TailFile(filename, tail.Config{
		Follow: true,
		ReOpen: true,
		Poll:   true,
	})
	
	if err != nil {
		logging.Error(logger, fmt.Sprintf("Error tailing file %s", filename), err)
	}
	
	//Set labels and metrics
	go func(siteCfg config.SiteConfig, metricCfg config.MetricConfig) {
		metricsCount := len(metricCfg.Labels)
		labelValues := make([]string, metricsCount)

		//Parse log line one by one
		for line := range t.Lines {
			line.Text = strings.TrimRight(line.Text, "\r")
			fields, err := grokParser.Parse("%{IIS}", line.Text)
			if err != nil {
				logging.Error(logger, fmt.Sprintf("Parsing line %s for site %s", siteCfg.Name, line.Text), err)
				continue
			}
			if len(fields) == 0 {
				logging.Info(logger, fmt.Sprintf("For site %s line not parsed: %s", siteCfg.Name, line.Text))
				continue
			}

			if logtime := fields["logtime"]; len(logtime) > 0 {
				timestamp, _ := dateparse.ParseAny(logtime)
				if !firstExecutionTime.IsZero() && timestamp.Before(firstExecutionTime) {
					continue
				}
			}

			if fields["status"] == "0" {
				logging.Info(logger, fmt.Sprintf("Status 0 for line: %s", line.Text))
			}

			for i := 0; i < metricsCount; i++ {
				labelValues[i] = "UNKNOWN"
			}

			for _, label := range siteCfg.LabelRules {
				var labelSource string
				if len(label.Source) < 1 {
					labelSource = label.Name
				} else {
					labelSource = label.Source
				}

				labelIndex := findIndexOfLabel(metricCfg.Labels, label.Name)
				if labelIndex == -1 {
					logging.Error(logger, fmt.Sprintf("There is no %s label configured in metric section. Change your config.", label.Name), nil)
				} else {
					labelValues[labelIndex] = SetLabel(&label, fields[labelSource])
				}
			}

			metrics.countTotal.WithLabelValues(labelValues...).Inc()

			if bytes := fields["bytes_sent"]; len(bytes) > 0 {
				if bytesFloat, err := strconv.ParseFloat(bytes, 64); err == nil {
					metrics.bytesSentTotal.WithLabelValues(labelValues...).Observe(bytesFloat)
				}
			}

			if bytes := fields["bytes_received"]; len(bytes) > 0 {
				if bytesReceivedFloat, err := strconv.ParseFloat(bytes, 64); err == nil {
					metrics.bytesReceivedTotal.WithLabelValues(labelValues...).Observe(bytesReceivedFloat)
				}
			}

			if responseTime := fields["time_taken"]; len(responseTime) > 0 {
				if responseTimeFloat, err := strconv.ParseFloat(responseTime, 64); err == nil {
					metrics.responseMiliSeconds.WithLabelValues(labelValues...).Observe(responseTimeFloat)
				}
			}
		}
	}(siteCfg, metricCfg)

	return t
}

func findIndexOfLabel(labels []string, label string) int {
	for i, v := range labels {
		if v == label {
			return i
		}
	}
	return -1
}
