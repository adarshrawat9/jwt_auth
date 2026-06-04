package helper

import (
	"context"
	"fmt"
	"jwt-auth/database"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)
	


type SignedDetail struct{
	Email            string
	First_Name       string     
	Last_Name        string    
	Uid              string
	User_type        string  
	ExpiresAt        time.Time 
	jwt.RegisteredClaims
}


var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRETE_KEY string = os.Getenv("SECRET_KEY")


func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string)(signedToken string, signedRefreshToken string , err error){
	claims := &SignedDetail{
		Email: email,
		First_Name: firstName,
		Last_Name: lastName,
		User_type: userType,
		Uid: uid,
		RegisteredClaims : jwt.RegisteredClaims{
			ExpiresAt:  jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(168))),
		},
	}

	refreshClaims :=  &SignedDetail{
		RegisteredClaims : jwt.RegisteredClaims{
			ExpiresAt:  jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(168))),
		},
	}

	token , err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRETE_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRETE_KEY))

	if err != nil{
		log.Fatal(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllToken(signedToken string, signedRefreshToken string, userId string){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj bson.D

		updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
		updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
		updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updated_at})

	


	upsert := true
	filter := bson.M{"user_id":userId}
	opt := options.UpdateOne().SetUpsert(upsert)


	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
		{Key: "$set", Value: updateObj},
		},
		opt,

	)
	defer cancel()

	if err != nil{
		log.Fatal(err)
		return
	}
	return
}


func ValidateToken(signedToken string)(claims *SignedDetail, msg string){
	token , err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetail{},
		func (token *jwt.Token)(interface{}, error)  {
			return []byte(SECRETE_KEY), nil
			
		},

	)
	if err != nil{
		msg = err.Error()
		return
	}

	claims , ok := token.Claims.(*SignedDetail)
	if !ok{
		msg = fmt.Sprintf("the token is invalid")
		return
	}

	if claims.ExpiresAt.Before(time.Now()){
		msg = fmt.Sprintf("token is expired")
		return
	}

	return claims , msg

}