package entity

import "time"

// 点赞目标类型（与 likes.target_type 一致）
const (
	LikeTargetArticle = "article"
	LikeTargetComment = "comment"
	LikeTargetImage   = "image"
)

// Like 统一点赞表（复合主键 user_id + target_type + target_id）
type Like struct {
	UserID     uint      `gorm:"primaryKey;column:user_id"`
	TargetType string    `gorm:"primaryKey;column:target_type;type:varchar(20);not null"`
	TargetID   uint      `gorm:"primaryKey;column:target_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (Like) TableName() string {
	return "likes"
}
