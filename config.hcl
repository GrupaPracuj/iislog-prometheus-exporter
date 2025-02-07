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

/*
site "Example1" {
  logs_dir = "C:\\example1"
  label_rules = [
    {
      label_name = "app"                    # Name pointing to label from metric section
      fixed_value = "Example1"              # This exact value will be set to label
    },
    {
      label_name = "method"
      source = "method"                     # Field name from grok. Label value will be set basing on this field. Default is value of label_name.
      rules = [
        {
          pattern = "copyFromSource"        # Exactly copy value from grok field to label
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

site "Example2" {
  logs_dir = "C:\\example2"
  label_rules = [
    {
      label_name = "app"
      fixed_value = "Example2"
    },
    {
      label_name = "method"
      source = "method"
      rules = [
        {
          pattern = "copyFromSource"
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
         pattern = "/some/endpoint"
         label_value = "someEndpoint"
        },
        {
         pattern = "/another/endpoint/{someId}"
         label_value = "anotherEndpoint"
        }]
    }]
  pattern = "%{TIMESTAMP_ISO8601:logtime} %{WORD:sitename} %{NOTSPACE:computername} %{IPORHOST:ip} %{WORD:method} %{NOTSPACE:uri} %{NOTSPACE:query} %{NUMBER:port} %{NOTSPACE:username} %{IPORHOST:client_ip} %{NOTSPACE:useragent} %{NOTSPACE:referrer} %{IPORHOST:address}(?:\\:%{NUMBER})? %{NUMBER:status} %{NUMBER:substatus} %{NUMBER:win32_status} %{NUMBER:bytes_sent} %{NUMBER:bytes_received} %{NUMBER:time_taken}"
}
*/
