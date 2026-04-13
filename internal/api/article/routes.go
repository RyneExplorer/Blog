package article

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册 /api/articles 路由
func (ctrl *ArticleController) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/articles")
	{
		g.POST("", middleware.Auth(), ctrl.Create)
		g.GET("", ctrl.List)
		g.GET("/:id/comments", ctrl.ListComments)
		g.GET("/:id", middleware.OptionalAuth(), ctrl.Detail)
		g.PUT("/:id", middleware.Auth(), ctrl.Update)
		g.POST("/:id/publish", middleware.Auth(), ctrl.Publish)
		g.DELETE("/:id", middleware.Auth(), ctrl.Delete)
		g.POST("/:id/view", ctrl.IncrView)
		g.POST("/:id/like", middleware.Auth(), ctrl.Like)
		g.DELETE("/:id/like", middleware.Auth(), ctrl.Unlike)
		g.POST("/:id/favorite", middleware.Auth(), ctrl.Favorite)
		g.DELETE("/:id/favorite", middleware.Auth(), ctrl.Unfavorite)
	}
}
