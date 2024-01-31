package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

// RecoverMiddleware recover any panics from upstream
// and return a 500 response.
//
// Recovered panics are logged with a stack trace
type RecoverMiddleware struct {
	response RecoveryResponse
}

// RecoveryResponse
// This is the signature of the function that will be called when a panic is recovered.
// The function should write the error, building the response as best suit, into the responseWriter and return.
type RecoveryResponse func(ctx context.Context, w http.ResponseWriter)

// stackTraceSize defines the size of the stack trace to be captured when a panic is recovered, in bytes.
const stackTraceSize = 64 << 10

// NewRecoverMiddleware initializes a new RecoverMiddleware
// with the default response
func NewRecoverMiddleware() func(next http.Handler) http.Handler {
	return NewRecoverMiddlewareWithCustomResponse(defaultRecoveryResponse())
}

// NewRecoverMiddlewareWithCustomResponse initializes a new RecoverMiddleware
// with a custom response
func NewRecoverMiddlewareWithCustomResponse(response RecoveryResponse) func(next http.Handler) http.Handler {
	return RecoverMiddleware{
		response: response,
	}.perform
}

// perform is the middleware handler itself
func (rm RecoverMiddleware) perform(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Recover from any panic upstream,
		// any panic not specifically handled by the upstream handlers/middlewares
		defer func() {
			if err := recover(); err != nil {

				stackTrace := make([]byte, stackTraceSize)
				stackTrace = stackTrace[:runtime.Stack(stackTrace, false)]

				switch x := err.(type) {
				case string:
					err = errors.New(x)
				default:
					err = fmt.Errorf("unknown panic: %w", x.(error))
				}

				// Log the error and stack trace
				log.Printf("ERROR: recovering from error. Error: %s. StackTrace: %s \n", err, stackTrace)

				rm.response(r.Context(), w)
			}
		}()

		// Call the next handler as a normal flow execution
		next.ServeHTTP(w, r)
	})
}

// defaultRecoveryResponse returns a default response for the RecoverMiddleware
// this response is a JSON string, and it is the default response.
// The response is a ErrorResponse
// The response pattern can be changed by calling NewRecoverMiddlewareWithCustomResponse, passing a new function
// that implements the same signature as RecoveryResponse
func defaultRecoveryResponse() RecoveryResponse {
	return func(ctx context.Context, w http.ResponseWriter) {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{"Internal error"})
	}
}

// StatusRecorder
// Source:
// https://upgear.io/blog/golang-tip-wrapping-http-response-writer-for-middleware/
type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	if status != 200 {
		r.ResponseWriter.WriteHeader(status)
	}
}

// LoggingMiddleware is a middleware that logs the request
type LoggingMiddleware struct{}

// NewLoggingMiddleware initializes a new LoggingMiddleware
func NewLoggingMiddleware() func(next http.Handler) http.Handler {
	return LoggingMiddleware{}.perform
}

// perform is the middleware handler itself
func (lm LoggingMiddleware) perform(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
		}

		start := time.Now().UnixNano()

		// Call the next handler as a normal flow execution
		next.ServeHTTP(recorder, r)

		end := time.Now().UnixNano()
		duration := (end - start) / int64(time.Millisecond)

		// Logs execution time of the request and other details
		log.Printf("INFO: %s \"%s %s\" %d %dms\n", r.RemoteAddr, r.Method, r.URL.Path, recorder.Status, duration)
	})
}
