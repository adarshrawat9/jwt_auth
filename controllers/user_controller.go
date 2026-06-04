package controllers

import (
	"context"
	"fmt"
	"jwt-auth/database"
	"jwt-auth/helper"
	"jwt-auth/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    

	"golang.org/x/crypto/bcrypt"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client , "user")
var validate = validator.New()

func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		log.Fatal(err)
	}
	return string(bytes)

}

func VerifyPassword(userPassword string, providedPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil{
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg

}

func Signin()gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user model.User
		var foundUser model.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}

		err := userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"email or passowrd is incorrect"})
			return 

		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return 
		}

		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})

		}
		token , refreshToken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, *foundUser.User_Type, *&foundUser.User_Id)
		helper.UpdateAllToken(token, refreshToken, foundUser.User_Id)
		userCollection.FindOne(ctx, bson.M{"userr_id": foundUser.User_Id}).Decode(&foundUser)

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
		c.JSON(http.StatusOK, foundUser)


	}
}

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

		password := HashPassword(*user.Password)
		user.Password = &password
		count, err = userCollection.CountDocuments(ctx , bson.M{"phone":user.Phone_No})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone count"})
			return 
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email or phone number already exists"})
		}

		user.Created_At, _ = time.Parse(time.RFC3339 , time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339 , time.Now().Format(time.RFC3339))
		user.ID = bson.NewObjectID()
		user.User_Id =  user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email , *user.First_Name, *user.Last_Name,*user.User_Type, user.User_Id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx , user)
		if insertErr != nil{
			msg := fmt.Sprintf("the user item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error" : msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)

		
	}
}

func GetUsers()gin.HandlerFunc{
	return func (c *gin.Context){
		if err := helper.CheckUserType(c, "ADMIN"); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1{
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1{
			page = 1
		}

		startIndex := (page - 1)* recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}

		groupStage := bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: "null"},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			},
		}}

		projectStage := bson.D{{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{
					Key:   "$slice",
					Value: []interface{}{"$data", startIndex, recordPerPage},
				}}},
			},
		}}
	

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured whilelisting user items"})
			return 

		}
		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

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

