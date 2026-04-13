package response

// AdminArticleDetailResponse 管理员审核详情
type AdminArticleDetailResponse struct {
	Article      ArticleDetailResponse `json:"article"`
	Author       AuthorProfile         `json:"author"`
	Categories   []CategoryBrief       `json:"categories"`
	RejectReason string                `json:"reject_reason"`
}
