package user

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户模块路由
func (ctrl *UserController) RegisterRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	userGroup.Use(middleware.Auth())
	{
		// 获取个人资料
		userGroup.GET("/profile", ctrl.GetProfile)
		// 更新个人资料
		userGroup.PUT("/profile", ctrl.UpdateProfile)
		// 上传头像
		userGroup.POST("/avatar", ctrl.UploadAvatar)
		// 修改密码
		userGroup.POST("/password", ctrl.ChangePassword)
	}
}
