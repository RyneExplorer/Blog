package api

import (
	"blog/internal/api/article"
	"blog/internal/api/auth"
	"blog/internal/api/category"
	"blog/internal/api/comment"
	"blog/internal/api/super"
	"blog/internal/api/user"
	"blog/internal/middleware"
	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

// Router 路由
type Router struct {
	userCtrl     *user.UserController
	authCtrl     *auth.AuthController
	articleCtrl  *article.ArticleController
	commentCtrl  *comment.CommentController
	categoryCtrl *category.CategoryController
	superCtrl    *super.ReviewController
}

// NewRouter 创建路由
func NewRouter(
	userService service.UserService,
	authService service.AuthService,
	articleService service.ArticleService,
	commentService service.CommentService,
	categoryService service.CategoryService,
	reviewService service.ReviewService,
) *Router {
	return &Router{
		userCtrl:     user.NewUserController(userService, articleService),
		authCtrl:     auth.NewAuthController(authService, userService),
		articleCtrl:  article.NewArticleController(articleService, commentService),
		commentCtrl:  comment.NewCommentController(commentService),
		categoryCtrl: category.NewCategoryController(categoryService),
		superCtrl:    super.NewReviewController(reviewService),
	}
}

// Setup 设置路由
func (r *Router) Setup(engine *gin.Engine) {
	// 全局中间件
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// 健康检查
	engine.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "blog API is running",
		})
	})

	// API 路由组
	v1 := engine.Group("/api")
	{
		// 认证路由
		r.authCtrl.RegisterRoutes(v1)
		// 用户路由
		r.userCtrl.RegisterRoutes(v1)
		// 文章路由
		r.articleCtrl.RegisterRoutes(v1)
		// 评论路由
		r.commentCtrl.RegisterRoutes(v1)
		// 分类路由
		r.categoryCtrl.RegisterRoutes(v1)
		// 管理员审核路由
		r.superCtrl.RegisterRoutes(v1)
	}
}
