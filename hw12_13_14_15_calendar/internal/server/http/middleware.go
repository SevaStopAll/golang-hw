package internalhttp

import (
	"fmt"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
		next.ServeHTTP(w, r)
	})
}
