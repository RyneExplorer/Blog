package service

import (
	"context"
	"strings"

	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/response"
	"blog/pkg/utils"
)

// ReviewService 管理员审核服务
type ReviewService interface {
	List(ctx context.Context, adminID uint, q *request.AdminArticleListRequest) (*response.PageResponse, error)
	Detail(ctx context.Context, adminID uint, articleID uint) (*dto.AdminArticleDetailResponse, error)
	Approve(ctx context.Context, adminID uint, articleID uint) error
	Reject(ctx context.Context, adminID uint, articleID uint, reason string) error
	Ban(ctx context.Context, adminID uint, articleID uint, reason string) error
	UpdateCategory(ctx context.Context, adminID uint, articleID uint, categoryIDs []uint) error
}

type reviewService struct {
	reviewRepo  repository.ReviewRepository
	articleRepo repository.ArticleRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, articleRepo repository.ArticleRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo, articleRepo: articleRepo}
}

func formatTime(t interface{ Format(string) string }) string {
	return t.Format("2006-01-02 15:04:05")
}

func (s *reviewService) List(ctx context.Context, adminID uint, q *request.AdminArticleListRequest) (*response.PageResponse, error) {
	_ = adminID

	filter := &repository.AdminListFilter{
		CategoryID: q.CategoryID,
		Username:   strings.TrimSpace(q.Username),
		Status:     q.Status,
	}

	offset := (q.Page - 1) * q.PageSize
	total, err := s.reviewRepo.CountForAdmin(ctx, filter)
	if err != nil {
		return nil, err
	}

	rows, err := s.reviewRepo.ListForAdmin(ctx, filter, offset, q.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]dto.ArticleListItem, 0, len(rows))
	for _, row := range rows {
		summary := strings.TrimSpace(row.Summary.String)
		if summary == "" {
			summary = utils.TruncateRunes(row.Content, 100)
		}

		cat := dto.CategoryBrief{}
		if row.CategoryRefID.Valid {
			cat = dto.CategoryBrief{ID: uint(row.CategoryRefID.Int64), Name: row.CategoryName, Slug: row.CategorySlug}
		}

		list = append(list, dto.ArticleListItem{
			Article: dto.ArticleBrief{
				ID:            row.ID,
				Title:         row.Title,
				Summary:       summary,
				CoverImage:    row.CoverImage,
				Status:        row.Status,
				ViewCount:     row.ViewCount,
				LikeCount:     int(row.LikeCount),
				FavoriteCount: int(row.FavoriteCount),
				CommentCount:  int(row.CommentCount),
				CreatedAt:     formatTime(row.CreatedAt),
				UpdatedAt:     formatTime(row.UpdatedAt),
			},
			Author: dto.AuthorProfile{
				ID:       row.UserID,
				Nickname: row.Nickname,
				Avatar:   row.Avatar,
				Bio:      row.Bio,
			},
			Category: cat,
		})
	}

	return response.NewPageResponse(list, total, q.Page, q.PageSize), nil
}

func (s *reviewService) Detail(ctx context.Context, adminID uint, articleID uint) (*dto.AdminArticleDetailResponse, error) {
	_ = adminID

	// 1. 先查询文章主信息，确认审核目标存在。
	row, err := s.reviewRepo.GetDetail(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}

	// 2. 再补充文章分类信息，便于审核端完整展示。
	a, err := s.articleRepo.GetByIDWithCategories(ctx, articleID)
	if err != nil {
		return nil, err
	}

	cats := make([]dto.CategoryBrief, 0)
	if a != nil {
		cats = make([]dto.CategoryBrief, 0, len(a.Categories))
		for _, c := range a.Categories {
			cats = append(cats, dto.CategoryBrief{ID: c.ID, Name: c.Name, Slug: c.Slug})
		}
	}

	cids := make([]uint, 0)
	if a != nil {
		cids = make([]uint, 0, len(a.Categories))
		for _, c := range a.Categories {
			cids = append(cids, c.ID)
		}
	}

	summary := strings.TrimSpace(row.Summary.String)
	if summary == "" {
		summary = utils.TruncateRunes(row.Content, 100)
	}
	rejectReason := strings.TrimSpace(row.RejectReason.String)

	// 3. 最后统一组装审核详情响应。
	return &dto.AdminArticleDetailResponse{
		Article: dto.ArticleDetailResponse{
			ID:            row.ID,
			Title:         row.Title,
			Summary:       summary,
			Content:       row.Content,
			CoverImage:    row.CoverImage,
			Status:        row.Status,
			ViewCount:     row.ViewCount,
			LikeCount:     int(row.LikeCount),
			FavoriteCount: int(row.FavoriteCount),
			CommentCount:  int(row.CommentCount),
			CategoryIDs:   cids,
			CreatedAt:     formatTime(row.CreatedAt),
			UpdatedAt:     formatTime(row.UpdatedAt),
		},
		Author: dto.AuthorProfile{
			ID:       row.UserID,
			Nickname: row.Nickname,
			Avatar:   row.Avatar,
			Bio:      row.Bio,
		},
		Categories:   cats,
		RejectReason: rejectReason,
	}, nil
}

func (s *reviewService) Approve(ctx context.Context, adminID uint, articleID uint) error {
	ok, err := s.reviewRepo.ApproveInTx(ctx, articleID, adminID)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeBadRequest, "仅待审核文章可通过审核")
	}
	return nil
}

func (s *reviewService) Reject(ctx context.Context, adminID uint, articleID uint, reason string) error {
	ok, err := s.reviewRepo.RejectInTx(ctx, articleID, adminID, reason)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeBadRequest, "仅待审核文章可驳回")
	}
	return nil
}

func (s *reviewService) Ban(ctx context.Context, adminID uint, articleID uint, reason string) error {
	ok, err := s.reviewRepo.BanInTx(ctx, articleID, adminID, reason)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	return nil
}

func (s *reviewService) UpdateCategory(ctx context.Context, adminID uint, articleID uint, categoryIDs []uint) error {
	ok, err := s.reviewRepo.UpdateCategoriesInTx(ctx, articleID, categoryIDs, adminID)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	return nil
}
