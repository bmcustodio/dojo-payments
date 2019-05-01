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

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

// APIServer serves APIs such as the Payments API.
type APIServer struct {
	// echo is the instance of Echo that powers the API server.
	echo *echo.Echo
}

// NewAPIServer returns a new instance of the API server
func NewAPIServer() *APIServer {
	// Create a new instance of the API server.
	s := &APIServer{
		echo: echo.New(),
	}
	// Register the root handler.
	s.echo.Add(http.MethodGet, "/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "")
	})
	// Disable Echo's banner.
	s.echo.HideBanner = true
	// Disable Echo's initial message.
	s.echo.HidePort = true
	// Activate logging of HTTP requests.
	s.echo.Use(middleware.Logger())
	// Assign an ID to each HTTP request.
	s.echo.Use(middleware.RequestID())
	// Return the instance of the API server to the caller.
	return s
}

// Run runs the API server at the specified address.
func (srv *APIServer) Run(bindAddress string) error {
	log.Infof("starting the api server at %s", bindAddress)
	return srv.echo.Start(bindAddress)
}
