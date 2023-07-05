package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Email     string             `json:"email,omitempty" validate:"required"`
	Password  string             `json:"password,omitempty" validate:"required"`
	CreatedAt primitive.DateTime `json:"createdAt,omitempty" validate:"required"`
	UpdatedAt primitive.DateTime `json:"updatedAt,omitempty" validate:"required"`
}
