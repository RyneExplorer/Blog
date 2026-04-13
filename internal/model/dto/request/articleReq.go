package request

// MyArticleListQuery 获取“我的文章列表”查询参数
type MyArticleListQuery struct {
	Page       int    `form:"page" binding:"required,min=1"`
	PageSize   int    `form:"pageSize" binding:"required,min=1,max=100"`
	CategoryID *uint  `form:"category_id"`
	Sort       string `form:"sort" binding:"omitempty,oneof=latest hottest"`
}

// ArticleListQuery 文章列表查询参数
type ArticleListQuery struct {
	Page       int    `form:"page" binding:"required,min=1"`
	PageSize   int    `form:"pageSize" binding:"required,min=1,max=100"`
	CategoryID *uint  `form:"category_id"`
	Sort       string `form:"sort" binding:"omitempty,oneof=latest hottest"`
}

// CreateArticleRequest 创建文章（草稿）
type CreateArticleRequest struct {
	Title       string `json:"title" binding:"omitempty,max=255"`
	Content     string `json:"content" binding:"omitempty"`
	Summary     string `json:"summary" binding:"omitempty,max=255"`
	CoverImage  string `json:"cover_image" binding:"omitempty,max=500"`
	CategoryIDs []uint `json:"category_ids" binding:"omitempty"`
}

// UpdateArticleRequest 更新文章（自动保存草稿）
type UpdateArticleRequest struct {
	Title       string `json:"title" binding:"omitempty,max=255"`
	Content     string `json:"content" binding:"omitempty"`
	Summary     string `json:"summary" binding:"omitempty,max=255"`
	CoverImage  string `json:"cover_image" binding:"omitempty,max=500"`
	CategoryIDs []uint `json:"category_ids" binding:"omitempty"`
}
