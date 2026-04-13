package user

import (
	"blog/internal/middleware"
	"blog/internal/model/dto/request"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService    service.UserService
	articleService service.ArticleService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService, articleService service.ArticleService) *UserController {
	return &UserController{
		userService:    userService,
		articleService: articleService,
	}
}

// GetProfile 获取当前用户信息
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		response.BizError(c, err)
		return
	}

	userResp := ctrl.userService.GetUserResponse(user)
	response.Success(c, userResp)
}

// UpdateProfile 更新用户信息
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.userService.UpdateUser(userID, &req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// ChangePassword 修改密码
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.userService.ChangePassword(userID, &req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// ListUsers 获取用户列表（分页）
func (ctrl *UserController) ListUsers(c *gin.Context) {
	var req request.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pageResp, err := ctrl.userService.ListUsers(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, pageResp)
}

// ListMyArticles 获取“我的文章列表”
func (ctrl *UserController) ListMyArticles(c *gin.Context) {
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
