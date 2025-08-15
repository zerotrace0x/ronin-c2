package middleware

import (
	"net/http"
	"os"
	"strings"
)

func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		required := os.Getenv("RONIN_C2_API_KEY")
		if required == "" {
			http.Error(w, "Server misconfigured: RONIN_C2_API_KEY empty", http.StatusInternalServerError)
			return
		}
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != required {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
