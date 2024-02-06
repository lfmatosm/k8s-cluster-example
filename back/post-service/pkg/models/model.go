package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostDTO struct {
	Mime  string `json:"mime"`
	Image []byte `json:"image"`
}

type Post struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Mime        string             `json:"mime" bson:"mime"`
	Image       primitive.Binary   `json:"image" bson:"image"`
	LastUpdated time.Time          `json:"lastUpdated" bson:"lastUpdated"`
}
