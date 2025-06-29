package router

import "github.com/gin-gonic/gin"

func Router(r *gin.Engine) {
	routerGroup := r.Group("/kitex")
	{
		apiGroup := routerGroup.Group("/api")
		{
			UserModel(apiGroup)
		}
	}
}
