package article

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文章模块路由
func (ctrl *ArticleController) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/articles")
	{
		// 创建文章
		g.POST("", middleware.Auth(), ctrl.Create)

		// 上传文章封面
		g.POST("/cover_image", middleware.Auth(), ctrl.UploadCover)

		// 上传正文图片
		g.POST("/content_image", middleware.Auth(), ctrl.UploadContentImage)

		// 获取文章列表
		g.GET("", ctrl.List)

		// 获取我的收藏列表
		g.GET("/favorites", middleware.Auth(), ctrl.ListFavorites)

		// 获取我的文章列表
		g.GET("/mine", middleware.Auth(), ctrl.ListMyArticles)

		// 获取我的文章详情
		g.GET("/mine/:id", middleware.Auth(), ctrl.MyDetail)

		// 获取文章评论列表
		g.GET("/:id/comments", ctrl.ListComments)

		// 获取首页文章详情
		g.GET("/:id", ctrl.Detail)

		// 更新文章
		g.PUT("/:id", middleware.Auth(), ctrl.Update)

		// 发布文章
		g.POST("/:id/publish", middleware.Auth(), ctrl.Publish)

		// 删除文章
		g.DELETE("/:id", middleware.Auth(), ctrl.Delete)

		// 增加浏览量
		g.POST("/:id/view", ctrl.IncrView)

		// 点赞文章
		g.POST("/:id/like", middleware.Auth(), ctrl.Like)

		// 取消点赞
		g.POST("/:id/unlike", middleware.Auth(), ctrl.Unlike)

		// 收藏文章
		g.POST("/:id/favorite", middleware.Auth(), ctrl.Favorite)

		// 取消收藏
		g.POST("/:id/unfavorite", middleware.Auth(), ctrl.Unfavorite)
	}
}
