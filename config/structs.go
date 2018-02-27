package config

// StartupFlags is a struct containing filename option that can be passed via the
// command line
type StartupFlags struct {
	ConfigFile string
}

// Config models the application's configuration
type Config struct {
	Listen ListenConfig `hcl:"listen"`
	Sites  []SiteConfig `hcl:"site"`
	Metric MetricConfig `hcl:"metric"`
	Logger LoggerConfig `hcl:"logger"`
	Consul ConsulConfig `hcl:"consul"`
}

// ListenConfig is a struct describing the built-in webserver configuration
type ListenConfig struct {
	Port    int
	Address string
}

type LoggerConfig struct {
	OutputDir     string `hcl:"output_log_dir"`
	RotateEveryMb int    `hcl:"rotate_every_mb"`
	FilesNumber   int    `hcl:"number_of_files"`
	MaxAge        int    `hcl:"files_max_age"`
}

type MetricConfig struct {
	MetricPrefix string   `hcl:"metric_prefix"`
	Labels       []string `hcl:"labels"`
}

type SiteConfig struct {
	Name       string        `hcl:",key"`
	LogsDir    string        `hcl:"logs_dir"`
	Pattern    string        `hcl:"pattern"`
	LabelRules []LabelConfig `hcl:"label_rules"`
}

type LabelConfig struct {
	Name       string       `hcl:"label_name"`
	Source     string       `hcl:"source"`
	FixedValue string       `hcl:"fixed_value"`
	Rules      []RuleConfig `hcl:"rules"`
}

type RuleConfig struct {
	Pattern    string `hcl:"pattern"`
	LabelValue string `hcl:"label_value"`
}

type ConsulConfig struct {
	Enable                  bool   `hcl:"enable"`
	Name                    string `hcl:"name"`
	Address                 string `hcl:"address"`
	DeregisterOnServiceStop bool   `hcl:"deregister_on_service_stop"`
}
