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
	// 1. 统一排序参数并计算分页偏移量。
	sort := q.Sort
	if sort == "" {
		sort = "latest"
	}
	offset := (q.Page - 1) * q.PageSize

	// 2. 先查询总数，再按分页条件查询当前页数据。
	total, err := s.articleRepo.CountPublished(ctx, q.CategoryID)
	if err != nil {
		return nil, err
	}
	rows, err := s.articleRepo.ListPublishedWithJoin(ctx, offset, q.PageSize, q.CategoryID, sort)
	if err != nil {
		return nil, err
	}

	// 3. 最后把仓储层返回的行数据转换成接口响应结构。
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
	// 1. 先按文章 ID 加载文章主体和分类信息。
	a, err := s.articleRepo.GetByIDWithCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, bizerrors.New(bizerrors.CodeNotFound, "文章不存在")
	}

	// 2. 再根据文章状态和查看人身份判断是否有权访问。
	if a.Status != 2 && a.UserID != viewerUserID {
		return nil, bizerrors.New(bizerrors.CodeForbidden, "无权限查看该文章")
	}

	cids := make([]uint, 0, len(a.Categories))
	categoryName := ""
	for _, c := range a.Categories {
		cids = append(cids, c.ID)
		if categoryName == "" {
			categoryName = c.Name
		}
	}

	summary := strings.TrimSpace(a.Summary)
	if summary == "" {
		summary = utils.TruncateRunes(a.Content, 100)
	}

	authorName := strings.TrimSpace(a.User.Nickname)
	if authorName == "" {
		authorName = strings.TrimSpace(a.User.Username)
	}

	// 3. 最后补齐摘要、分类 ID 等衍生字段并返回详情响应。
	return &dto.ArticleDetailResponse{
		ID:            a.ID,
		Title:         a.Title,
		Summary:       summary,
		Content:       a.Content,
		CoverImage:    a.CoverImage,
		CategoryName:  categoryName,
		Username:      a.User.Username,
		Nickname:      authorName,
		Bio:           a.User.Bio,
		Avatar:        a.User.Avatar,
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
	// 1. 统一我的文章列表的排序与分页参数。
	sort := q.Sort
	if sort == "" {
		sort = "latest"
	}
	offset := (q.Page - 1) * q.PageSize

	// 2. 查询总数和当前页数据，保证分页信息完整。
	total, err := s.articleRepo.CountByUser(ctx, userID, q.CategoryID)
	if err != nil {
		return nil, err
	}
	rows, err := s.articleRepo.ListByUserWithJoin(ctx, userID, offset, q.PageSize, q.CategoryID, sort)
	if err != nil {
		return nil, err
	}

	// 3. 将仓储层结果转换为前端需要的列表结构。
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
	// 1. 统一收藏列表的排序与分页参数。
	sort := q.Sort
	if sort == "" {
		sort = "latest"
	}
	offset := (q.Page - 1) * q.PageSize

	// 2. 查询收藏总数和当前页文章数据。
	total, err := s.articleRepo.CountFavorites(ctx, userID, q.CategoryID)
	if err != nil {
		return nil, err
	}
	rows, err := s.articleRepo.ListFavoritesWithJoin(ctx, userID, offset, q.PageSize, q.CategoryID, sort)
	if err != nil {
		return nil, err
	}

	// 3. 组装成与“我的文章”一致的列表结构，便于前端复用。
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
	// 1. 先对标题、正文、摘要、封面等输入做去空格规范化。
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	summary := strings.TrimSpace(req.Summary)

	// 2. 如果摘要为空，则从正文中截取一段内容生成默认摘要。
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
	// 3. 最后在事务里创建草稿文章并写入文章分类关联。
	if err := s.articleRepo.CreateWithCategoriesInTx(ctx, a, req.CategoryIDs); err != nil {
		return 0, err
	}
	return a.ID, nil
}

func (s *articleService) UpdateDraft(ctx context.Context, userID uint, articleID uint, req *request.UpdateArticleRequest) error {
	// 1. 先规范化请求中的标题、正文、摘要和封面字段。
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	summary := strings.TrimSpace(req.Summary)
	if summary == "" {
		summary = utils.TruncateRunes(content, 100)
	}

	// 2. 再构造待更新字段集合，避免把业务判断下沉到仓储层。
	updates := map[string]interface{}{
		"title":       title,
		"content":     content,
		"summary":     summary,
		"cover_image": strings.TrimSpace(req.CoverImage),
	}

	// 3. 最后在事务里同时更新文章主体和分类关联，保证数据一致。
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
	// 1. 先读取文章并校验文章存在且属于当前作者。
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

	// 2. 再校验标题和正文等发布前置条件，避免空内容进入审核流。
	if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.Content) == "" {
		return bizerrors.New(bizerrors.CodeBadRequest, "标题和内容不能为空")
	}

	// 3. 条件满足后把文章状态推进到待审核。
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
	// 1. 通过作者 ID 和文章 ID 限定删除范围，避免越权删除他人文章。
	// 2. 在事务中删除文章主体及其分类关联。
	ok, err := s.articleRepo.DeleteByAuthorInTx(ctx, articleID, userID)
	if err != nil {
		return err
	}

	// 3. 若未命中记录，则统一按业务语义返回错误结果。
	if !ok {
		return bizerrors.New(bizerrors.CodeForbidden, "仅作者本人可删除该文章")
	}
	return nil
}
