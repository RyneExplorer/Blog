package article

import (
	"strconv"

	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// ArticleController 文章控制器
type ArticleController struct {
	articleService service.ArticleService
}

// NewArticleController 创建文章控制器
func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

func parseUintParam(c *gin.Context, key string) (uint, bool) {
	s := c.Param(key)
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil || v == 0 {
		return 0, false
	}
	return uint(v), true
}

// List 获取文章列表
func (ctrl *ArticleController) List(c *gin.Context) {
	var q request.ArticleListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确："+err.Error())
		return
	}
	page, err := ctrl.articleService.ListArticles(c.Request.Context(), &q)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, page)
}

// Detail 文章详情
func (ctrl *ArticleController) Detail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	data, err := ctrl.articleService.GetArticleDetail(c.Request.Context(), id)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

// IncrView 浏览量 +1
func (ctrl *ArticleController) IncrView(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	if err := ctrl.articleService.IncrementView(c.Request.Context(), id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Like 点赞文章
func (ctrl *ArticleController) Like(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.articleService.LikeArticle(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Unlike 取消点赞
func (ctrl *ArticleController) Unlike(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.articleService.UnlikeArticle(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Favorite 收藏文章
func (ctrl *ArticleController) Favorite(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.articleService.FavoriteArticle(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Unfavorite 取消收藏
func (ctrl *ArticleController) Unfavorite(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.articleService.UnfavoriteArticle(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
