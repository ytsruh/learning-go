package engine

import (
	"fmt"
	"net/http"
	"time"

	"ytsruh.com/saas/cache"
)

// RateLimiter middleware used to prevent too many call in short time span
func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var keys Auth

		ctx := r.Context()
		v := ctx.Value(ContextAuth)
		if v == nil {
			keys = Auth{}
		} else {
			a, ok := v.(Auth)
			if ok {
				keys = a
			}
		}

		key := fmt.Sprintf("%v", keys.AccountID)

		// TODO: Make this configurable
		count, err := cache.RateLimit(key, 1*time.Minute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Make this configurable
		if count >= 60 {
			// we get the expiration duration of this key so we can notify the user
			d, err := cache.GetThrottleExpiration(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if d.Seconds() > 0 {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(d.Seconds())))
			}
			http.Error(w, fmt.Sprintf("you've reached your rate limit, retry in %v", d), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
