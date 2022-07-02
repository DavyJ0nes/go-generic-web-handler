package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/davyj0nes/generic-web-handler/server"
)

type Logger interface {
	Printf(string, ...any)
	Fatal(v ...any)
}

type Client struct {
	addr   string
	logger Logger
}

func NewClient(addr string, logger Logger) *Client {
	if logger == nil {
		logger = log.New(os.Stdout, "client: ", 0)
	}

	return &Client{
		addr:   addr,
		logger: logger,
	}
}

func (c *Client) Run() {
	intInput := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 14, 15,
	}
	intExpected := float64(96)

	intRes, err := sumRequest(c.addr, intInput)
	if err != nil {
		c.logger.Fatal(err)
	}

	logResponse(intInput, c.logger, intExpected, intRes.Sum)

	floatInput := []float64{
		1.9, 2.8, 3.7, 4.6, 5.5, 6.4, 7.3, 8.2, 9.1,
	}
	floatExpected := 49.0

	floatRes, err := sumRequest(c.addr, floatInput)
	if err != nil {
		c.logger.Fatal(err)
	}

	logResponse(floatInput, c.logger, floatExpected, floatRes.Sum)
}

// make sumRequest to Server
func sumRequest[T server.Number](addr string, input []T) (server.SumResponse, error) {
	var reqBody bytes.Buffer
	if err := json.NewEncoder(&reqBody).Encode(server.SumRequest[T]{Input: input}); err != nil {
		return server.SumResponse{}, fmt.Errorf("encoding int request: %w", err)
	}

	res, err := http.Post(addr, "application/json", &reqBody)
	if err != nil {
		return server.SumResponse{}, fmt.Errorf("int request failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return server.SumResponse{},
			fmt.Errorf("int request failed with err: %s", res.Status)
	}

	return parseSumResponse(res.Body)
}

// logResponse from Server.
// logger needs to be parsed as a parameter because methods cannot have type parameters.
// more info: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#no-parameterized-methods
func logResponse[Num server.Number](input []Num, logger Logger, expected, sum float64) {
	logger.Printf("sum of %v = %.2f\n", input, sum)
	if sum != expected {
		logger.Printf("output does not match! %f != %f", sum, expected)
	}
}

func parseSumResponse(responseBody io.ReadCloser) (server.SumResponse, error) {
	var output server.SumResponse
	if err := json.NewDecoder(responseBody).Decode(&output); err != nil {
		return server.SumResponse{}, fmt.Errorf("reading body failed: %w", err)
	}

	if err := responseBody.Close(); err != nil {
		return server.SumResponse{}, fmt.Errorf("closing body failed: %w", err)
	}

	return output, nil
}
