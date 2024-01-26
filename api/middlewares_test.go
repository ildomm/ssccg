package api

import (
	"github.com/ildomm/ssccg/test_helpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRecoverMiddlewarePanicRecovery tests that the RecoverMiddleware handles panics and logs them.
func TestRecoverMiddlewarePanicRecovery(t *testing.T) {
	logBuf, restoreLog := test_helpers.CaptureOutput()
	defer restoreLog()

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	recoverMiddleware := NewRecoverMiddleware()

	testServer := httptest.NewServer(recoverMiddleware(panicHandler))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "middleware did not handle panic correctly")

	assert.Contains(t, logBuf.String(), "test panic", "log does not contain the panic message")
	assert.Contains(t, logBuf.String(), "ERROR", "log does not contain the error level")
}

// TestRecoverMiddlewareNoPanic tests the RecoverMiddleware without panics.
func TestRecoverMiddlewareNoPanic(t *testing.T) {
	normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	recoverMiddleware := NewRecoverMiddleware()

	testServer := httptest.NewServer(recoverMiddleware(normalHandler))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "middleware incorrectly handled normal request")
}

// TestLoggingMiddleware tests the LoggingMiddleware's logging of requests.
func TestLoggingMiddleware(t *testing.T) {
	logBuf, restoreLog := test_helpers.CaptureOutput()
	defer restoreLog()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	loggingMiddleware := NewLoggingMiddleware()

	testServer := httptest.NewServer(loggingMiddleware(testHandler))
	defer testServer.Close()

	_, err := http.Get(testServer.URL)
	assert.NoError(t, err)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "INFO", "log does not contain info level")
	assert.Contains(t, logOutput, "202", "log does not contain correct status code")
	assert.Contains(t, logOutput, "ms", "log does not contain execution time")
}
