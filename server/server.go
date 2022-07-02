package server

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"golang.org/x/exp/constraints"
)

type Server struct {
	addr   string
	logger *log.Logger
	mux    *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func NewServer(port string) *Server {
	s := &Server{
		addr:   ":" + port,
		logger: log.New(os.Stdout, "server: ", 0),
		mux:    http.NewServeMux(),
	}

	s.mux.Handle("/", http.NotFoundHandler())
	s.mux.HandleFunc("/sum", s.SumHandler)

	return s
}

func (s *Server) Start() error {
	s.logger.Println("starting server...")
	return http.ListenAndServe(s.addr, s.mux)
}

// Number is a custom constraint to communicate that this implementation only
// caters for Integers and Floats.
type Number interface {
	constraints.Integer | constraints.Float
}

// SumRequest defines the input for SumHandler
type SumRequest[T Number] struct {
	Input []T `json:"input"`
}

// SumResponse defines the output from SumHandler
type SumResponse struct {
	Sum float64 `json:"sum"`
}

// SumHandler is a POST request handler for summing slices of float64s or ints together
func (s *Server) SumHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := parseRequestBody(r.Body)
	if err != nil {
		s.respondError(w, err, http.StatusBadRequest)
		return
	}

	response, err := newSumResponse(reqBody)
	if err != nil {
		s.respondError(w, err, http.StatusBadRequest)
		return
	}

	s.logger.Printf("returning: %+v", response)

	if err := respondOK(w, response); err != nil {
		s.respondError(w, err, http.StatusInternalServerError)
	}
}

func respondOK(w http.ResponseWriter, res SumResponse) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func (s *Server) respondError(w http.ResponseWriter, err error, code int) {
	s.logger.Println(err)
	http.Error(w, err.Error(), code)
}

// This is a bit of a hack but to be able to use a type parameter on the SumRequest
// we need to see what fits, so this tries to unmarshal into each of the types
// and once successful it returns.
func parseRequestBody(body io.ReadCloser) (any, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var f SumRequest[float64]
	var i SumRequest[int]
	if err := json.Unmarshal(b, &i); err == nil {
		return i, nil
	}

	if err := json.Unmarshal(b, &f); err == nil {
		return f, nil
	}

	if err := body.Close(); err != nil {
		return nil, err
	}

	return nil, err
}

// This works in conjunction with parseRequestBody to allow for converting from
// an unknown request body into a struct with type parameters. This isn't very
// pretty and does look a lot like how this could have been handled pre Generics.
func newSumResponse(reqBody any) (res SumResponse, err error) {
	switch t := reqBody.(type) {
	case SumRequest[float64]:
		res.Sum = sumInput(t.Input)
	case SumRequest[int]:
		res.Sum = sumInput(t.Input)
	default:
		err = errors.New("unknown request type")
	}

	return
}

// Once the types are known then type parameters becomes easy to read.
// am returning a float64 here rather than T because of not wanting to do the
// dance again of having a generic SumResponse struct.
func sumInput[T Number](input []T) float64 {
	var sum T
	for _, i := range input {
		sum += i
	}

	return math.Round((float64(sum) * 100) / 100)
}
