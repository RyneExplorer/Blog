package entity

import "time"

// Image 图片表
type Image struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null;comment:上传者ID" json:"user_id"`
	URL       string    `gorm:"type:varchar(255);not null;comment:图片URL" json:"url"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Image) TableName() string {
	return "images"
}
