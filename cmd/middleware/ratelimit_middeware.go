// Package middleware responsible for application middlewares
package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/rodrigosuco/go-proxy/cmd/ratelimit"
)

var Buckets sync.Map

func RateLimitMiddleware(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid IP address",
			})
			return
		}

		var currentBucket *ratelimit.Bucket

		if val, ok := Buckets.Load(ip); ok {
			currentBucket = val.(*ratelimit.Bucket)
		} else {
			newBucket := &ratelimit.Bucket{
				Tokens:                    ratelimit.AvailableTokensForNewBucket,
				LastTokenRestoreTimeStamp: time.Now(),
			}

			val, _ := Buckets.LoadOrStore(ip, newBucket)
			currentBucket = val.(*ratelimit.Bucket)

		}

		if err := ratelimit.CheckLimit(currentBucket); err != nil {
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
