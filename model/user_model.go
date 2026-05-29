package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID    `bson:"_id"`
	First_Name    *string               `json:"first_name" validate: "required, min= 2, max = 50"`
	Last_Name     *string               `json:"last_name" validate: "required, min= 2, max = 50"`
	Email         *string               `json:"email" validate: "email, required"`
	Password      *string			    `json:"password" validate: "required, min = 8"`	
	Phone_No      *string				`json:"phone_no" validate:"required"`			
	Token         *string				`json:"token"`
	User_Type     *string				`json:"user_type" validate: "required, eq=ADMIN|eq=USER"`
	Refresh_Token *string				`json:"refresh_token"`
	Created_At    time.Time 			`json:"created_at"`
	Updated_At    time.Time				`json:"updated_at"`
	User_Id       time.Time				`json:"user_id"`
}