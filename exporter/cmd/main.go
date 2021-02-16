package main

import (
	"bufio"
	"flag"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/fnikolai/go-collection/terminal"

	"github.com/fnikolai/go-collection/exporter"
)

// [{"field": "key1",  "metric": "value1"}, {"field": "key2",  "metric": "value2"}]
func main() {
	terminal.SetLogger("debug", true)

	addr := flag.String("listen-address", ":9443", "Address on which to expose metrics")
	filter := flag.String("filter", "{}", "JSON filters for field extraction")
	flag.Parse()

	if *filter == "{}" {
		log.Fatal("No filter was specified")
	}

	// create metric endpoints
	if err := exporter.CreateFilter(*filter); err != nil {
		log.Fatal("filter creation failed ", err)
	}

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()



	log.Println("Press ctrl + c to terminate")
	terminal.HandleSignals(nil, func() error {
		// Read from input and export filtered fields
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			exporter.ApplyFilter(exporter.SpaceStringsBuilder(scanner.Text()))
		}

		return scanner.Err()
	})
}
