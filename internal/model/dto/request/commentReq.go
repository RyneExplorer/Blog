package request

// CommentListQuery 评论列表分页
type CommentListQuery struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"pageSize" binding:"required,min=1,max=100"`
}

// CreateCommentRequest 发表评论
type CreateCommentRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required,min=1,max=5000"`
}

// ReplyCommentRequest 回复评论
type ReplyCommentRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required,min=1,max=5000"`
	ParentID  uint   `json:"parent_id" binding:"required"`
	RootID    uint   `json:"root_id" binding:"required"`
}
