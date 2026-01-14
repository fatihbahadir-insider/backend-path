package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	once     sync.Once
	Registry *prometheus.Registry

	HttpRequestsTotal     *prometheus.CounterVec
	HttpRequestDuration   *prometheus.HistogramVec
	HttpRequestsInFlight  prometheus.Gauge
	TransactionsTotal     *prometheus.CounterVec
	TransactionAmount     *prometheus.HistogramVec
	ActiveUsers           prometheus.Gauge
	DatabaseQueriesTotal  *prometheus.CounterVec
	DatabaseQueryDuration *prometheus.HistogramVec
)

func Init() {
	once.Do(func() {
		Registry = prometheus.NewRegistry()

		Registry.MustRegister(collectors.NewGoCollector())
		Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

		HttpRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		)

		HttpRequestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "app_http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		)

		HttpRequestsInFlight = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
		)

		TransactionsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "transactions_total",
				Help: "Total number of transactions",
			},
			[]string{"type", "status"},
		)

		TransactionAmount = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "transaction_amount",
				Help:    "Transaction amounts distribution",
				Buckets: []float64{10, 50, 100, 500, 1000, 5000, 10000},
			},
			[]string{"type"},
		)

		ActiveUsers = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users_total",
				Help: "Total number of active users",
			},
		)

		DatabaseQueriesTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"operation", "table"},
		)

		DatabaseQueryDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_query_duration_seconds",
				Help:    "Duration of database queries in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		)

		Registry.MustRegister(
			HttpRequestsTotal,
			HttpRequestDuration,
			HttpRequestsInFlight,
			TransactionsTotal,
			TransactionAmount,
			ActiveUsers,
			DatabaseQueriesTotal,
			DatabaseQueryDuration,
		)
	})
}
