package response

// CommentBlock 评论主体
type CommentBlock struct {
	ID         uint   `json:"id"`
	ParentID   int    `json:"parent_id"` // 0 表示一级评论（JSON 用 int 表达可空语义）
	RootID     int    `json:"root_id"`
	Content    string `json:"content"`
	LikeCount  int    `json:"like_count"`
	ReplyCount int    `json:"reply_count"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// CommentAuthor 评论作者简要信息
type CommentAuthor struct {
	ID       uint   `json:"id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

// CommentTreeNode 评论树节点（含嵌套回复）
type CommentTreeNode struct {
	Comment CommentBlock       `json:"comment"`
	Author  CommentAuthor      `json:"author"`
	Replies []*CommentTreeNode `json:"replies"`
}
