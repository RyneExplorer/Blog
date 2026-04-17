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
		// 发布评论
		g.POST("", ctrl.Create)
		// 回复评论
		g.POST("/replies", ctrl.Reply)
		// 点赞评论
		g.POST("/:id/like", ctrl.Like)
		// 取消点赞评论
		g.POST("/:id/unlike", ctrl.Unlike)
		// 删除评论
		g.DELETE("/:id", ctrl.Delete)
	}
}
