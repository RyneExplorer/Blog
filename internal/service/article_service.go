package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/response"
	"blog/pkg/utils"

	"gorm.io/gorm"
)

// ArticleService 文章业务
type ArticleService interface {
	ListArticles(ctx context.Context, q *request.ArticleListQuery) (*response.PageResponse, error)
	GetArticleDetail(ctx context.Context, id uint, viewerUserID uint) (*dto.ArticleDetailResponse, error)
	IncrementView(ctx context.Context, id uint) error
	LikeArticle(ctx context.Context, userID, articleID uint) error
	UnlikeArticle(ctx context.Context, userID, articleID uint) error
	FavoriteArticle(ctx context.Context, userID, articleID uint) error
	UnfavoriteArticle(ctx context.Context, userID, articleID uint) error

	// 用户文章模块
	ListMyArticles(ctx context.Context, userID uint, q *request.MyArticleListQuery) (*response.PageResponse, error)
	ListMyFavorites(ctx context.Context, userID uint, q *request.MyArticleListQuery) (*response.PageResponse, error)
	CreateDraft(ctx context.Context, userID uint, req *request.CreateArticleRequest) (uint, error)
	UpdateDraft(ctx context.Context, userID uint, articleID uint, req *request.UpdateArticleRequest) error
	Publish(ctx context.Context, userID uint, articleID uint) error
	Delete(ctx context.Context, userID uint, articleID uint) error
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

func (s *articleService) GetArticleDetail(ctx context.Context, id uint, viewerUserID uint) (*dto.ArticleDetailResponse, error) {
	a, err := s.articleRepo.GetByIDWithCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	if a.Status != 2 && a.UserID != viewerUserID {
		return nil, bizerrors.New(bizerrors.CodeForbidden, "无权限查看该文章")
	}
	cids := make([]uint, 0, len(a.Categories))
	for _, c := range a.Categories {
		cids = append(cids, c.ID)
	}
	summary := strings.TrimSpace(a.Summary)
	if summary == "" {
		summary = utils.TruncateRunes(a.Content, 100)
	}
	return &dto.ArticleDetailResponse{
		ID:            a.ID,
		Title:         a.Title,
		Summary:       summary,
		Content:       a.Content,
		CoverImage:    a.CoverImage,
		Status:        a.Status,
		ViewCount:     a.ViewCount,
		LikeCount:     int(a.LikeCount),
		FavoriteCount: int(a.FavoriteCount),
		CommentCount:  int(a.CommentCount),
		CategoryIDs:   cids,
		CreatedAt:     formatDateTime(a.CreatedAt),
		UpdatedAt:     formatDateTime(a.UpdatedAt),
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

func (s *articleService) ListMyArticles(ctx context.Context, userID uint, q *request.MyArticleListQuery) (*response.PageResponse, error) {
	sort := q.Sort
	if sort == "" {
		sort = "latest"
	}
	offset := (q.Page - 1) * q.PageSize
	total, err := s.articleRepo.CountByUser(ctx, userID, q.CategoryID)
	if err != nil {
		return nil, err
	}
	rows, err := s.articleRepo.ListByUserWithJoin(ctx, userID, offset, q.PageSize, q.CategoryID, sort)
	if err != nil {
		return nil, err
	}
	list := make([]dto.MyArticleListItem, 0, len(rows))
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
		list = append(list, dto.MyArticleListItem{
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
			Category: cat,
		})
	}
	return response.NewPageResponse(list, total, q.Page, q.PageSize), nil
}

func (s *articleService) ListMyFavorites(ctx context.Context, userID uint, q *request.MyArticleListQuery) (*response.PageResponse, error) {
	sort := q.Sort
	if sort == "" {
		sort = "latest"
	}
	offset := (q.Page - 1) * q.PageSize
	total, err := s.articleRepo.CountFavorites(ctx, userID, q.CategoryID)
	if err != nil {
		return nil, err
	}
	rows, err := s.articleRepo.ListFavoritesWithJoin(ctx, userID, offset, q.PageSize, q.CategoryID, sort)
	if err != nil {
		return nil, err
	}
	list := make([]dto.MyArticleListItem, 0, len(rows))
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
		list = append(list, dto.MyArticleListItem{
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
			Category: cat,
		})
	}
	return response.NewPageResponse(list, total, q.Page, q.PageSize), nil
}

func (s *articleService) CreateDraft(ctx context.Context, userID uint, req *request.CreateArticleRequest) (uint, error) {
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	summary := strings.TrimSpace(req.Summary)
	if summary == "" {
		summary = utils.TruncateRunes(content, 100)
	}
	a := &entity.Article{
		UserID:     userID,
		Title:      title,
		Content:    content,
		Summary:    summary,
		CoverImage: strings.TrimSpace(req.CoverImage),
		Status:     0,
	}
	if err := s.articleRepo.CreateWithCategoriesInTx(ctx, a, req.CategoryIDs); err != nil {
		return 0, err
	}
	return a.ID, nil
}

func (s *articleService) UpdateDraft(ctx context.Context, userID uint, articleID uint, req *request.UpdateArticleRequest) error {
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	summary := strings.TrimSpace(req.Summary)
	if summary == "" {
		summary = utils.TruncateRunes(content, 100)
	}
	updates := map[string]interface{}{
		"title":       title,
		"content":     content,
		"summary":     summary,
		"cover_image": strings.TrimSpace(req.CoverImage),
	}
	ok, err := s.articleRepo.UpdateByAuthorWithCategoriesInTx(ctx, articleID, userID, updates, req.CategoryIDs)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeForbidden, "仅作者本人可修改该文章")
	}
	return nil
}

func (s *articleService) Publish(ctx context.Context, userID uint, articleID uint) error {
	a, err := s.articleRepo.GetByIDWithCategories(ctx, articleID)
	if err != nil {
		return err
	}
	if a == nil {
		return bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}
	if a.UserID != userID {
		return bizerrors.New(bizerrors.CodeForbidden, "仅作者本人可发布该文章")
	}
	if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.Content) == "" {
		return bizerrors.New(bizerrors.CodeBadRequest, "标题和内容不能为空")
	}
	ok, err := s.articleRepo.UpdateStatusByAuthor(ctx, articleID, userID, 1)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeForbidden, "仅作者本人可发布该文章")
	}
	return nil
}

func (s *articleService) Delete(ctx context.Context, userID uint, articleID uint) error {
	ok, err := s.articleRepo.DeleteByAuthorInTx(ctx, articleID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return bizerrors.New(bizerrors.CodeForbidden, "仅作者本人可删除该文章")
	}
	return nil
}
