module exporter

go 1.15

require (
	github.com/fnikolai/go-collection/exporter v0.0.0-20210215172911-1a42f4f00308 // indirect
	github.com/fnikolai/go-collection/terminal v0.0.0-00010101000000-000000000000 // indirect
	github.com/prometheus/client_golang v1.9.0
)

replace github.com/fnikolai/go-collection/terminal => ../terminal
