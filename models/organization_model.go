package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Organization struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" validate:"required"`
	Email       string             `bson:"email" validate:"required"`
	Password    string             `bson:"password" validate:"required"`
	Role        string             `bson:"role,omitempty" validate:"required"`
	Description string             `bson:"description,omitempty" validate:"required"`
	City        string             `bson:"address,omitempty" validate:"required"`
	Logo        string             `bson:"logo,omitempty"`
	Instagram   string             `bson:"instagram,omitempty"`
	Contact     Contact            `bson:"contact,omitempty" validate:"required"`
	Activities  []Activity         `bson:"activities"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" validate:"required"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" validate:"required"`
}

type Contact struct {
	Name  string `bson:"name,omitempty" validate:"required"`
	Phone string `bson:"phone,omitempty" validate:"required"`
}
