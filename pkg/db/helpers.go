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
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// deletedAtFieldName is the name of the field that holds the deletion date of a given record.
	deletedAtFieldName = "deleted_at"
	// idFieldName is the name of the field that holds the ID of a given record.
	idFieldName = "_id"
)

const (
	// eqOp represents the "$eq" operator.
	eqOp = "$eq"
	// setOp represents the "$set" operator.
	setOp = "$set"
)

// existing is a helper method that allows for selecting existing (i.e. not deleted) objects.
func existing() primitive.M {
	return primitive.M{
		deletedAtFieldName: primitive.M{
			eqOp: nil,
		},
	}
}

// existingByID is a helper method that allows for selecting an existing (i.e. not deleted) object by its ID.
func existingByID(id primitive.ObjectID) primitive.M {
	return primitive.M{
		idFieldName: id,
		deletedAtFieldName: primitive.M{
			eqOp: nil,
		},
	}
}

// markDeleted is a helper method that allows for marking an object as deleted.
func markDeleted(time time.Time) primitive.M {
	return primitive.M{
		setOp: primitive.M{
			deletedAtFieldName: time,
		},
	}
}
