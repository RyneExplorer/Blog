package article

import (
	"blog/internal/api/comment"
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 /api/articles 路由
func (ctrl *ArticleController) RegisterRoutes(r *gin.RouterGroup, commentCtrl *comment.CommentController) {
	g := r.Group("/articles")
	{
		g.GET("", ctrl.List)
		g.GET("/:id/comments", commentCtrl.ListByArticle)
		g.GET("/:id", ctrl.Detail)
		g.POST("/:id/view", ctrl.IncrView)
		g.POST("/:id/like", middleware.Auth(), ctrl.Like)
		g.DELETE("/:id/like", middleware.Auth(), ctrl.Unlike)
		g.POST("/:id/favorite", middleware.Auth(), ctrl.Favorite)
		g.DELETE("/:id/favorite", middleware.Auth(), ctrl.Unfavorite)
	}
}
