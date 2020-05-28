module github.com/pberndro/smartpi_exporter

go 1.14

require (
	github.com/pberndro/smartpi_exporter/ade7878 v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.6.0
	github.com/prometheus/common v0.10.0
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/exp v0.0.0-20200513190911-00229845015e
)

replace github.com/pberndro/smartpi_exporter/ade7878 => ./ade7878
