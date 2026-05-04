package repository

import "context"

// LikeRepository 统一点赞仓储接口（likes 表，target_type 区分目标）
type LikeRepository interface {
	LikeArticleInTx(ctx context.Context, userID, articleID uint) error
	UnlikeArticleInTx(ctx context.Context, userID, articleID uint) error
	LikeCommentInTx(ctx context.Context, userID, commentID uint) error
	UnlikeCommentInTx(ctx context.Context, userID, commentID uint) error
	ListLikedArticleIDs(ctx context.Context, userID uint, articleIDs []uint) (map[uint]bool, error)
	ListLikedCommentIDs(ctx context.Context, userID uint, commentIDs []uint) (map[uint]bool, error)
}
