package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty" validate:"required"`
	Phone     string             `bson:"phone,omitempty" validate:"required"`
	Email     string             `bson:"email,omitempty" validate:"required"`
	Password  string             `bson:"password,omitempty" validate:"required"`
	Address   Address            `bson:"address,omitempty"`
	Points    Points             `bson:"points,omitempty" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" validate:"required"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" validate:"required"`
}

type Address struct {
	Phone       string `bson:"phone,omitempty" validate:"required"`
	Province    string `bson:"province,omitempty" validate:"required"`
	SubDistrict string `bson:"sub_district,omitempty" validate:"required"`
	City        string `bson:"city,omitempty" validate:"required"`
	FullAddress string `bson:"fullAddress,omitempty" validate:"required"`
	PostalCode  string `bson:"postalCode,omitempty" validate:"required"`
}

type Points struct {
	TotalPoints  float64      `bson:"totalPoints" validate:"required"`
	PointHistory PointHistory `bson:"pointOutHistory,omitempty"`
}

type PointHistory struct {
	ProductName string    `bson:"productName"`
	PointOut    float64   `bson:"pointOut" validate:"required"`
	PointIn     float64   `bson:"pointIn" validate:"required"`
	CreatedAt   time.Time `bson:"createdAt" validate:"required"`
}
