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
	"time"

	request "github.com/imroc/req"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/bmcstdio/dojo-payments/pkg/db/models"
	"github.com/bmcstdio/dojo-payments/pkg/server"
	"github.com/bmcstdio/dojo-payments/pkg/server/apis/payments"
	"github.com/bmcstdio/dojo-payments/test/e2e/util"
)

const (
	// paymentIDFieldName is the name of the "id" field of a Payment object.
	paymentIDFieldName = "ID"
)

var _ = Describe("API Server", func() {
	When(`receiving a "GET /" HTTP request`, func() {
		var (
			err error
			res *request.Resp
		)

		BeforeEach(func() {
			// Make a "GET /" request and make sure no errors have occurred.
			res, err = request.Get(baseUrl)
			Expect(err).NotTo(HaveOccurred())
		})

		It(`returns "200 OK"`, func() {
			// Make sure that a "200 OK" response was returned.
			Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
		})

		It("returns a value indicating that the database is online", func() {
			// Make sure that there is a key on the response's body indicating that the database is online.
			body := server.APIServerRootResponse{}
			err = res.ToJSON(&body)
			Expect(err).NotTo(HaveOccurred())
			Expect(body.DatabaseStatus).To(Equal(server.DatabaseStatusOnline))
		})
	})

	Context("serving the Payments API", func() {
		When(`receiving a "POST /payments" request`, func() {
			var (
				payment models.Payment
			)

			BeforeEach(func() {
				// Make sure we start with a valid payment.
				payment = models.Payment{
					Amount:      314.15,
					Currency:    "EUR",
					Date:        util.MustParseRFC3339Time("2019-04-30T22:30:00Z"),
					Description: "Order #1",
					Beneficiary: models.Entity{
						AccountNumber: "1234",
						BankID:        "4321",
						Name:          "John",
					},
					Debtor: models.Entity{
						AccountNumber: "5678",
						BankID:        "8765",
						Name:          "Dave",
					},
				}
			})

			Context("containing an invalid payment", func() {
				DescribeTable(`returns "400 BAD REQUEST" and a meaningful error message`,

					// The following function represents the test itself, which will be executed for each test case.

					func(fn func(*models.Payment), expectedErrorMessage string) {
						// Apply the transformation function to the base payment object in order to make it invalid.
						fn(&payment)
						// Attempt to create the payment and make sure that "400 BAD REQUEST" is returned.
						res, err := request.Post(baseUrl+payments.BasePath, request.BodyJSON(payment))
						Expect(err).NotTo(HaveOccurred())
						Expect(res.Response().StatusCode).To(Equal(http.StatusBadRequest))
						// Make sure that the expected error message was returned.
						resBody := echo.HTTPError{}
						err = res.ToJSON(&resBody)
						Expect(err).NotTo(HaveOccurred())
						Expect(resBody.Message).To(Equal(expectedErrorMessage))
					},

					// The following entries represent the test cases.

					Entry("when the beneficiary's account number is empty", func(p *models.Payment) {
						p.Beneficiary.AccountNumber = ""
					}, "beneficiary: the entity's account number must not be empty"),

					Entry("when the beneficiary's bank id is empty", func(p *models.Payment) {
						p.Beneficiary.BankID = ""
					}, "beneficiary: the entity's bank id must not be empty"),

					Entry("when the beneficiary's name is empty", func(p *models.Payment) {
						p.Beneficiary.Name = ""
					}, "beneficiary: the entity's name must not be empty"),

					Entry("when the debtors's account number is empty", func(p *models.Payment) {
						p.Debtor.AccountNumber = ""
					}, "debtor: the entity's account number must not be empty"),

					Entry("when the debtors's bank id is empty", func(p *models.Payment) {
						p.Debtor.BankID = ""
					}, "debtor: the entity's bank id must not be empty"),

					Entry("when the debtor's name is empty", func(p *models.Payment) {
						p.Debtor.Name = ""
					}, "debtor: the entity's name must not be empty"),

					Entry("when the amount is negative", func(p *models.Payment) {
						p.Amount = -1.34
					}, "the amount must be positive"),

					Entry("when the amount is empty", func(p *models.Payment) {
						p.Amount = 0
					}, "the amount must be positive"),

					Entry("when the currency is empty", func(p *models.Payment) {
						p.Currency = ""
					}, "the currency must not be empty"),

					Entry("when the date is empty", func(p *models.Payment) {
						p.Date = time.Time{}
					}, "the date must not be empty"),

					Entry("when the description is empty", func(p *models.Payment) {
						p.Description = ""
					}, "the description must not be empty"),
				)
			})

			Context("containing a valid payment", func() {
				It("creates the payment and returns its ID in the response's body", func() {
					req, err := request.Post(baseUrl+payments.BasePath, request.BodyJSON(payment))
					Expect(err).NotTo(HaveOccurred())
					Expect(req.Response().StatusCode).To(Equal(http.StatusCreated))
					err = req.ToJSON(&payment)
					Expect(err).NotTo(HaveOccurred())
					Expect(payment.ID).NotTo(BeEmpty())
				})
			})
		})

		When("more than one payment exists in the database", func() {
			var (
				payment1 models.Payment
				payment2 models.Payment
			)

			BeforeEach(func() {
				payment1 = models.Payment{
					Amount:      314.15,
					Currency:    "EUR",
					Date:        util.MustParseRFC3339Time("2019-04-30T22:30:00Z"),
					Description: "Order #1",
					Beneficiary: models.Entity{
						AccountNumber: "1234",
						BankID:        "4321",
						Name:          "John",
					},
					Debtor: models.Entity{
						AccountNumber: "5678",
						BankID:        "8765",
						Name:          "Dave",
					},
				}
				payment2 = models.Payment{
					Amount:      412.32,
					Currency:    "USD",
					Date:        util.MustParseRFC3339Time("2019-04-30T22:30:00Z"),
					Description: "Order #2",
					Beneficiary: models.Entity{
						AccountNumber: "9182",
						BankID:        "1312",
						Name:          "John",
					},
					Debtor: models.Entity{
						AccountNumber: "3424",
						BankID:        "2131",
						Name:          "Dave",
					},
				}

				// Create the first payment.
				res, err := request.Post(baseUrl+payments.BasePath, request.BodyJSON(payment1))
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusCreated))
				err = res.ToJSON(&payment1)
				Expect(err).NotTo(HaveOccurred())
				Expect(payment1.ID).NotTo(BeEmpty())

				// Create the second payment.
				res, err = request.Post(baseUrl+payments.BasePath, request.BodyJSON(payment2))
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusCreated))
				err = res.ToJSON(&payment2)
				Expect(err).NotTo(HaveOccurred())
				Expect(payment2.ID).NotTo(BeEmpty())
			})

			It("can find an existing payment by its ID", func() {
				// Try to get one of the payments by its ID and make sure no error has been returned.
				res, err := request.Get(baseUrl + payments.BasePath + "/" + payment1.ID.Hex())
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				// Decode the response's body and make sure the correct payment has been returned.
				result := models.Payment{}
				err = res.ToJSON(&result)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.ID.Hex()).To(Equal(payment1.ID.Hex()))
			})

			It("can list all registered payments", func() {
				// List all registered payments and make sure no error has been returned.
				res, err := request.Get(baseUrl + payments.BasePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				// Decode the response's body and make sure that both payments have been returned.
				result := make([]models.Payment, 0)
				err = res.ToJSON(&result)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(ContainElement(MatchFields(IgnoreExtras, Fields{
					paymentIDFieldName: Equal(payment1.ID),
				})))
				Expect(result).To(ContainElement(MatchFields(IgnoreExtras, Fields{
					paymentIDFieldName: Equal(payment2.ID),
				})))
			})

			It("can delete a payment by its ID and does not further list it", func() {
				// Delete the first payment and make sure no error has been returned.
				res, err := request.Delete(baseUrl + payments.BasePath + "/" + payment1.ID.Hex())
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusNoContent))

				// Make sure that the first payment can no longer be retrieved by its ID.
				res, err = request.Get(baseUrl + payments.BasePath + "/" + payment1.ID.Hex())
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusNotFound))

				// Make sure that the second payment can still be retrieved by its ID.
				res, err = request.Get(baseUrl + payments.BasePath + "/" + payment2.ID.Hex())
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))

				// Make sure that the first payment is no longer listed, but that the second one is.
				res, err = request.Get(baseUrl + payments.BasePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				result := make([]models.Payment, 0)
				err = res.ToJSON(&result)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(ContainElement(MatchFields(IgnoreExtras, Fields{
					paymentIDFieldName: Equal(payment1.ID),
				})))
				Expect(result).To(ContainElement(MatchFields(IgnoreExtras, Fields{
					paymentIDFieldName: Equal(payment2.ID),
				})))
			})

			It("can update a payment by its ID", func() {
				// Grab the original ID of the first payment so that we can check that IDs cannot be overwritten.
				originalID := payment1.ID.Hex()
				// Change some of the fields of the first payment, including its ID.
				payment1.Amount = 1200.41
				payment1.ID = payment2.ID
				// Update the first payment and make sure that only the ".amount" field was updated.
				res, err := request.Put(baseUrl+payments.BasePath+"/"+originalID, request.BodyJSON(payment1))
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Response().StatusCode).To(Equal(http.StatusOK))
				result := models.Payment{}
				err = res.ToJSON(&result)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Amount).To(Equal(payment1.Amount))
				Expect(result.ID.Hex()).To(Equal(originalID))
			})
		})
	})
})
