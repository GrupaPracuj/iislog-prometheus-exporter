IIS-to-Prometheus log file exporter
=====================================

Service that continuously reads newest IIS log files from given directory, parses them using [Grok](https://github.com/vjeantet/grok) and exports them to [Prometheus](https://github.com/prometheus/client_golang).

Configuration file
------------------

Configuration needs to be placed in the config.hcl file. [HCL](https://github.com/hashicorp/hcl) format is must have. Here's an example file:

```hcl
listen {
#  port = 9511               # Default is 9511 
#  address = "0.0.0.0"       # Default is "0.0.0.0"
}

consul {
#  enable = true                        # Default is true
#  name = "IISLogExporter"              # Default is "IISLogExporter"
#  address = "http://localhost:8500"    # Default is "http://localhost:8500"
#  deregister_on_service_stop = true    # Default is false
}

metric {
  metric_prefix = "iis"
  labels = ["app", "method", "status", "uri", "endpoint"]
}

logger {
#  output_log_dir = "C:\\logs"    # Default is "<exeDirectory>\logs"
#  rotate_every_mb = 10           # Default is 10
#  number_of_files = 0            # Default is 0
#  files_max_age = 0              # Default is 0
}

site "Example1" {
  logs_dir = "C:\\example1"
  label_rules = [
    {
      label_name = "app"                    # Name pointing to label from metric section
      fixed_value = "Example1"              # This exact value will be set to label
    },
    {
      label_name = "method"
      source = "method"                     # Field name from grok. Label value will be set basing on this field. Default is value of label_name

      rules = [
        {
          pattern = "copyFromSource"        # Copy exact value from grok field to label
        }
      ]
    },
    {
      label_name = "status"
      source = "status"
      rules = [
        {
          pattern = "copyFromSource"
        }
      ]
    },
    {
      label_name = "uri"
      source = "uri"
      rules = [
        {
          pattern = "copyFromSource"
        }
      ]
    },
    {
      label_name = "endpoint"
      source = "uri"
      rules = [
        {
         pattern = "/example/{someId}/endpoint"
         label_value = "exampleEndpoint"              # If source matches pattern then label will be set to this value
        },
        {
         pattern = "/help{+path}"
         label_value = "anotherExampleEndpoint"
        }]
    }]
  pattern = "%{TIMESTAMP_ISO8601:logtime} %{WORD:sitename} %{NOTSPACE:computername} %{IPORHOST:ip} %{WORD:method} %{NOTSPACE:uri} %{NOTSPACE:query} %{NUMBER:port} %{NOTSPACE:username} %{IPORHOST:client_ip} %{NOTSPACE:useragent} %{NOTSPACE:referrer} %{IPORHOST:address}(?:\\:%{NUMBER})? %{NUMBER:status} %{NUMBER:substatus} %{NUMBER:win32_status} %{NUMBER:bytes_sent} %{NUMBER:bytes_received} %{NUMBER:time_taken}"
}
```

Consul
------------
In `config.hcl` you can configure integration with Consul. You can provide values for following options:

- `enable` - Enables Consul registration
- `name` - Name of the service
- `address` - Consul adrress
- `deregister_on_service_stop = true` - Enables deregistration from Consul when service is stopped

Metrics
------------

You can also configure labels for your metrics and prefix value.

This exporter collects the following metrics:

- `<prefix>_http_response_count_total` - The total amount of processed HTTP requests/responses.
- `<prefix>_http_response_size_bytes` - A summary vector of the total amount of transferred content in response in bytes
- `<prefix>_http_request_size_bytes` - A summary vector of the total amount of transferred content in request in bytes
- `<prefix>_http_response_time_miliseconds` - A summary vector of the total response times in miliseconds.

Sites
--------

You can define multiple site sections in config file. It's mandatory to specify directory with IIS log files to observe. Also you have to provide grok pattern to parse log lines. Make sure you have `bytes_sent`, `bytes_received` and `time_taken` fields in your grok pattern.
For every label specified in `metric` section you should specify a rule. There are 3 options:
1. Set label value to fixed string (only `label_name` and `fixed_value` are mandatory)
2. Copy value from one of the grok fields (`label_name` and `pattern = "copyFromSource"` are mandatory)
3. Use uri patterns (rfc6570) to match one of the grok fields and set label to custom value (`label_name`, `pattern` and `label_value` are mandatory)

Notice: In options 2 and 3 `source` is set the same as `label_name` by default. Eg. for label "method" exporter will be looking for grok field "method". You can overwrite this behaviour providing your own `source` field in `config.hcl`.

Labels without any rule are set to "UNKNOWN".

In 3 option label will be matched according to [rfc6570](https://tools.ietf.org/html/rfc6570) standard. For instance:
- label pattern `example/{someId}/endpoint` will match `example/123/endpoint`, `example/456/endpoint`, replacing id's with {someId} token
- label pattern `help{+path}` will match `help`, `help/example`, `help/123/files` and any other uri with prefix `help`

Build, install and run
-------------------------

To build the exporter just run `build.ps1` script. It will create a bin folder with exe file and a copy of config.hcl file.

To install exporter as Windows service, just run `IISLogExporter.exe` with `install` argument.

You can run the service by clicking on the `Start` in Windows Services Manager.

Metric will be available on `<address>:<port>/Metrics` endpoint.

Credits
-------

- [tail](https://github.com/hpcloud/tail), MIT License
- [Prometheus Go client library](https://github.com/prometheus/client_golang), Apache License
- [HashiCorp configuration language](https://github.com/hashicorp/hcl), Mozilla Public License 2.0
- [RFC 6570 Uri template matcher](https://github.com/yosida95/uritemplate), BSD 3-clause "New" or "Revised" License
- [Grok parser](https://github.com/vjeantet/grok), Apache License 2.0
- [File system notifications for Go](https://github.com/fsnotify/fsnotify), BSD 3-clause "New" or "Revised" License
- [Go Date Parser](https://github.com/araddon/dateparse), MIT License
- [Go CLI Library](https://github.com/mitchellh/cli), Mozilla Public License 2.0
