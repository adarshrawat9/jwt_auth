package controllers

import (
	"jwt-auth/database"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client , "user")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signin()

func Signup()

func GetUsers()

func GetUser()

