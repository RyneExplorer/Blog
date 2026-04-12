package entity

// Category 文章分类
// Articles 与 Article 为多对多，中间表 article_categories
type Category struct {
	BaseEntity
	Name string `gorm:"type:varchar(100);not null;comment:分类名称" json:"name"`
	Slug string `gorm:"type:varchar(100);not null;uniqueIndex;comment:URL 标识" json:"slug"`

	Articles []Article `gorm:"many2many:article_categories" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}
