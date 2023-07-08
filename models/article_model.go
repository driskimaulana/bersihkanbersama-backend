package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Article struct {
	Id        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	Writer    string             `bson:"writer"`
	Summary   string             `bson:"summary"`
	Content   string             `bson:"content"`
	Cover     string             `bson:"cover"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
