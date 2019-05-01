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

package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/bmcstdio/dojo-payments/pkg/constants"
	"github.com/bmcstdio/dojo-payments/pkg/db"
	"github.com/bmcstdio/dojo-payments/pkg/server/apis/payments"
)

const (
	// DatabaseStatusOffline indicates that the database cannot be reached.
	DatabaseStatusOffline = "OFFLINE"
	// DatabaseStatusOnline indicates that the database can be reached.
	DatabaseStatusOnline = "ONLINE"
)

// APIServerRootResponse represents a response returned by the root handler.
type APIServerRootResponse struct {
	// DatabaseStatus is the current status of the database.
	DatabaseStatus string `json:"database_status"`
	// Timestamp is the current timestamp.
	Timestamp time.Time `json:"time"`
}

// APIServer serves APIs such as the Payments API.
type APIServer struct {
	// echo is the instance of Echo that powers the API server.
	echo *echo.Echo
}

// NewAPIServer returns a new instance of the API server that uses the specified database for storage.
func NewAPIServer(database db.Database) *APIServer {
	// Create a new instance of the API server.
	s := &APIServer{
		echo: echo.New(),
	}
	// Register the root handler.
	s.echo.Add(http.MethodGet, "/", func(ctx echo.Context) error {
		var (
			status string
		)
		if ctx.Get(constants.DatabaseContextKey).(db.Database).IsOnline() {
			status = DatabaseStatusOnline
		} else {
			status = DatabaseStatusOffline
		}
		return ctx.JSON(http.StatusOK, APIServerRootResponse{
			DatabaseStatus: status,
			Timestamp:      time.Now(),
		})
	})
	// Disable Echo's banner.
	s.echo.HideBanner = true
	// Disable Echo's initial message.
	s.echo.HidePort = true
	// Activate logging of HTTP requests.
	s.echo.Use(middleware.Logger())
	// Assign an ID to each HTTP request.
	s.echo.Use(middleware.RequestID())
	// Add the database to the context so that HTTP handlers can use it to actually store data.
	s.echo.Use(func(fn echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(constants.DatabaseContextKey, database)
			return fn(ctx)
		}
	})
	// Register the Payments API.
	payments.Register(s.echo)
	// Return the instance of the API server to the caller.
	return s
}

// Run runs the API server at the specified address.
func (srv *APIServer) Run(bindAddress string) error {
	log.Infof("starting the api server at %s", bindAddress)
	return srv.echo.Start(bindAddress)
}
