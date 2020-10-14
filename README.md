# tempCheck

Read temperature from a DS18B20 thermometer, format, and write for [Prometheus Node Exporter](https://github.com/prometheus/node_exporter) textfile directory.


## Running

`/usr/local/bin/tempCheck -d /sys/bus/w1/devices -t /var/lib/node_exporter/textfile_collector/probe_temps.prom`

### Args:
- -d, --devicePath "Path to directory with device info."
- -t, --textfileExporterPath "Path to directory for node_exporter textfile collector."
