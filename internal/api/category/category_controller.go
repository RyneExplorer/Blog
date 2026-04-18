package category

import (
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryService service.CategoryService
}

// NewCategoryController 创建分类控制器
func NewCategoryController(categoryService service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

// List 获取分类列表接口
func (ctrl *CategoryController) List(c *gin.Context) {
	data, err := ctrl.categoryService.ListAll(c.Request.Context())
	if err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, data)
}
