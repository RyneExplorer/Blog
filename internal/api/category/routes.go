package category

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册 /api/categories
func (ctrl *CategoryController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/categories", ctrl.List)
}
