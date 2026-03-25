package http_helper

import (
	"io"
	"net/http"
	"testing"

	"QA-TA-SPRING-2026/Task2.1/internal/utils"

	"github.com/stretchr/testify/require"
)

// ReadResponseBody reads the response body and closes the connection.
func ReadResponseBody(t testing.TB, resp *http.Response) string {
	defer resp.Body.Close() //nolint:errcheck

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")

	return string(body)
}

// AssertStatusCode checks the HTTP status code and returns the body.
// Also asserts Content-Type is application/json when body is non-empty.
func AssertStatusCode(t testing.TB, resp *http.Response, expectedStatusCode int) string {
	body := ReadResponseBody(t, resp)

	if len(body) > 0 {
		contentType := resp.Header.Get("Content-Type")
		require.Contains(t, contentType, "application/json",
			"expected Content-Type application/json, got %s", contentType)
	}

	require.Equalf(t, expectedStatusCode, resp.StatusCode,
		"expected HTTP status %d, got %d. Body: %s",
		expectedStatusCode, resp.StatusCode, body)

	utils.LogWithLabelAndTimeStamp("HTTP", "Status: "+resp.Status)

	return body
}
