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

package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/bmcstdio/dojo-payments/pkg/constants"
)

// Database represents the database where data will be stored.
type Database interface {
	// IsOnline returns a value indicating whether the database is online.
	IsOnline() bool
	// Payments allows for accessing methods used to perform CRUD operations on payments.
	Payments() PaymentsDatabase
}

// mongodbDatabase is an implementation of Database powered by MongoDB.
type mongodbDatabase struct {
	// db is the actual MongoDB database in which to store data.
	db *mongo.Database
}

// NewMongoDDatabase returns a new instance of Database powered by MongoDB.
func NewMongoDDatabase(mongodbURL, databaseName string) (Database, error) {
	ctx, fn := context.WithTimeout(context.Background(), constants.MongoDBOperationTimeout)
	defer fn()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %v", err)
	}
	return &mongodbDatabase{
		db: c.Database(databaseName),
	}, nil
}

// IsOnline returns a value indicating whether the database is online.
func (m *mongodbDatabase) IsOnline() bool {
	ctx, fn := context.WithTimeout(context.Background(), constants.MongoDBOperationTimeout)
	defer fn()
	err := m.db.Client().Ping(ctx, readpref.Primary())
	return err == nil
}

// Payments allows for accessing methods used to perform CRUD operations on payments.
func (m *mongodbDatabase) Payments() PaymentsDatabase {
	return &mongodbPaymentsDatabase{
		c: m.db.Collection("payments"),
	}
}
