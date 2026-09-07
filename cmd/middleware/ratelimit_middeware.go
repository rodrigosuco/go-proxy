// Package middleware responsible for application middlewares
package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/rodrigosuco/go-proxy/cmd/ratelimit"
)

var Buckets = make(map[string]*ratelimit.Bucket, 0)

func RateLimitMiddleware(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid IP address",
			})
			return
		}

		bucket, exists := Buckets[ip]

		if !exists {
			Buckets[ip] = &ratelimit.Bucket{
				IP:                        ip,
				Tokens:                    15,
				LastTokenRestoreTimeStamp: time.Now(),
			}
		} else {
			if err := ratelimit.CheckLimit(bucket); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "rate limit exceeded",
				})
				return
			}
		}
		proxy.ServeHTTP(w, r)
	})
}
