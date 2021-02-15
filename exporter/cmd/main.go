package main

import (
	"bufio"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/fnikolai/go-collection/exporter"
)

// [{"field": "key1",  "metric": "value1"}, {"field": "key2",  "metric": "value2"}]
func main() {
	addr := flag.String("listen-address", ":9080", "Address on which to expose metrics")
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

	// Read from input and export filtered fields
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		exporter.ApplyFilter(exporter.SpaceStringsBuilder(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	log.Println("Press ctrl + c to terminate")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
	}
}
