package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	backends []*Backend
	current  uint64
}

func (lb *LoadBalancer) getNextBackend() *Backend {
	next := atomic.AddUint64(&lb.current, 1)
	return lb.backends[next%uint64(len(lb.backends))]
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend()
	if backend != nil && backend.Alive {
		backend.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func newBackend(urlStr string) *Backend {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	return &Backend{
		URL:          parsedURL,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}
}

func newLoadBalancer(urls []string) *LoadBalancer {
	backends := make([]*Backend, len(urls))
	for i, url := range urls {
		backends[i] = newBackend(url)
	}
	return &LoadBalancer{
		backends: backends,
	}
}

func main() {
	backendURLs := []string{
		"http://192.168.1.2:8000", // Nodo del Servidor 1
		"http://192.168.1.3:8000", // Nodo del Servidor 2
	}

	lb := newLoadBalancer(backendURLs)

	server := http.Server{
		Addr:    ":8000",
		Handler: lb,
	}

	log.Println("Load Balancer started on :8000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
