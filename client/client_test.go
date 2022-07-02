package client_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davyj0nes/generic-web-handler/client"
	"gotest.tools/v3/assert"
)

func TestClient(t *testing.T) {
	var callCount int
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sum float64
		if callCount == 0 {
			sum = 96
		} else {
			sum = 49
		}

		res := struct {
			Sum float64 `json:"sum"`
		}{
			Sum: sum,
		}

		assert.NilError(t, json.NewEncoder(w).Encode(res))
		callCount++
	})

	srv := httptest.NewServer(handler)
	defer srv.Close()

	logger := &mockLogger{}
	c := client.NewClient(srv.URL, logger)
	c.Run()
	want := []string{
		"sum of [1 2 3 4 5 6 7 8 9 10 12 14 15] = 96.00\n",
		"sum of [1.9 2.8 3.7 4.6 5.5 6.4 7.3 8.2 9.1] = 49.00\n",
	}
	assert.DeepEqual(t, want, logger.got)
}

type mockLogger struct {
	got []string
}

// noop
func (m *mockLogger) Fatal(...any) {}

func (m *mockLogger) Printf(format string, a ...any) {
	m.got = append(m.got, fmt.Sprintf(format, a...))
}
