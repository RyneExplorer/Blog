package comment

import (
	"strconv"

	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// CommentController 评论控制器
type CommentController struct {
	commentService service.CommentService
}

// NewCommentController 创建评论控制器
func NewCommentController(commentService service.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
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

// ListByArticle 获取文章评论列表接口
func (ctrl *CommentController) ListByArticle(c *gin.Context) {
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

// Create 发表评论接口
func (ctrl *CommentController) Create(c *gin.Context) {
	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体格式错误："+err.Error())
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.commentService.CreateComment(c.Request.Context(), userID, &req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Reply 回复评论接口
func (ctrl *CommentController) Reply(c *gin.Context) {
	var req request.ReplyCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体格式错误："+err.Error())
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.commentService.ReplyComment(c.Request.Context(), userID, &req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Delete 删除评论接口
func (ctrl *CommentController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "评论 ID 无效")
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.commentService.DeleteComment(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Like 评论点赞接口
func (ctrl *CommentController) Like(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "评论 ID 无效")
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.commentService.LikeComment(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Unlike 取消评论点赞接口
func (ctrl *CommentController) Unlike(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "评论 ID 无效")
		return
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	if err := ctrl.commentService.UnlikeComment(c.Request.Context(), userID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
