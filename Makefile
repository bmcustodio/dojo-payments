# ROOT holds the absolute path to the root of the repository.
ROOT := $(shell git rev-parse --show-toplevel)

# run runs the API server.
.PHONY: run
run: BIND_ADDR?=localhost:8080
run:
	@go run $(ROOT)/cmd/main.go --bind-addr $(BIND_ADDR)

# test.e2e runs the end-to-end test suite.
.PHONY: test.e2e
test.e2e: BASE_URL ?= http://localhost:8080
test.e2e:
	@go test $(ROOT)/test/e2e --ginkgo.v --test.v --base-url $(BASE_URL)
