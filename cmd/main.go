package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rodrigosuco/go-proxy/cmd/middleware"
	"github.com/rodrigosuco/go-proxy/cmd/ratelimit"
)

func main() {
	var buckets []ratelimit.Bucket
	target, err := url.Parse("http://localhost:3000/api/v1")
	if err != nil {
		log.Fatal("Error parsing target URL: ", err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			fmt.Println(r.Out)
			r.SetURL(target)
			r.Out.Host = r.In.Host
			r.SetXForwarded()
		},
	}

	handler := middleware.RateLimitMiddleware(proxy, buckets)
	fmt.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
