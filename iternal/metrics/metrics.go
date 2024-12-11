package metrics

import "github.com/prometheus/client_golang/prometheus"

var RequestCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "request_count",
	Help: "Number of requests received.",
})
var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "endpoint"},
)

func init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestDuration)
}
