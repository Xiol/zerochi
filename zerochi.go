package zerochi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// LogMessage is the log message associated with each log line.
var LogMessage = "handled request"

var silentPaths map[string]struct{}

// SetSilentPaths disables logging for the provided paths, e.g. "/health", "/metrics"
func SetSilentPaths(paths ...string) {
	silentPaths = make(map[string]struct{}, len(paths))
	for _, path := range paths {
		silentPaths[path] = struct{}{}
	}
}

func Logger(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			t2 := time.Now()

			if rec := recover(); rec != nil {
				log.Error().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("ip", r.RemoteAddr).
					Str("user_agent", r.UserAgent()).
					Dur("latency", t2.Sub(t1)).
					Interface("recovered", rec).
					Msg("panic recovered")
				http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if _, ok := silentPaths[r.URL.Path]; ok {
				return
			}

			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", ww.Status()).
				Str("ip", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Int64("bytes_in", r.ContentLength).
				Int("bytes_out", ww.BytesWritten()).
				Dur("latency", t2.Sub(t1)).
				Msg(LogMessage)
		}
		return http.HandlerFunc(fn)
	}
}
