package middleware

import (
	"fmt"
	"jwt-auth/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)


func Authenticate() gin.HandlerFunc{
	return func (c *gin.Context)  {
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("token is not provided")})
			c.Abort()
			return 
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return 
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_Name)
		c.Set("last_name", claims.Last_Name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_type)
		c.Next()


	}
}

