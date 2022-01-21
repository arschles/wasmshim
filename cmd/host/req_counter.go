package main

import (
	"net/http"

	prom "github.com/prometheus/client_golang/prometheus"
)

func reqCounterGauge() prom.Gauge {
	g := prom.NewGauge(prom.GaugeOpts{
		Namespace: "wasmshim_host",
		Subsystem: "http",
		Name:      "pending_requests",
		Help:      "The number of requests currently awaiting response from a WASM module.",
	})
	prom.MustRegister(g)
	return g
}

func reqCounterMiddleware(h http.Handler, g prom.Gauge) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.Inc()
		h.ServeHTTP(w, r)
		g.Dec()
	})
}
