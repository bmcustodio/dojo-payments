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

package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/bmcstdio/dojo-payments/pkg/db"
	"github.com/bmcstdio/dojo-payments/pkg/server"
)

var (
	// bindAddr is the "host:port" combination at which to serve the API server.
	bindAddr string
	// mongodbDatabase is the name of the MongoDB database to use for storage.
	mongodbDatabase string
	// mongodbUrl is the URL at which MongoDB can be reached.
	mongodbURL string
)

func init() {
	flag.StringVar(&bindAddr, "bind-addr", ":8080", `the "host:port" combination at which to serve the api server`)
	flag.StringVar(&mongodbDatabase, "mongodb-database", "dojo-payments", "the name of the mongodb database to use for storage")
	flag.StringVar(&mongodbURL, "mongodb-url", "mongodb://localhost:27017", "the url at which mongodb can be reached")
}

func main() {
	// Parse the provided command-line flags.
	flag.Parse()

	// Initialize the the database.
	database, err := db.NewMongoDDatabase(mongodbURL, mongodbDatabase)
	if err != nil {
		log.Fatalf("failed to initialize the database: %v", err)
	}

	// Initialize and run the API server using this database for storage.
	srv := server.NewAPIServer(database)
	if err := srv.Run(bindAddr); err != nil {
		log.Fatalf("failed to run the api server: %v", err)
	}
}
