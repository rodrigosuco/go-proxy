package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rodrigosuco/go-proxy/cmd/middleware"
)

func main() {
	target, err := url.Parse("http://localhost:3000/api/v1")
	if err != nil {
		log.Fatal("Error parsing target URL: ", err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.SetXForwarded()
		},
	}

	handler := middleware.RateLimitMiddleware(proxy)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
