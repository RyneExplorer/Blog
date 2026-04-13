package entity

import "time"

// Category 文章分类
type Category struct {
	BaseEntity
	Name string `gorm:"type:varchar(100);not null;comment:分类名称" json:"name"`
	Slug string `gorm:"type:varchar(100);not null;uniqueIndex;comment:URL 标识" json:"slug"`

	Articles []Article `gorm:"many2many:article_categories;" json:"articles"`
}

func (Category) TableName() string {
	return "categories"
}

// ArticleCategory 文章分类中间关联表
type ArticleCategory struct {
	ArticleID  uint      `gorm:"primaryKey;column:article_id;comment:文章ID"`
	CategoryID uint      `gorm:"primaryKey;column:category_id;comment:分类ID"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ArticleCategory) TableName() string {
	return "article_categories"
}
