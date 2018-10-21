package lib

import (
	"fmt"
	"net/http"

  "github.com/GrupaPracuj/iislog-prometheus-exporter/config"
  "github.com/GrupaPracuj/iislog-prometheus-exporter/version"
)

const helpPageTemplate string = `<!doctype html>
<html>  
  <head>
    <meta charset="utf-8">
    <title>IISLogExporter</title>
    <style type="text/css">
      body { font-family: Calibri; }
    </style>
  </head>

  <body>
    <h2>You are running IISLogExporter version %s</h2>
    <p>Click <a href="/metrics">here</a> to view your metrics.</p>
  </body>
</html>`

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, helpPageTemplate, version.Version)
}

func enableHelpPage(appConfig *config.Config) {
	http.HandleFunc("/", handler)
}
