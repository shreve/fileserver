package main

import (
	"log"
	"time"
	"net/http"
)

type Middleware func(http.Handler) http.Handler;

func Logging(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Started [%s] %s", r.Method, r.URL.Path)
		start := time.Now()
		f.ServeHTTP(w, r)
		end := time.Now()
		diff := float64(end.Sub(start)) / float64(time.Microsecond)
		log.Printf("Completed [%s] %s (%.2f μs)", r.Method, r.URL.Path, diff)
		log.Printf("")
	})
}

func CORS(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*");
		w.Header().Set("Content-Type", "application/json");
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,POST,PUT,OPTIONS");
		if r.Method == "OPTIONS" { return; }
		f.ServeHTTP(w, r)
	})
}

func M(action http.HandlerFunc, middles ...Middleware) http.HandlerFunc {
	output := http.Handler(action)
	for _, middle := range middles {
		output = middle(output)
	}
	return output.ServeHTTP
}
