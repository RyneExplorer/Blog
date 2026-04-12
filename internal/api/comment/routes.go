package comment

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 /api/comments 相关路由
func (ctrl *CommentController) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/comments")
	g.Use(middleware.Auth())
	{
		g.POST("", ctrl.Create)
		g.POST("/reply", ctrl.Reply)
		g.DELETE("/:id", ctrl.Delete)
		g.POST("/:id/like", ctrl.Like)
		g.DELETE("/:id/like", ctrl.Unlike)
	}
}
