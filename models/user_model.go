package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id              primitive.ObjectID   `bson:"_id,omitempty"`
	Name            string               `bson:"name,omitempty" validate:"required"`
	Phone           string               `bson:"phone,omitempty" validate:"required"`
	Role            string               `bson:"role,omitempty" validate:"required"`
	Email           string               `bson:"email,omitempty" validate:"required"`
	Password        string               `bson:"password,omitempty" validate:"required"`
	Address         Address              `bson:"address,omitempty"`
	Points          Points               `bson:"points,omitempty" validate:"required"`
	Activity        []primitive.ObjectID `bson:"activity"`
	DonationHistory []primitive.ObjectID `bson:"donationHistory"`
	CreatedAt       time.Time            `bson:"createdAt,omitempty" validate:"required"`
	UpdatedAt       time.Time            `bson:"updatedAt,omitempty" validate:"required"`
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
	TotalPoints  int            `bson:"totalPoints" validate:"required"`
	PointHistory []PointHistory `bson:"pointHistory,omitempty"`
}

type PointHistory struct {
	Title       string    `bson:"title"`
	ProductName string    `bson:"productName"`
	PointOut    int       `bson:"pointOut" validate:"required"`
	PointIn     int       `bson:"pointIn" validate:"required"`
	CreatedAt   time.Time `bson:"createdAt" validate:"required"`
}
