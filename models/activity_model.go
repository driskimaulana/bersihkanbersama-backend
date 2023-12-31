package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Activity struct {
	Id               primitive.ObjectID `bson:"_id"`
	OrganizationId   primitive.ObjectID `bson:"organizationId" validate:"required"`
	Title            string             `bson:"title" validate:"required"`
	Description      string             `bson:"description" validate:"required"`
	EventDate        string             `bson:"eventDate" validate:"required"`
	Location         Location           `bson:"location" validate:"required"`
	CoverImage       string             `bson:"coverImage"`
	Volunteer        Volunteer          `bson:"volunteer"`
	Status           string             `bson:"status" validate:"required"`
	Rewards          Rewards            `bson:"rewards" validate:"required"`
	DonationActivity DonationActivity   `bson:"donationActivity" validate:"required"`
	CreatedAt        time.Time          `bson:"createdAt" validate:"required"`
	UpdatedAt        time.Time          `bson:"updatedAt" validate:"required"`
}

type Location struct {
	City        string `bson:"city"`
	FullAddress string `bson:"fullAddress"`
	Latitude    string `bson:"latitude"`
	Longitude   string `bson:"longitude"`
}

type Volunteer struct {
	Count          int              `bson:"count" validate:"required"`
	UserRegistered []UserRegistered `bson:"userRegistered"`
	Teams          []Team           `bson:"teams"`
}

type UserRegistered struct {
	Id    primitive.ObjectID `bson:"_id" validate:"required"`
	Name  string             `bson:"name" validate:"required"`
	Phone string             `bson:"phone" validate:"required"`
}

type Team struct {
	Name         string           `bson:"name" validate:"required"`
	Members      []UserRegistered `bson:"members" validate:"required"`
	TrashResults float64          `bson:"trashResult" validate:"required"`
}

type Rewards struct {
	Participation int `bson:"participation" validate:"required"`
	First         int `bson:"first" validate:"required"`
	Second        int `bson:"second" validate:"required"`
	Third         int `bson:"third" validate:"required"`
}

type DonationActivity struct {
	TotalDonation   float64           `bson:"totalDonation" validate:"required"`
	DonationHistory []DonationSummary `bson:"donationHistory"`
	//ReceivedDonation []DonationItem `bson:"receivedDonation"`
}

type DonationSummary struct {
	DonationId    primitive.ObjectID `bson:"donationId"`
	UserName      string             `bson:"userName"`
	TotalDonation float64            `bson:"totalDonation"`
}
