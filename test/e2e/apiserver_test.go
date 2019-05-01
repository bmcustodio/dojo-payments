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
	"net/http"

	request "github.com/imroc/req"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bmcstdio/dojo-payments/pkg/server"
)

var _ = Describe("API Server", func() {
	When(`receiving a "GET /" HTTP request`, func() {
		var (
			err error
			res *request.Resp
		)

		JustBeforeEach(func() {
			// Make a "GET /" request.
			res, err = request.Get(baseUrl)
		})

		It(`returns "200 OK"`, func() {
			// Make sure that no errors have occurred, and that "200 OK" was returned.
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
		})

		It("returns a value indicating that the database is online", func() {
			// Make sure that no errors have occurred, and that there is a key indicating that the database is online.
			Expect(err).NotTo(HaveOccurred())
			body := server.APIServerRootResponse{}
			err = res.ToJSON(&body)
			Expect(err).NotTo(HaveOccurred())
			Expect(body.DatabaseStatus).To(Equal(server.DatabaseStatusOnline))
		})
	})
})
