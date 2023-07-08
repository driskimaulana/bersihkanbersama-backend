package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Donation struct {
	Id            primitive.ObjectID `bson:"_id" validate:"required"`
	UserId        primitive.ObjectID `bson:"userId" validate:"required"`
	ActivityId    primitive.ObjectID `bson:"activityId" validate:"required"`
	Items         []DonationItem     `bson:"items" validate:"required"`
	TotalDonation float64            `bson:"totalDonation" validate:"required"`
	PaymentId     string             `bson:"paymentId" validate:"required"`
	Status        string             `bson:"status" validate:"required"`
	IsAnonim      bool               `bson:"isAnonim" validate:"required"`
	UserName      string             `bson:"userName" validate:"required"`
}

type DonationItem struct {
	Name  string  `bson:"name" validate:"required"`
	Count int     `bson:"count" validate:"required"`
	Price float64 `bson:"price" validate:"required"`
}
