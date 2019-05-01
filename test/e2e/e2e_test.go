// Copyright 2019 Bruno Miguel Custodio
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var (
	baseUrl string
)

func init() {
	flag.StringVar(&baseUrl, "base-url", "http://localhost:8080", "the base url at which the api server can be reached")
}

var _ = BeforeSuite(func() {
	log.Infof("running the end-to-end test suite against the api server at %q", baseUrl)
})

func TestEndToEnd(t *testing.T) {
	// Parse the provided command-line flags.
	flag.Parse()
	// Register a failure handler and run the test suite.
	RegisterFailHandler(Fail)
	RunSpecs(t, "end-to-end test suite")
}
