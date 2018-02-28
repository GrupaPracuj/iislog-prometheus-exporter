package lib

import (
	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"

	"github.com/vjeantet/grok"
)

func newGrok(cfg config.SiteConfig) (grokParser *grok.Grok) {
	g, _ := grok.NewWithConfig(&grok.Config{NamedCapturesOnly: true})
	g.AddPattern("IIS", cfg.Pattern)

	return g
}
