package response

// ArticleDetailResponse 文章详情（用于“我的文章详情”与公开详情）
type ArticleDetailResponse struct {
	ID            uint            `json:"id"`
	Title         string          `json:"title"`
	Summary       string          `json:"summary"`
	Content       string          `json:"content"`
	CoverImage    string          `json:"cover_image"`
	Categories    []CategoryBrief `json:"categories"`
	Username      string          `json:"username"`
	Nickname      string          `json:"nickname"`
	Bio           string          `json:"bio"`
	Avatar        string          `json:"avatar"`
	Status        int             `json:"status"`
	ViewCount     int             `json:"view_count"`
	LikeCount     int             `json:"like_count"`
	FavoriteCount int             `json:"favorite_count"`
	CommentCount  int             `json:"comment_count"`
	Liked         bool            `json:"liked"`
	Favorited     bool            `json:"favorited"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

// MyArticleListItem “我的文章列表”单项
type MyArticleListItem struct {
	Article    ArticleBrief    `json:"article"`
	Categories []CategoryBrief `json:"categories"`
}

// ArticleBrief 文章列表中的 article 块
type ArticleBrief struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	CoverImage    string `json:"cover_image"`
	Status        int    `json:"status"`
	ViewCount     int    `json:"view_count"`
	LikeCount     int    `json:"like_count"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	Liked         bool   `json:"liked"`
	Favorited     bool   `json:"favorited"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// AuthorProfile 作者信息
type AuthorProfile struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

// CategoryBrief 分类信息
type CategoryBrief struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ArticleListItem 文章列表单项
type ArticleListItem struct {
	Article    ArticleBrief    `json:"article"`
	Author     AuthorProfile   `json:"author"`
	Categories []CategoryBrief `json:"categories"`
}

// ArticleDetail 文章详情
type ArticleDetail struct {
	ID            uint            `json:"id"`
	Title         string          `json:"title"`
	Categories    []CategoryBrief `json:"categories"`
	Nickname      string          `json:"nickname"`
	Bio           string          `json:"bio"`
	Avatar        string          `json:"avatar"`
	Content       string          `json:"content"`
	CoverImage    string          `json:"cover_image"`
	Status        int             `json:"status"`
	ViewCount     int             `json:"view_count"`
	LikeCount     int             `json:"like_count"`
	FavoriteCount int             `json:"favorite_count"`
	CommentCount  int             `json:"comment_count"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}
