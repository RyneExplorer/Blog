package entity

import "time"

// ReviewLog 审核日志
type ReviewLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"index;not null;comment:文章ID" json:"article_id"`
	AdminID   uint      `gorm:"index;not null;comment:审核人ID" json:"admin_id"`
	Action    string    `gorm:"type:varchar(20);not null;comment:approve/reject/ban" json:"action"`
	Reason    string    `gorm:"type:varchar(255);default:'';comment:原因/备注" json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

func (ReviewLog) TableName() string {
	return "review_logs"
}
