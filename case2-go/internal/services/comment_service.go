package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"majootest/case2-go/internal/models"
	"majootest/case2-go/internal/repositories"
)

// CommentService defines the business logic for comment management.
type CommentService interface {
	Create(ctx context.Context, postID, userID int64, req *models.CreateCommentRequest) (*models.Comment, error)
	GetByID(ctx context.Context, id int64) (*models.Comment, error)
	ListByPostID(ctx context.Context, postID int64, page, limit int) ([]*models.Comment, *models.PaginationMeta, error)
	Update(ctx context.Context, commentID, userID int64, req *models.UpdateCommentRequest) (*models.Comment, error)
	Delete(ctx context.Context, commentID, userID int64) error
}

type commentService struct {
	commentRepo repositories.CommentRepository
	postRepo    repositories.PostRepository
	txManager   repositories.TxManager
}

// NewCommentService creates a new CommentService instance.
func NewCommentService(
	commentRepo repositories.CommentRepository,
	postRepo repositories.PostRepository,
	txManager repositories.TxManager,
) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		txManager:   txManager,
	}
}

func (s *commentService) Create(ctx context.Context, postID, userID int64, req *models.CreateCommentRequest) (*models.Comment, error) {
	// Verify post exists
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		return nil, err
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	// Atomically insert comment and increment post comment count
	err := s.txManager.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := s.commentRepo.CreateWithTx(ctx, tx, comment); err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
		if err := s.postRepo.IncrementCommentCount(ctx, tx, postID, 1); err != nil {
			return fmt.Errorf("failed to increment post comment count: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.commentRepo.GetByID(ctx, comment.ID)
}

func (s *commentService) GetByID(ctx context.Context, id int64) (*models.Comment, error) {
	return s.commentRepo.GetByID(ctx, id)
}

func (s *commentService) ListByPostID(ctx context.Context, postID int64, page, limit int) ([]*models.Comment, *models.PaginationMeta, error) {
	// Verify post exists
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		return nil, nil, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	comments, totalItems, err := s.commentRepo.ListByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list comments: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	return comments, meta, nil
}

func (s *commentService) Update(ctx context.Context, commentID, userID int64, req *models.UpdateCommentRequest) (*models.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	// Authorization check: only author can update comment
	if comment.UserID != userID {
		return nil, models.ErrForbidden
	}

	comment.Content = req.Content
	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return s.commentRepo.GetByID(ctx, commentID)
}

func (s *commentService) Delete(ctx context.Context, commentID, userID int64) error {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	// Authorization check: only author can delete comment
	if comment.UserID != userID {
		return models.ErrForbidden
	}

	// Atomically delete comment and decrement post comment count
	return s.txManager.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := s.commentRepo.DeleteWithTx(ctx, tx, commentID); err != nil {
			return fmt.Errorf("failed to delete comment: %w", err)
		}
		if err := s.postRepo.IncrementCommentCount(ctx, tx, comment.PostID, -1); err != nil {
			return fmt.Errorf("failed to decrement comment count: %w", err)
		}
		return nil
	})
}
