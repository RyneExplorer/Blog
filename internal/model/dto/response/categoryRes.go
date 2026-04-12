package response

// CategoryItem 分类列表项
type CategoryItem struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
