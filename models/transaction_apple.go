package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionApple struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AppleAppId            int64              `bson:"appleAppId" json:"appleAppId"`
	OriginalTransactionId string             `bson:"originalTransactionId" json:"originalTransactionId"`
	Status                int32              `bson:"status" json:"status"`
	StatusText            string             `bson:"statusText" json:"statusText"`
	BundleId              string             `bson:"bundleId,omitempty" json:"bundleId,omitempty"`
	Environment           Environment        `bson:"environment,omitempty" json:"environment,omitempty"`
	UpdatedAt             time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (t *TransactionApple) CollectionName() string {
	return "transactionApple"
}
