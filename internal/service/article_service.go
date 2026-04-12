package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/response"
	"blog/pkg/utils"

	"gorm.io/gorm"
)

// ArticleService 文章业务
type ArticleService interface {
	ListArticles(ctx context.Context, q *request.ArticleListQuery) (*response.PageResponse, error)
	GetArticleDetail(ctx context.Context, id uint) (*dto.ArticleDetail, error)
	IncrementView(ctx context.Context, id uint) error
	LikeArticle(ctx context.Context, userID, articleID uint) error
	UnlikeArticle(ctx context.Context, userID, articleID uint) error
	FavoriteArticle(ctx context.Context, userID, articleID uint) error
	UnfavoriteArticle(ctx context.Context, userID, articleID uint) error
}

type articleService struct {
	articleRepo  repository.ArticleRepository
	likeRepo     repository.LikeRepository
	favoriteRepo repository.FavoriteRepository
}

// NewArticleService 创建文章服务
func NewArticleService(
	articleRepo repository.ArticleRepository,
	likeRepo repository.LikeRepository,
	favoriteRepo repository.FavoriteRepository,
) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		likeRepo:     likeRepo,
		favoriteRepo: favoriteRepo,
	}
}

func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func (s *articleService) ListArticles(ctx context.Context, q *request.ArticleListQuery) (*response.PageResponse, error) {
	sort := q.Sort
	if sort == "" {
		sort = "latest"
	}
	offset := (q.Page - 1) * q.PageSize
	total, err := s.articleRepo.CountPublished(ctx, q.CategoryID)
	if err != nil {
		return nil, err
	}
	rows, err := s.articleRepo.ListPublishedWithJoin(ctx, offset, q.PageSize, q.CategoryID, sort)
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
			cat = dto.CategoryBrief{
				ID:   uint(row.CategoryRefID.Int64),
				Name: row.CategoryName,
				Slug: row.CategorySlug,
			}
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
				CreatedAt:     formatDateTime(row.CreatedAt),
				UpdatedAt:     formatDateTime(row.UpdatedAt),
			},
			Author: dto.AuthorProfile{
				ID:       row.UserID,
				Nickname: row.AuthorNickname,
				Avatar:   row.AuthorAvatar,
				Bio:      row.AuthorBio,
			},
			Category: cat,
		})
	}
	return response.NewPageResponse(list, total, q.Page, q.PageSize), nil
}

func (s *articleService) GetArticleDetail(ctx context.Context, id uint) (*dto.ArticleDetail, error) {
	row, err := s.articleRepo.GetPublishedDetailJoin(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	return &dto.ArticleDetail{
		ID:            row.ID,
		Title:         row.Title,
		CategoryName:  row.CategoryName,
		Nickname:      row.Nickname,
		Bio:           row.Bio,
		Avatar:        row.Avatar,
		Content:       row.Content,
		CoverImage:    row.CoverImage,
		Status:        row.Status,
		ViewCount:     row.ViewCount,
		LikeCount:     int(row.LikeCount),
		FavoriteCount: int(row.FavoriteCount),
		CommentCount:  int(row.CommentCount),
		CreatedAt:     formatDateTime(row.CreatedAt),
		UpdatedAt:     formatDateTime(row.UpdatedAt),
	}, nil
}

func (s *articleService) IncrementView(ctx context.Context, id uint) error {
	if err := s.articleRepo.IncrementViewInTx(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
		}
		return err
	}
	return nil
}

func (s *articleService) LikeArticle(ctx context.Context, userID, articleID uint) error {
	err := s.likeRepo.LikeArticleInTx(ctx, userID, articleID)
	if err == nil {
		return nil
	}
	if utils.IsMySQLDuplicateKey(err) {
		return bizerrors.New(bizerrors.CodeConflict, "您已点赞该文章，请勿重复点赞")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	return err
}

func (s *articleService) UnlikeArticle(ctx context.Context, userID, articleID uint) error {
	err := s.likeRepo.UnlikeArticleInTx(ctx, userID, articleID)
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrArticleUnlikeMissing) {
		return bizerrors.New(bizerrors.CodeBadRequest, "您尚未点赞该文章")
	}
	return err
}

func (s *articleService) FavoriteArticle(ctx context.Context, userID, articleID uint) error {
	err := s.favoriteRepo.FavoriteArticleInTx(ctx, userID, articleID)
	if err == nil {
		return nil
	}
	if utils.IsMySQLDuplicateKey(err) {
		return bizerrors.New(bizerrors.CodeConflict, "您已收藏该文章")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	return err
}

func (s *articleService) UnfavoriteArticle(ctx context.Context, userID, articleID uint) error {
	err := s.favoriteRepo.UnfavoriteArticleInTx(ctx, userID, articleID)
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrUnfavoriteMissing) {
		return bizerrors.New(bizerrors.CodeBadRequest, "您尚未收藏该文章")
	}
	return err
}
