package middleware

import (
	"errors"
	"expvar"
	"net/http"
	"strconv"

	"github.com/felixge/httpsnoop"
	"github.com/seanflannery10/core/internal/auth"
	"github.com/seanflannery10/core/internal/httperrors"
)

func RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := auth.GetUser(r)

		if authenticatedUser == "" {
			httperrors.AuthenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalRequestsReceived.Add(1)
		totalResponsesSent.Add(1)
		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())
		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err, _ := rvr.(error)
				if !errors.Is(err, http.ErrAbortHandler) {
					httperrors.ServerError(w, r, err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
