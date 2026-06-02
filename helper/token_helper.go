package helper

import (
	"jwt-auth/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v5"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)


type SignedDetail struct{
	Email            string
	First_Name       string     
	Last_Name        string    
	Uid              string
	User_type        string   
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
			ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(168)),
		},
	}

	refreshClaims :=  &SignedDetail{
		RegisteredClaims : jwt.RegisteredClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(168)),
		},
	}

	token , err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRETE_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString([]byte(SECRETE_KEY))

	if err != nil{
		log.Fatal(err)
		return
	}

	return token, refreshToken, err
}