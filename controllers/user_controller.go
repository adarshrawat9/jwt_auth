package controllers

import (
	"context"
	"jwt-auth/database"
	"jwt-auth/helper"
	"jwt-auth/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client , "user")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signin()

func Signup()gin.HandlerFunc{
	return func (c *gin.Context)  {
		var ctx , cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user model.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error ": validationErr.Error()})
			return 

		}
		count, err := userCollection.CountDocuments(ctx , bson.M{"email":user.Email})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}

		count, err = userCollection.CountDocuments(ctx , bson.M{"phone":user.Phone_No})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone count"})
			return 
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email or phone number already exists"})
		}
		
	}
}

func GetUsers()

func GetUser()(gin.HandlerFunc){
	return func (c *gin.Context)  {

		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil{
			c.JSON(http.StatusBadRequest , gin.H{"Error ":  err.Error()})
			return 
		}
		var ctx , cancel = context.WithTimeout(context.Background() , 100*time.Second)

		var user model.User
		err := userCollection.FindOne(ctx , bson.M{"user_id" : userId}).Decode(&user)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return 
		}
		c.JSON(http.StatusOK, user)
		
	}
}

