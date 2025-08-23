package middleware

import (
	"net/http"

	loggertag "taskL0/internal/entity/logger_tag"
	"taskL0/pkg/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.AddValuesToContext(r.Context(),
			loggertag.RequestPath, r.URL.Path,
			loggertag.RequestMethod, r.Method,
			loggertag.RequestRemoteAddr, r.RemoteAddr,
		)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
