package article

import (
	"strconv"

	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"blog/pkg/upload"

	"github.com/gin-gonic/gin"
)

// ArticleController 处理文章相关接口
type ArticleController struct {
	articleService service.ArticleService
	commentService service.CommentService
}

// NewArticleController 创建文章控制器
func NewArticleController(articleService service.ArticleService, commentService service.CommentService) *ArticleController {
	return &ArticleController{
		articleService: articleService,
		commentService: commentService,
	}
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

// UploadCover 上传文章封面
func (ctrl *ArticleController) UploadCover(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := upload.SaveImage(c, "file", "article")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cover_image": result.URL,
		"url":         result.URL,
	})
}

// UploadContentImage 上传正文图片
func (ctrl *ArticleController) UploadContentImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	result, err := upload.SaveImage(c, "file", "article-content")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"content_image": result.URL,
		"url":           result.URL,
	})
}

// Detail 获取首页文章详情
func (ctrl *ArticleController) Detail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	data, err := ctrl.articleService.GetArticleDetail(c.Request.Context(), id, 0)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

// MyDetail 获取我的文章详情
func (ctrl *ArticleController) MyDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	data, err := ctrl.articleService.GetArticleDetail(c.Request.Context(), id, userID)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

// ListMyArticles 获取我的文章列表
func (ctrl *ArticleController) ListMyArticles(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}
	var q request.MyArticleListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确："+err.Error())
		return
	}
	page, err := ctrl.articleService.ListMyArticles(c.Request.Context(), userID, &q)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, page)
}

// ListFavorites 获取我的收藏列表
func (ctrl *ArticleController) ListFavorites(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}
	var q request.MyArticleListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确："+err.Error())
		return
	}
	page, err := ctrl.articleService.ListMyFavorites(c.Request.Context(), userID, &q)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, page)
}

// Create 创建文章草稿
func (ctrl *ArticleController) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	var req request.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体格式错误："+err.Error())
		return
	}
	id, err := ctrl.articleService.CreateDraft(c.Request.Context(), userID, &req)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, gin.H{"article_id": id})
}

// Update 更新文章
func (ctrl *ArticleController) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	var req request.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体格式错误："+err.Error())
		return
	}
	if err := ctrl.articleService.UpdateDraft(c.Request.Context(), userID, id, &req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Publish 发布文章
func (ctrl *ArticleController) Publish(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	if err := ctrl.articleService.Publish(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Delete 删除文章
func (ctrl *ArticleController) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	if err := ctrl.articleService.Delete(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// ListComments 获取文章评论列表
func (ctrl *ArticleController) ListComments(c *gin.Context) {
	articleID, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	var q request.CommentListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确："+err.Error())
		return
	}
	page, err := ctrl.commentService.ListByArticle(c.Request.Context(), articleID, &q)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, page)
}

// IncrView 增加文章浏览量
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
