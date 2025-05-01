package server

import (
	"log"
	"net/http"
	"net/http/httputil"

	"TestCloudCamp/internal/balancer"
)

func NewReverseProxy(pool *balancer.ServerPool) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			back := pool.GetNextBackend()
			if back != nil {
				log.Printf("Proxying request to: %s", back.URL)
				req.URL.Scheme = back.URL.Scheme
				req.URL.Host = back.URL.Host
				req.Host = back.URL.Host
			}
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			http.Error(w, "Server unavailable", http.StatusServiceUnavailable)
		},
	}
}
