# ROOT holds the absolute path to the root of the repository.
ROOT := $(shell git rev-parse --show-toplevel)

# test.e2e runs the end-to-end test suite.
.PHONY: test.e2e
test.e2e: BASE_URL ?= http://localhost:8080
test.e2e:
	@go test $(ROOT)/test/e2e --ginkgo.v --test.v --base-url $(BASE_URL)
