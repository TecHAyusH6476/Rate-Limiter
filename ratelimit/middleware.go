package ratelimit

import (
	"net/http"
	"strings"
)

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: extract domain and descriptors from headers or URL
		domain := r.Header.Get("X-RateLimit-Domain")
		if domain == "" {
			domain = "default"
		}
		descriptors := map[string]string{}
		for k, v := range r.Header {
			if strings.HasPrefix(strings.ToLower(k), "x-ratelimit-desc-") && len(v) > 0 {
				// Extract the part after the prefix and convert to lower-case
				descKey := strings.ToLower(k[len("X-RateLimit-Desc-"):])
				descriptors[descKey] = v[0]
			}
		}
		allowed, _ := rl.Allow(domain, descriptors)

		if !allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("429 - Too Many Requests"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
