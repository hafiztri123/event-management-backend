package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/hafiztri123/src/internal/pkg/logger"
)

type Logger struct {

}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger)  LoggerMiddleware(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			requestID := r.Header.Get("X-Request-ID")
			if requestID =="" {
				requestID = middleware.GetReqID(r.Context())
			}

			ctx := logger.WithRequestID(r.Context(), requestID)

			if claims, ok := ctx.Value("user").(map[string]interface{}); ok {
				if userID, ok := claims["user_id"].(string); ok {
					ctx = logger.WithUserID(ctx, userID)
				}
			}

			traceID := r.Header.Get("X-Trace-ID")
			ctx = logger.WithTraceID(ctx, traceID)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			log.Info(ctx, "Request started", map[string]interface{}{
				"method": 		r.Method,
				"path": 		r.URL.Path,
				"remote_ip": 	r.RemoteAddr,
				"user_agent": 	r.UserAgent(),
				"referer": 		r.Referer(),
			})

			next.ServeHTTP(ww, r.WithContext(ctx))


			duration := time.Since(start)

			data := map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      ww.Status(),
				"duration_ms": duration.Milliseconds(),
				"bytes":       ww.BytesWritten(),
			}

			if ww.Status() >= 500 {
				log.Error(ctx, "Request completed with server error", nil, data)
			} else if ww.Status() >= 400 {
				log.Warn(ctx, "Request completed with client error", data)
			} else {
				log.Info(ctx, "Request completed successfully", data)
			}

		})
	}
}