package exporter

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode"

	//dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var collectors = map[string]interface{}{}
var PairDelimiter = ","
var KVDelimiter = ":"

// CreateFilter generates filters according to the specific format.
// If metric is undefined, then it takes the value of Field. Practically, Metric is used for renaming
// erroneous fields like 99th(us)
func CreateFilter(jsonFilters string) error {

	type filter struct {
		Field     string `json:"field"`
		Metric    string `json:"metric"`
		Collector string `json:"collector"`
	}

	var filters []filter
	if err := json.Unmarshal([]byte(jsonFilters), &filters); err != nil {
		return err
	}

	// create and register prometheus Collectors
	for _, f := range filters {

		if f.Field == "" || f.Collector == "" {
			log.Warn("Invalid filter ", f)
			continue
		}

		if f.Metric == "" {
			f.Metric = f.Field
		}

		var collector prometheus.Collector
		switch f.Collector {
		case "gauge":
			collector = prometheus.NewGauge(prometheus.GaugeOpts{
				Name: f.Metric,
			})

		case "counter":
			collector = prometheus.NewCounter(
				prometheus.CounterOpts{
					Name: f.Metric,
				},
			)

		case "histogram":
			collector = prometheus.NewHistogram(
				prometheus.HistogramOpts{
					Name: f.Metric,
				},
			)

		case "summary":
			collector = prometheus.NewSummary(
				prometheus.SummaryOpts{
					Name: f.Metric,
				},
			)

		default:
			log.Warnf("Unknown collector %s for field %s", f.Collector, f.Field)
			continue
		}

		prometheus.MustRegister(collector)
		collectors[f.Field] = collector
		log.Printf("register %s collector (%s) for field %s", f.Collector, f.Metric, f.Field)
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

/* https://github.com/prometheus/client_golang/blob/master/prometheus/examples_test.go
func ExampleNewMetricWithTimestamp() {
	desc := prometheus.NewDesc(
		"temperature_kelvin",
		"Current temperature in Kelvin.",
		nil, nil,
	)

	// NewTestingEnvironment a constant gauge from values we got from an external
	// temperature reporting system. Those values are reported with a slight
	// delay, so we want to add the timestamp of the actual measurement.
	temperatureReportedByExternalSystem := 298.15
	timeReportedByExternalSystem := time.Date(2009, time.November, 10, 23, 0, 0, 12345678, time.UTC)
	s := prometheus.NewMetricWithTimestamp(
		timeReportedByExternalSystem,
		prometheus.MustNewConstMetric(
			desc, prometheus.GaugeValue, temperatureReportedByExternalSystem,
		),
	)

	// Just for demonstration, let's check the state of the gauge by
	// (ab)using its Write method (which is usually only used by prometheus
	// internally).
	metric := &dto.Metric{}
	s.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

	// Output:
	// gauge: <
	//   value: 298.15
	// >
	// timestamp_ms: 1257894000012
}

//hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()


*/
