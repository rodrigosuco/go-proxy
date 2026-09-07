// Package middleware responsible for application middlewares
package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/rodrigosuco/go-proxy/cmd/ratelimit"
)

func RateLimitMiddleware(proxy *httputil.ReverseProxy, buckets []ratelimit.Bucket) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := ratelimit.CheckLimit(r.RemoteAddr, buckets)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "rate limit exceeded",
			})
			return
		}

		proxy.ServeHTTP(w, r)
	})
}
