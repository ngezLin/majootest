package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"majootest/case2-go/internal/models"
	"majootest/case2-go/internal/repositories"
)

// PostService defines the business logic for post management.
type PostService interface {
	Create(ctx context.Context, userID int64, req *models.CreatePostRequest) (*models.Post, error)
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	List(ctx context.Context, page, limit int, search string) ([]*models.Post, *models.PaginationMeta, error)
	Update(ctx context.Context, postID, userID int64, req *models.UpdatePostRequest) (*models.Post, error)
	Delete(ctx context.Context, postID, userID int64) error
}

type postService struct {
	postRepo    repositories.PostRepository
	commentRepo repositories.CommentRepository
	txManager   repositories.TxManager
}

// NewPostService creates a new PostService instance.
func NewPostService(
	postRepo repositories.PostRepository,
	commentRepo repositories.CommentRepository,
	txManager repositories.TxManager,
) PostService {
	return &postService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		txManager:   txManager,
	}
}

func (s *postService) Create(ctx context.Context, userID int64, req *models.CreatePostRequest) (*models.Post, error) {
	post := &models.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Fetch full post with author profile
	return s.postRepo.GetByID(ctx, post.ID)
}

func (s *postService) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	return s.postRepo.GetByID(ctx, id)
}

func (s *postService) List(ctx context.Context, page, limit int, search string) ([]*models.Post, *models.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	posts, totalItems, err := s.postRepo.List(ctx, limit, offset, search)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list posts: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	return posts, meta, nil
}

func (s *postService) Update(ctx context.Context, postID, userID int64, req *models.UpdatePostRequest) (*models.Post, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Authorization: only the author can update their post
	if post.UserID != userID {
		return nil, models.ErrForbidden
	}

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}

	if err := s.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return s.postRepo.GetByID(ctx, postID)
}

func (s *postService) Delete(ctx context.Context, postID, userID int64) error {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	// Authorization: only the author can delete their post
	if post.UserID != userID {
		return models.ErrForbidden
	}

	// Atomically delete comments then post inside a single transaction
	return s.txManager.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := s.commentRepo.DeleteByPostIDWithTx(ctx, tx, postID); err != nil {
			return fmt.Errorf("failed to delete post comments: %w", err)
		}
		if err := s.postRepo.DeleteWithTx(ctx, tx, postID); err != nil {
			return fmt.Errorf("failed to delete post: %w", err)
		}
		return nil
	})
}
