package routes

import (
	controller "jwt-auth/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.POST("user/signup" , controller.Signin())
	incomingRoutes.POST("user/signin" , controller.Signup())
}