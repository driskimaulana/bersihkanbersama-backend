package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Activity struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" validate:"required"`
	OrganizationId primitive.ObjectID `bson:"organizationId,omitempty" validate:"required"`
	Title          string             `bson:"title,omitempty" validate:"required"`
	Description    string             `bson:"description,omitempty" validate:"required"`
	EventDate      string             `bson:"eventDate,omitempty" validate:"required"`
	Location       Location           `bson:"location,omitempty" validate:"required"`
	CoverImage     string             `bson:"coverImage,omitempty" validate:"required"`
	Volunteer      Volunteer          `bson:"volunteer,omitempty"`
	Status         string             `bson:"status,omitempty" validate:"required"`
	Rewards        Rewards            `bson:"rewards" validate:"required"`
	Donation       Donation           `bson:"donation" validate:"required"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" validate:"required"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" validate:"required"`
}

type Location struct {
	Latitude  string `bson:"latitude,omitempty" validate:"required"`
	Longitude string `bson:"longitude,omitempty" validate:"required"`
}

type Volunteer struct {
	Count          int              `bson:"count" validate:"required"`
	UserRegistered []UserRegistered `bson:"userRegistered,omitempty"`
	Teams          Teams            `bson:"teams"`
}

type UserRegistered struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" validate:"required"`
	Name  string             `bson:"name,omitempty" validate:"required"`
	Phone string             `bson:"phone,omitempty" validate:"required"`
}

type Teams struct {
	Name         string           `bson:"name,omitempty" validate:"required"`
	Members      []UserRegistered `bson:"members,omitempty" validate:"required"`
	TrashResults float32          `bson:"trashResult" validate:"required"`
}

type Rewards struct {
	Participation int `bson:"participation" validate:"required"`
	First         int `bson:"first" validate:"required"`
	Second        int `bson:"second" validate:"required"`
	Third         int `bson:"third" validate:"required"`
}

type Donation struct {
	TotalDonation    float64           `bson:"totalDonation" validate:"required"`
	ReceivedDonation []DonationItem    `bson:"receivedDonation"`
	DonationHistory  []DonationHistory `bson:"donationHistory"`
}

type DonationHistory struct {
	Id            primitive.ObjectID `bson:"id,omitempty" validate:"required"`
	UserId        primitive.ObjectID `bson:"userId,omitempty" validate:"required"`
	Items         []DonationItem     `bson:"items" validate:"required"`
	TotalDonation float64            `bson:"totalDonation" validate:"required"`
	PaymentId     string             `bson:"paymentId" validate:"required"`
	IsAnonim      bool               `bson:"isAnonim" validate:"required"`
	UserName      string             `bson:"userName,omitempty" validate:"required"`
}

type DonationItem struct {
	Name  string `bson:"type" validate:"required"`
	Count int    `bson:"count" validate:"required"`
}
