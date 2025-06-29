package router

import (
	"github.com/gin-gonic/gin"
	"kitexTwo/client/handler/api"
)

func UserModel(r *gin.RouterGroup) {
	userModel := r.Group("/user")
	{
		userModel.POST("/createUser", api.CreateUser)
		userModel.POST("/loginUser", api.LoginUser)
	}
}
