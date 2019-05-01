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
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bmcstdio/dojo-payments/pkg/constants"
	"github.com/bmcstdio/dojo-payments/pkg/db/models"
)

// PaymentsDatabase contains methods used to perform CRUD operations on payments.
type PaymentsDatabase interface {
	// CreatePayment creates the provided payment.
	CreatePayment(models.Payment) (models.Payment, error)
	// DeletePayment deletes the payment with the specified ID.
	DeletePayment(string) (bool, error)
	// GetPayment returns the payment with the specified ID.
	GetPayment(string) (models.Payment, error)
	// ListPayments lists all registered payments.
	ListPayments() ([]models.Payment, error)
	// UpdatePayment updates the payment with the specified ID.
	UpdatePayment(string, models.Payment) (models.Payment, error)
}

// mongodbPaymentsDatabase is an implementation of PaymentsDatabase powered by MongoDB.
type mongodbPaymentsDatabase struct {
	// c is the MongoDB collection to use for storing payments.
	c *mongo.Collection
}

// CreatePayment creates the provided payment.
func (db *mongodbPaymentsDatabase) CreatePayment(p models.Payment) (models.Payment, error) {
	// Grab the current timestamp and set the modification date.
	now := time.Now()
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

// DeletePayment deletes the payment with the specified ID.
func (db *mongodbPaymentsDatabase) DeletePayment(id string) (bool, error) {
	// Grab the current timestamp so we can set the deletion date.
	now := time.Now()
	// Grab the ObjectID that corresponds to the provided ID.
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("%q is not a valid payment ID", id)
	}
	// Try to mark the payment as having been deleted.
	ctx, fn := context.WithTimeout(context.Background(), constants.MongoDBOperationTimeout)
	defer fn()
	r, err := db.c.UpdateOne(ctx, existingByID(objectID), markDeleted(now))
	if err != nil {
		return false, fmt.Errorf("failed to delete payment with id %q: %v", id, err)
	}
	return r.ModifiedCount != 0, nil
}

// GetPayment returns the payment with the provided ID.
func (db *mongodbPaymentsDatabase) GetPayment(id string) (models.Payment, error) {
	// Grab the ObjectID that corresponds to the provided ID.
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Payment{}, fmt.Errorf("%q is not a valid payment ID", id)
	}
	// Try to retrieve the payment with the provided ID, excluding deleted payments.
	ctx, fn := context.WithTimeout(context.Background(), constants.MongoDBOperationTimeout)
	defer fn()
	r := db.c.FindOne(ctx, existingByID(objectID))
	if r.Err() != nil {
		return models.Payment{}, fmt.Errorf("failed to get payment with id %q: %v", id, r.Err())
	}
	// Check whether a payment with the provided ID was found, and return it if it does.
	p := models.Payment{}
	if err := r.Decode(&p); err != nil {
		if err != mongo.ErrNoDocuments {
			// The payment might exist or not, but we've got an unexpected error which we must propagate.
			return models.Payment{}, fmt.Errorf("failed to get payment with id %q: %v", id, err)
		}
		// The payment was not found, so we just return an empty payment (and error).
		return models.Payment{}, nil
	}
	return p, nil
}

// ListPayments lists all registered payments.
func (db *mongodbPaymentsDatabase) ListPayments() ([]models.Payment, error) {
	// Try to retrieve all registered payments, excluding deleted ones.
	ctx, fn := context.WithTimeout(context.Background(), constants.MongoDBOperationTimeout)
	defer fn()
	c, err := db.c.Find(ctx, existing())
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %v", err)
	}
	defer c.Close(ctx)
	// Build the list of payments and return it back to the caller.
	r := make([]models.Payment, 0)
	for c.Next(ctx) {
		p := models.Payment{}
		if err := c.Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to list payments: %v", err)
		}
		r = append(r, p)
	}
	if c.Err() != nil {
		return nil, fmt.Errorf("failed to list payments: %v", err)
	}
	return r, nil
}

// UpdatePayment updates the payment with the specified ID.
func (db *mongodbPaymentsDatabase) UpdatePayment(id string, p models.Payment) (models.Payment, error) {
	// Grab the current timestamp so we can set the modification date.
	now := time.Now()
	// Grab the ObjectID that corresponds to the provided ID.
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Payment{}, fmt.Errorf("%q is not a valid payment ID", id)
	}
	// Force-overwrite the payment's ID so that it is not possibly changed during the update.
	p.ID = objectID
	// Set the payment's modification date.
	p.UpdatedAt = now
	// Try to update the payment with the specified ID, requesting for the new (updated) document to be returned.
	opts := &options.FindOneAndReplaceOptions{}
	opts.SetReturnDocument(options.After)
	ctx, fn := context.WithTimeout(context.Background(), constants.MongoDBOperationTimeout)
	defer fn()
	r := db.c.FindOneAndReplace(ctx, existingByID(objectID), p, opts)
	if r.Err() != nil {
		return models.Payment{}, fmt.Errorf("failed to update payment: %v", r.Err())
	}
	// Check whether a payment with the provided ID was found, and return it if it does.
	res := models.Payment{}
	if err := r.Decode(&res); err != nil {
		if err != mongo.ErrNoDocuments {
			// The payment might exist or not, but we've got an unexpected error which we must propagate.
			return models.Payment{}, fmt.Errorf("failed to update payment: %v", r.Err())
		}
		// The payment was not found, so we just return an empty payment (and error).
		return models.Payment{}, nil
	}
	return res, nil
}
