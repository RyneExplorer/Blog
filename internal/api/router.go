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
	"blog/pkg/upload"

	"github.com/gin-gonic/gin"
)

// Router 应用路由集合
type Router struct {
	userCtrl     *user.UserController
	authCtrl     *auth.AuthController
	articleCtrl  *article.ArticleController
	commentCtrl  *comment.CommentController
	categoryCtrl *category.CategoryController
	superCtrl    *super.ReviewController
}

// NewRouter 创建路由对象
func NewRouter(
	userService service.UserService,
	authService service.AuthService,
	articleService service.ArticleService,
	commentService service.CommentService,
	categoryService service.CategoryService,
	reviewService service.ReviewService,
) *Router {
	return &Router{
		userCtrl:     user.NewUserController(userService),
		authCtrl:     auth.NewAuthController(authService, userService),
		articleCtrl:  article.NewArticleController(articleService, commentService),
		commentCtrl:  comment.NewCommentController(commentService),
		categoryCtrl: category.NewCategoryController(categoryService),
		superCtrl:    super.NewReviewController(reviewService, userService),
	}
}

// Setup 注册所有路由
func (r *Router) Setup(engine *gin.Engine) {
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())
	if err := upload.EnsureBaseDir(); err != nil {
		panic(err)
	}
	engine.Static("/uploads", "./uploads")

	engine.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "blog API is running",
		})
	})

	v1 := engine.Group("/api")
	{
		r.authCtrl.RegisterRoutes(v1)
		r.userCtrl.RegisterRoutes(v1)
		r.articleCtrl.RegisterRoutes(v1)
		r.commentCtrl.RegisterRoutes(v1)
		r.categoryCtrl.RegisterRoutes(v1)
		r.superCtrl.RegisterRoutes(v1)
	}
}
