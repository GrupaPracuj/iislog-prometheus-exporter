package lib

import (
	"strings"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"

	"github.com/yosida95/uritemplate"
)

func SetLabel(config *config.LabelConfig, labelSource string) string {
	if config.FixedValue == "" {
		for _, rule := range config.Rules {
			if rule.Pattern == "copyFromSource" {
				return labelSource
			}
			template := uritemplate.MustNew(strings.ToLower(rule.Pattern))

			match := template.Match(strings.ToLower(strings.TrimSuffix(labelSource, "/")))
			if match != nil {
				return rule.LabelValue
			}
		}
	} else {
		return config.FixedValue
	}

	return "NOT MATCHED"
}
