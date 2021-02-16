package exporter

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var collectors = map[string]interface{}{}
var PairDelimiter = ","
var KVDelimiter = ":"

func CreateFilter(jsonFilters string) error {

	log.Print("dasdsa 3")

	type filter struct {
		Field     string `json:"field"`
		Collector string `json:"collector"`
	}

	var filters []filter
	if err := json.Unmarshal([]byte(jsonFilters), &filters); err != nil {
		return err
	}

	// create and register Prometheus Collectors
	for _, f := range filters {
		var collector prometheus.Collector
		switch f.Collector {
		case "gauge":
			collector = prometheus.NewGauge(prometheus.GaugeOpts{
				Name: f.Field,
			})

		case "counter":
			collector = prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: f.Field,
				},
			)

		case "histogram":
			collector = prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Name: f.Field,
				},
			)

		case "summary":
			collector = prometheus.NewSummary(
				prometheus.SummaryOpts{
					Name: f.Field,
				},
			)

		default:
			log.Warnf("Unknown collector %s for field %s", f.Collector, f.Field)
			continue
		}

		prometheus.MustRegister(collector)
		collectors[f.Field] = collector
		log.Printf("register %s collector for field %s", f.Collector, f.Field)
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
		collector, ok := collectors[key]
		if !ok {
			return
		}

		// validate value
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("Invalid pair %s=%s", key, value)
			return
		}

		// update collector with value
		switch v := collector.(type) {
		case prometheus.Gauge:
			v.Set(val)
		case prometheus.Counter:
			v.Add(val)
		case prometheus.Histogram:
			v.Observe(val)
		case prometheus.Summary:
			v.Observe(val)
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
