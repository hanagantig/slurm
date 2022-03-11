package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "myapp_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

func main() {

	r := mux.NewRouter()
	r.Use(prometheusMiddleware)

	r.Path("/metrics").Handler(promhttp.Handler())

	r.Path("/ping").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("serving ping")
			fmt.Fprintln(w, "pong")
		},
	)

	srv := &http.Server{Addr: ":8080", Handler: r}
	log.Fatal(srv.ListenAndServe())
}
