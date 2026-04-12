package request

// ArticleListQuery 文章列表查询参数
type ArticleListQuery struct {
	Page       int    `form:"page" binding:"required,min=1"`
	PageSize   int    `form:"pageSize" binding:"required,min=1,max=100"`
	CategoryID *uint  `form:"category_id"`
	Sort       string `form:"sort" binding:"omitempty,oneof=latest hottest"`
}
