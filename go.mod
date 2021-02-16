module github.com/fnikolai/go-collection

go 1.15

replace (
	github.com/fnikolai/go-collection/terminal  => ./terminal
	github.com/fnikolai/go-collection/exporter   => ./exporter

)

require (
	github.com/prometheus/client_golang v1.9.0
	github.com/sirupsen/logrus v1.7.1
	github.com/urfave/cli/v2 v2.3.0
)
