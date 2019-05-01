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
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bmcstdio/dojo-payments/pkg/constants"
	"github.com/bmcstdio/dojo-payments/pkg/db/models"
)

// PaymentsDatabase contains methods used to perform CRUD operations on payments.
type PaymentsDatabase interface {
	// CreatePayment creates the provided payment.
	CreatePayment(models.Payment) (models.Payment, error)
}

// mongodbPaymentsDatabase is an implementation of PaymentsDatabase powered by MongoDB.
type mongodbPaymentsDatabase struct {
	// c is the MongoDB collection to use for storing payments.
	c *mongo.Collection
}

// CreatePayment creates the provided payment.
func (db *mongodbPaymentsDatabase) CreatePayment(p models.Payment) (models.Payment, error) {
	// Grab the current timestamp and set the creation/modification date.
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	// Create the payment.
	ctx, fn := context.WithTimeout(context.Background(), constants.MongoDBOperationTimeout)
	defer fn()
	r, err := db.c.InsertOne(ctx, p)
	if err != nil {
		return models.Payment{}, fmt.Errorf("failed to create payment: %v", err)
	}
	// Return the full payment back to the caller.
	p.ID = r.InsertedID.(primitive.ObjectID)
	return p, nil
}
