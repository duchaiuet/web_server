package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Model struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	DeletedAt *time.Time         `json:"deleted_at" bson:"deleted_at"`
	CreatedAt *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at"`
}

type User struct {
	Model        `bson:"inline"`
	UserName     string `json:"user_name" bson:"user_name"`
	Password     string `json:"password" bson:"password"`
	EmailAddress string `json:"email_address" bson:"email_address"`
	FirstName    string `json:"first_name" bson:"first_name"`
	LastName     string `json:"last_name" bson:"last_name"`
	Status       string `json:"status" bson:"status"`
	Role         int    `json:"role" bson:"role"`
}
