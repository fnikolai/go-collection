package exporter

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/prometheus/client_golang/prometheus"
)

var metrics = map[string]interface{}{}
var PairDelimiter = ","
var KVDelimiter = ":"

func CreateFilter(jsonFilters string) error {

	type filter struct {
		Field  string `json:"field"`
		Metric string `json:"metric"`
	}

	var filters []filter
	if err := json.Unmarshal([]byte(jsonFilters), &filters); err != nil {
		return err
	}

	// create and register Prometheus Collectors
	for _, f := range filters {
		switch f.Metric {
		case "gauge":
			metric := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: f.Field,
			})
			prometheus.MustRegister(metric)
			metrics[f.Field] = metric

			log.Print("register gauge for field ", f.Field)

		case "counter":
			metric := prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: f.Field,
				},
			)
			prometheus.MustRegister(metric)
			metrics[f.Field] = metric

			log.Print("register counter for field ", f.Field)
		default:
			log.Printf("Unknown metric %s for field %s", f.Metric, f.Field)
		}
	}

	return nil
}

func SpaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func ApplyFilter(line string) {

	filter := func(key, value string) {
		// validate key
		metric, ok := metrics[key]
		if !ok {
			return
		}

		// validate value
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("Invalid pair %s=%s", key, value)
			return
		}

		// update metric with value
		switch v := metric.(type) {
		case prometheus.Gauge:
			v.Set(val)
		case prometheus.Counter:
			v.Add(val)
		}
	}

	for _, pair := range strings.Split(line, PairDelimiter) {
		x := strings.Split(pair, KVDelimiter)
		if len(x) != 2 {
			continue
		}
		filter(x[0], x[1])
	}
}

//hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
