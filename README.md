# Go Generic Web Handler

This is an exploration of how Go type parameters, introduced in Go 1.18, can be used
with an HTTP handler as well as in table test cases where the input can be generic.

None of this is intended to be a best practice or even a good idea but purely as
an experiment to see how to use type parameters for generic programming.

For more information on generic programming with Go see:

- [An Introduction to Generics](https://go.dev/blog/intro-generics)
- [Tutorial: Getting started with generics](https://go.dev/doc/tutorial/generics)
- [Refactor Cloud applications in Go 1.18 with generics](https://www.youtube.com/watch?v=-F2t3oInqKE)
- [Type Parameter Discussion](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)

## Structure

This example includes a simple server that handles summing slices of float64 or int
from a POST request. The code for this can be found in [server](./server).

In [client](./client) there is a client that makes a couple of requests with
floats and ints to easily verify the server and be able to experiment with
type parameters from both server and client side.

I've added comments in both packages with some oddities and where I've found 
boundaries with what Go Generics can be used at the moment.

## Running

There is a [Makefile](./Makefile) to make running the experiment and its tests.

- `make run` - Run [main.go](./main.go)
- `make test` - Run the tests for server and client packages.
