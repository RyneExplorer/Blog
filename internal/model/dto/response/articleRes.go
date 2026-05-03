package response

// ArticleDetailResponse 文章详情（用于“我的文章详情”与公开详情）
type ArticleDetailResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	Content       string `json:"content"`
	CoverImage    string `json:"cover_image"`
	CategoryName  string `json:"category_name"`
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	Bio           string `json:"bio"`
	Avatar        string `json:"avatar"`
	Status        int    `json:"status"`
	ViewCount     int    `json:"view_count"`
	LikeCount     int    `json:"like_count"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	CategoryIDs   []uint `json:"category_ids"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// MyArticleListItem “我的文章列表”单项
type MyArticleListItem struct {
	Article  ArticleBrief  `json:"article"`
	Category CategoryBrief `json:"category"`
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
	Article  ArticleBrief  `json:"article"`
	Author   AuthorProfile `json:"author"`
	Category CategoryBrief `json:"category"`
}

// ArticleDetail 文章详情
type ArticleDetail struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	CategoryName  string `json:"category_name"`
	Nickname      string `json:"nickname"`
	Bio           string `json:"bio"`
	Avatar        string `json:"avatar"`
	Content       string `json:"content"`
	CoverImage    string `json:"cover_image"`
	Status        int    `json:"status"`
	ViewCount     int    `json:"view_count"`
	LikeCount     int    `json:"like_count"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
