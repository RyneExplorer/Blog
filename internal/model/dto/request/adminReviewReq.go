package request

// AdminArticleListRequest 获取审核文章列表参数
type AdminArticleListRequest struct {
	Page       int    `form:"page" binding:"required,min=1"`
	PageSize   int    `form:"pageSize" binding:"required,min=1,max=100"`
	CategoryID *uint  `form:"category_id"`
	Username   string `form:"username" binding:"omitempty,max=50"`
	Status     *int   `form:"status" binding:"omitempty,min=0,max=4"`
}

// ReviewRequest 审核驳回请求体
type ReviewRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=255"`
}

// UpdateCategoryRequest 修改文章分类请求体
type UpdateCategoryRequest struct {
	CategoryIDs []uint `json:"category_ids" binding:"required"`
}
