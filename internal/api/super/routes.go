package super

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 /api/super 路由
func (ctrl *ReviewController) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/super")
	g.Use(middleware.Auth())
	{
		g.GET("/articles", ctrl.List)
		g.GET("/articles/:id", ctrl.Detail)
		g.POST("/articles/:id/approve", ctrl.Approve)
		g.POST("/articles/:id/reject", ctrl.Reject)
		g.POST("/articles/:id/ban", ctrl.Ban)
		g.PUT("/articles/:id/category", ctrl.UpdateCategory)
	}
}
