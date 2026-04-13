package super

import (
	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReviewController 管理员审核控制器
type ReviewController struct {
	reviewService service.ReviewService
}

func NewReviewController(reviewService service.ReviewService) *ReviewController {
	return &ReviewController{reviewService: reviewService}
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

// List GET /api/super/articles
func (ctrl *ReviewController) List(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	var q request.AdminArticleListRequest
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确："+err.Error())
		return
	}
	page, err := ctrl.reviewService.List(c.Request.Context(), adminID, &q)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, page)
}

// Detail GET /api/super/articles/:id
func (ctrl *ReviewController) Detail(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	data, err := ctrl.reviewService.Detail(c.Request.Context(), adminID, id)
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}

// Approve POST /api/super/articles/:id/approve
func (ctrl *ReviewController) Approve(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	if err := ctrl.reviewService.Approve(c.Request.Context(), adminID, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Reject POST /api/super/articles/:id/reject
func (ctrl *ReviewController) Reject(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	var req request.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体格式错误："+err.Error())
		return
	}
	if err := ctrl.reviewService.Reject(c.Request.Context(), adminID, id, req.Reason); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// Ban POST /api/super/articles/:id/ban
func (ctrl *ReviewController) Ban(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	var req request.ReviewRequest
	_ = c.ShouldBindJSON(&req)
	if err := ctrl.reviewService.Ban(c.Request.Context(), adminID, id, req.Reason); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

// UpdateCategory PUT /api/super/articles/:id/category
func (ctrl *ReviewController) UpdateCategory(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	if adminID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.BadRequest(c, "文章 ID 无效")
		return
	}
	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体格式错误："+err.Error())
		return
	}
	if err := ctrl.reviewService.UpdateCategory(c.Request.Context(), adminID, id, req.CategoryIDs); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}
