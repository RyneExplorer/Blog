package service

import (
	"context"
	"errors"
	"sort"

	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/response"
	"blog/pkg/utils"

	"gorm.io/gorm"
)

// CommentService 评论业务
type CommentService interface {
	ListByArticle(ctx context.Context, articleID uint, q *request.CommentListQuery) (*response.PageResponse, error)
	CreateComment(ctx context.Context, userID uint, req *request.CreateCommentRequest) error
	ReplyComment(ctx context.Context, userID uint, req *request.ReplyCommentRequest) error
	DeleteComment(ctx context.Context, userID, commentID uint) error
	LikeComment(ctx context.Context, userID, commentID uint) error
	UnlikeComment(ctx context.Context, userID, commentID uint) error
}

type commentService struct {
	commentRepo repository.CommentRepository
	articleRepo repository.ArticleRepository
	likeRepo    repository.LikeRepository
}

// NewCommentService 创建评论服务
func NewCommentService(
	commentRepo repository.CommentRepository,
	articleRepo repository.ArticleRepository,
	likeRepo repository.LikeRepository,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		likeRepo:    likeRepo,
	}
}

func toCommentBlock(row repository.CommentJoinRow) dto.CommentBlock {
	parent := 0
	if row.ParentID.Valid {
		parent = int(row.ParentID.Int64)
	}
	root := 0
	if row.RootID.Valid {
		root = int(row.RootID.Int64)
	}
	return dto.CommentBlock{
		ID:         row.ID,
		ParentID:   parent,
		RootID:     root,
		Content:    row.Content,
		LikeCount:  row.LikeCount,
		ReplyCount: int(row.ReplyCount),
		CreatedAt:  row.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  row.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ListByArticle 分页一级评论 + 子评论树
func (s *commentService) ListByArticle(ctx context.Context, articleID uint, q *request.CommentListQuery) (*response.PageResponse, error) {
	ok, err := s.articleRepo.ExistsPublished(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	total, err := s.commentRepo.CountRootsByArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}
	offset := (q.Page - 1) * q.PageSize
	rows, err := s.commentRepo.ListJoinedByArticlePage(ctx, articleID, q.PageSize, offset)
	if err != nil {
		return nil, err
	}
	type wrap struct {
		node *dto.CommentTreeNode
		row  repository.CommentJoinRow
	}
	byID := make(map[uint]*wrap)
	children := make(map[uint][]uint)
	var rootIDs []uint

	for _, row := range rows {
		node := &dto.CommentTreeNode{
			Comment: toCommentBlock(row),
			Author: dto.CommentAuthor{
				ID:       row.UserID,
				Avatar:   row.UserAvatar,
				Nickname: row.UserNickname,
			},
			Replies: make([]*dto.CommentTreeNode, 0),
		}
		byID[row.ID] = &wrap{node: node, row: row}
		if !row.ParentID.Valid {
			rootIDs = append(rootIDs, row.ID)
		} else {
			pid := uint(row.ParentID.Int64)
			children[pid] = append(children[pid], row.ID)
		}
	}

	var build func(id uint) *dto.CommentTreeNode
	build = func(id uint) *dto.CommentTreeNode {
		w := byID[id]
		if w == nil {
			return nil
		}
		childIDs := children[id]
		sort.Slice(childIDs, func(i, j int) bool {
			return byID[childIDs[i]].row.CreatedAt.Before(byID[childIDs[j]].row.CreatedAt)
		})
		for _, cid := range childIDs {
			ch := build(cid)
			if ch != nil {
				w.node.Replies = append(w.node.Replies, ch)
			}
		}
		return w.node
	}

	sort.Slice(rootIDs, func(i, j int) bool {
		return byID[rootIDs[i]].row.CreatedAt.After(byID[rootIDs[j]].row.CreatedAt)
	})

	list := make([]*dto.CommentTreeNode, 0, len(rootIDs))
	for _, rid := range rootIDs {
		if n := build(rid); n != nil {
			list = append(list, n)
		}
	}

	return response.NewPageResponse(list, total, q.Page, q.PageSize), nil
}

func (s *commentService) CreateComment(ctx context.Context, userID uint, req *request.CreateCommentRequest) error {
	ok, err := s.articleRepo.ExistsPublished(ctx, req.ArticleID)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	c := &entity.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		Content:   utils.TrimSpace(req.Content),
		Status:    1,
	}
	return s.commentRepo.CreateWithCountersInTx(ctx, c)
}

func expectedRootID(parent *entity.Comment) uint {
	if parent.ParentID == nil {
		return parent.ID
	}
	if parent.RootID != nil {
		return *parent.RootID
	}
	return parent.ID
}

func (s *commentService) ReplyComment(ctx context.Context, userID uint, req *request.ReplyCommentRequest) error {
	ok, err := s.articleRepo.ExistsPublished(ctx, req.ArticleID)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	parent, err := s.commentRepo.GetByID(ctx, req.ParentID)
	if err != nil {
		return err
	}
	if parent == nil || parent.ArticleID != req.ArticleID || parent.Status != 1 {
		return bizerrors.New(bizerrors.CodeBadRequest, "父评论不存在或不属于该文章")
	}
	expRoot := expectedRootID(parent)
	if req.RootID != expRoot {
		return bizerrors.New(bizerrors.CodeBadRequest, "root_id 与评论层级不匹配")
	}
	pid := req.ParentID
	rid := req.RootID
	c := &entity.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		ParentID:  &pid,
		RootID:    &rid,
		Content:   utils.TrimSpace(req.Content),
		Status:    1,
	}
	return s.commentRepo.CreateWithCountersInTx(ctx, c)
}

func (s *commentService) DeleteComment(ctx context.Context, userID, commentID uint) error {
	c, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if c == nil {
		return bizerrors.New(bizerrors.CodeNotFound, "评论不存在")
	}
	if c.UserID != userID {
		return bizerrors.New(bizerrors.CodeForbidden, "无权删除该评论")
	}
	n, err := s.commentRepo.CountChildren(ctx, commentID)
	if err != nil {
		return err
	}
	if n > 0 {
		return bizerrors.New(bizerrors.CodeBadRequest, "请先删除子评论后再删除本条评论")
	}
	return s.commentRepo.DeleteWithCountersInTx(ctx, c)
}

func (s *commentService) LikeComment(ctx context.Context, userID, commentID uint) error {
	err := s.likeRepo.LikeCommentInTx(ctx, userID, commentID)
	if err == nil {
		return nil
	}
	if utils.IsMySQLDuplicateKey(err) {
		return bizerrors.New(bizerrors.CodeConflict, "您已点赞该评论")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return bizerrors.New(bizerrors.CodeNotFound, "评论不存在")
	}
	return err
}

func (s *commentService) UnlikeComment(ctx context.Context, userID, commentID uint) error {
	err := s.likeRepo.UnlikeCommentInTx(ctx, userID, commentID)
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrCommentUnlikeMissing) {
		return bizerrors.New(bizerrors.CodeBadRequest, "您尚未点赞该评论")
	}
	return err
}
