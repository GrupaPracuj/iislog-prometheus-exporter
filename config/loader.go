package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl"
)

const configFile string = "config.hcl"

func LoadConfig() (*Config, error) {
	configFile := getConfigFileDir()
	cfg, err := loadConfigFromFile(configFile)

	if cfg.Listen.Port == 0 {
		cfg.Listen.Port = 9511
	}

	if len(cfg.Listen.Address) == 0 {
		cfg.Listen.Address = "0.0.0.0"
	}

	if len(cfg.Metric.MetricPrefix) == 0 {
		cfg.Metric.MetricPrefix = "iis"
	}

	if cfg.Consul.Address == "" {
		cfg.Consul.Address = "http://localhost:8500"
	}

	if cfg.Consul.Name == "" {
		cfg.Consul.Name = "IISLogExporter"
	}

	cfg.Logger.OutputDir = strings.TrimSuffix(cfg.Logger.OutputDir, "\\")
	cfg.Logger.OutputDir = strings.TrimSuffix(cfg.Logger.OutputDir, "/")

	for _, site := range cfg.Sites {
		site.LogsDir = strings.TrimSuffix(site.LogsDir, "\\")
		site.LogsDir = strings.TrimSuffix(site.LogsDir, "/")

		for _, label := range site.LabelRules {
			if label.Rules != nil {
				for i := 0; i < len(label.Rules); i++ {
					rule := &label.Rules[i]
					rule.Pattern = strings.TrimSuffix(rule.Pattern, "/")
				}
			}
		}
	}

	return cfg, err
}

// LoadConfigFromFile fills a configuration object (passed as parameter) with
// values read from a configuration file (pass as parameter by filename). The
// configuration file needs to be in HCL format.
func loadConfigFromFile(filepath string) (*Config, error) {
	cfg := Config{Consul: ConsulConfig{Enable: true}}

	if _, err := os.Stat(filepath); err == nil {
		buf, err := ioutil.ReadFile(filepath)
		if err != nil {
			return &cfg, err
		}

		hclText := string(buf)

		err = hcl.Decode(&cfg, hclText)
		if err != nil {
			return &cfg, err
		}
	}

	return &cfg, nil
}

func getConfigFileDir() (dir string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return dir + "\\" + configFile
}
