package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

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

	handler := rateLimitMiddleware(proxy, buckets)
	fmt.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func rateLimitMiddleware(proxy *httputil.ReverseProxy, buckets []ratelimit.Bucket) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ratelimit.CheckLimit(r.RemoteAddr, buckets)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "rate limit exceeded",
			})
			return
		}

		proxy.ServeHTTP(w, r)
	})
}
