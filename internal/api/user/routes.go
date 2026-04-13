package user

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户路由
func (ctrl *UserController) RegisterRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	userGroup.Use(middleware.Auth())
	{
		userGroup.GET("/articles", ctrl.ListMyArticles)
		userGroup.GET("/profile", ctrl.GetProfile)
		userGroup.PUT("/profile", ctrl.UpdateProfile)
		userGroup.POST("/password", ctrl.ChangePassword)
		userGroup.GET("/list", ctrl.ListUsers)
	}
}
