package metrics

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func Register() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		glog.Error(err)
	}

}

const (
	MetricsNamespace = "httpserver"
)

var functionLatency = CreateExecutionTimeMetric(MetricsNamespace, "Time spend")

type ExecutionTimer struct {
	histogram *prometheus.HistogramVec
	start     time.Time
	last      time.Time
}

func (t *ExecutionTimer) ObserveTotal() {
	(*t.histogram).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

func NewExecutionTimer(histogram *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histogram: histogram,
		start:     now,
		last:      now,
	}

}

func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)

}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"step"},
	)

}
