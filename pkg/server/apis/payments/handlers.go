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

package payments

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/bmcstdio/dojo-payments/pkg/constants"
	"github.com/bmcstdio/dojo-payments/pkg/db"
	"github.com/bmcstdio/dojo-payments/pkg/db/models"
)

const (
	// BasePath is the base path of the Payments API.
	BasePath = "/payments"
)

// Register registers the handlers for the Payments API to the provided Echo instance.
func Register(echo *echo.Echo) {
	echo.Add(http.MethodPost, BasePath, createPayment)
	echo.Add(http.MethodDelete, BasePath+"/:id", deletePayment)
	echo.Add(http.MethodGet, BasePath+"/:id", getPayment)
	echo.Add(http.MethodGet, BasePath, listPayments)
	echo.Add(http.MethodPut, BasePath+"/:id", updatePayment)
}

// createPayment creates a payment.
func createPayment(ctx echo.Context) error {
	var (
		err error
		p   models.Payment
	)
	// Grab the Payment object details from the request's body and validate it.
	if err := ctx.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := p.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	p, err = ctx.Get(constants.DatabaseContextKey).(db.Database).Payments().CreatePayment(p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, p)
}

// deletePayment deletes a payment by ID.
func deletePayment(ctx echo.Context) error {
	d, err := ctx.Get(constants.DatabaseContextKey).(db.Database).Payments().DeletePayment(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if !d {
		return echo.NewHTTPError(http.StatusNotFound, "payment not found")
	}
	return ctx.String(http.StatusNoContent, "")
}

// getPayment gets a payment by ID.
func getPayment(ctx echo.Context) error {
	p, err := ctx.Get(constants.DatabaseContextKey).(db.Database).Payments().GetPayment(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if p == (models.Payment{}) {
		return echo.NewHTTPError(http.StatusNotFound, "payment not found")
	}
	return ctx.JSON(http.StatusOK, p)
}

// listPayments lists payments.
func listPayments(ctx echo.Context) error {
	r, err := ctx.Get(constants.DatabaseContextKey).(db.Database).Payments().ListPayments()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, r)
}

// updatePayment updates a payment by ID.
func updatePayment(ctx echo.Context) error {
	var (
		err error
		p   models.Payment
		r   models.Payment
	)
	// Grab the Payment object details from the request's body and validate it.
	if err := ctx.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := p.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	r, err = ctx.Get(constants.DatabaseContextKey).(db.Database).Payments().UpdatePayment(ctx.Param("id"), p)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if p == (models.Payment{}) {
		return echo.NewHTTPError(http.StatusNotFound, "payment not found")
	}
	return ctx.JSON(http.StatusOK, r)
}
