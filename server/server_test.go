package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davyj0nes/generic-web-handler/server"
	"gotest.tools/v3/assert"
)

// To allow for a slice of testCases with generic input we need to use an interface
// and then specific structs for the different cases with different input types.
type testCases interface {
	GetName() string
	GetExpected() float64
	GetInput() []any
}

type intTestCase struct {
	name     string
	input    []int
	expected float64
}

func (i intTestCase) GetName() string {
	return i.name
}

func (i intTestCase) GetExpected() float64 {
	return i.expected
}

func (i intTestCase) GetInput() []any {
	var in []any
	for _, a := range i.input {
		in = append(in, a)
	}

	return in
}

type floatTestCase struct {
	name     string
	input    []float64
	expected float64
}

func (f floatTestCase) GetName() string {
	return f.name
}

func (f floatTestCase) GetExpected() float64 {
	return f.expected
}

func (f floatTestCase) GetInput() []any {
	var in []any
	for _, e := range f.input {
		in = append(in, e)
	}

	return in
}

func TestServer(t *testing.T) {
	tt := []testCases{
		&intTestCase{
			name:     "withInts",
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8},
			expected: 36,
		},
		&floatTestCase{
			name:     "withFloats",
			input:    []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8},
			expected: 40,
		},
	}

	for _, tc := range tt {
		t.Run(tc.GetName(), func(t *testing.T) {
			s := server.NewServer("not used")
			srv := httptest.NewServer(s)
			defer srv.Close()

			got := makeRequest(t, srv, tc.GetInput())
			assert.Equal(t, got.Sum, tc.GetExpected())
		})
	}
}

func TestServerNotFound(t *testing.T) {
	s := server.NewServer("not used")
	srv := httptest.NewServer(s)
	defer srv.Close()

	res, err := srv.Client().Get(srv.URL)
	assert.NilError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func makeRequest(t *testing.T, srv *httptest.Server, input []any) server.SumResponse {
	var reqBody bytes.Buffer

	// This is a bit of a hack to get the input into the correct type. It's not clear
	// to me at the moment why this is needed so will need to do some more research.
	switch input[0].(type) {
	case int:
		var in []int

		for _, a := range input {
			in = append(in, a.(int))
		}
		err := json.NewEncoder(&reqBody).Encode(server.SumRequest[int]{Input: in})
		assert.NilError(t, err)
	case float64:
		var in []float64
		for _, a := range input {
			in = append(in, a.(float64))
		}
		err := json.NewEncoder(&reqBody).Encode(server.SumRequest[float64]{Input: in})
		assert.NilError(t, err)
	}

	res, err := srv.Client().Post(srv.URL+"/sum", "application/json", &reqBody)
	assert.NilError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	var got server.SumResponse
	err = json.NewDecoder(res.Body).Decode(&got)
	assert.NilError(t, err)

	return got
}
