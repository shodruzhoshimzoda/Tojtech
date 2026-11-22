package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

// LimitBody - a request body bigger than maxBytes is cut: reading it fails with an error,
// so a huge body cannot eat the memory of the server.
func LimitBody(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

type visitor struct {
	count   int
	resetAt time.Time
}

// RateLimit - lets one IP address make at most `limit` requests per `window`.
// One RateLimit(...) value has one shared counter, so the same value on several routes
// counts them together.
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*visitor)
	)

	// allow says whether this client may make one more request and, if not, how long to wait
	allow := func(key string) (bool, time.Duration) {
		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		// throw away expired entries now and then, so the map does not grow forever
		if len(visitors) > 10_000 {
			for k, v := range visitors {
				if now.After(v.resetAt) {
					delete(visitors, k)
				}
			}
		}

		v, ok := visitors[key]
		if !ok || now.After(v.resetAt) {
			visitors[key] = &visitor{count: 1, resetAt: now.Add(window)}
			return true, 0
		}
		if v.count >= limit {
			return false, time.Until(v.resetAt)
		}
		v.count++
		return true, 0
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			if ok, wait := allow(ip); !ok {
				w.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds())+1))
				httphelpers.RespondWarn(r.Context(), w, r, http.StatusTooManyRequests,
					"rate limit exceeded", "too many attempts, please try again in a minute")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}