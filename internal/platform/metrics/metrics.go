package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics owns a private registry so process and application metrics are not
// mixed with the default global registry.
type Metrics struct {
	registry *prometheus.Registry
}

func New() *Metrics {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		prometheus.NewBuildInfoCollector(),
	)
	return &Metrics{registry: registry}
}

func (m *Metrics) Registry() *prometheus.Registry {
	return m.registry
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// Mux exposes metrics on /metrics only. It is intended for a separate internal
// listener, never for the user facing router.
func (m *Metrics) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/metrics", m.Handler())
	return mux
}
