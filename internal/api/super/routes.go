package super

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册管理模块路由
func (ctrl *ReviewController) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/super")
	g.Use(middleware.Auth(), middleware.RequireRole(0))
	{
		// 获取审核文章列表
		g.GET("/articles", ctrl.List)
		// 获取审核文章详情
		g.GET("/articles/:id", ctrl.Detail)
		// 获取用户列表
		g.GET("/userlist", ctrl.ListUsers)
		// 审核通过文章
		g.POST("/articles/:id/approve", ctrl.Approve)
		// 驳回文章
		g.POST("/articles/:id/reject", ctrl.Reject)
		// 封禁文章
		g.POST("/articles/:id/ban", ctrl.Ban)
		// 更新文章分类
		g.PUT("/articles/:id/category", ctrl.UpdateCategory)
	}
}
