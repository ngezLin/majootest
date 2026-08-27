package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"majootest/case2-go/internal/models"
)

// CommentRepository defines database operations for comments.
type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	CreateWithTx(ctx context.Context, tx *sql.Tx, comment *models.Comment) error
	GetByID(ctx context.Context, id int64) (*models.Comment, error)
	ListByPostID(ctx context.Context, postID int64, limit, offset int) ([]*models.Comment, int64, error)
	Update(ctx context.Context, comment *models.Comment) error
	DeleteWithTx(ctx context.Context, tx *sql.Tx, id int64) error
	DeleteByPostIDWithTx(ctx context.Context, tx *sql.Tx, postID int64) error
}

type mysqlCommentRepository struct {
	db *sql.DB
}

// NewCommentRepository returns a new MySQL comment repository instance.
func NewCommentRepository(db *sql.DB) CommentRepository {
	return &mysqlCommentRepository{db: db}
}

func (r *mysqlCommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	return r.CreateWithTx(ctx, nil, comment)
}

func (r *mysqlCommentRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, comment *models.Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`
	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(ctx, query, comment.PostID, comment.UserID, comment.Content)
	} else {
		result, err = r.db.ExecContext(ctx, query, comment.PostID, comment.UserID, comment.Content)
	}
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get comment last insert id: %w", err)
	}
	comment.ID = id
	return nil
}

func (r *mysqlCommentRepository) GetByID(ctx context.Context, id int64) (*models.Comment, error) {
	query := `
		SELECT 
			c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at,
			u.id, u.name, u.email, u.created_at, u.updated_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`
	comment := &models.Comment{Author: &models.User{}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Author.ID,
		&comment.Author.Name,
		&comment.Author.Email,
		&comment.Author.CreatedAt,
		&comment.Author.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("failed to query comment by id: %w", err)
	}
	return comment, nil
}

func (r *mysqlCommentRepository) ListByPostID(ctx context.Context, postID int64, limit, offset int) ([]*models.Comment, int64, error) {
	countQuery := `SELECT COUNT(*) FROM comments WHERE post_id = ?`
	var totalItems int64
	if err := r.db.QueryRowContext(ctx, countQuery, postID).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	listQuery := `
		SELECT 
			c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at,
			u.id, u.name, u.email, u.created_at, u.updated_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, listQuery, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query comments list: %w", err)
	}
	defer rows.Close()

	comments := make([]*models.Comment, 0)
	for rows.Next() {
		comment := &models.Comment{Author: &models.User{}}
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author.ID,
			&comment.Author.Name,
			&comment.Author.Email,
			&comment.Author.CreatedAt,
			&comment.Author.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan comment row: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, totalItems, nil
}

func (r *mysqlCommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE comments
		SET content = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, comment.Content, comment.ID)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

func (r *mysqlCommentRepository) DeleteWithTx(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM comments WHERE id = ?`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id)
	} else {
		_, err = r.db.ExecContext(ctx, query, id)
	}
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

func (r *mysqlCommentRepository) DeleteByPostIDWithTx(ctx context.Context, tx *sql.Tx, postID int64) error {
	query := `DELETE FROM comments WHERE post_id = ?`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, postID)
	} else {
		_, err = r.db.ExecContext(ctx, query, postID)
	}
	if err != nil {
		return fmt.Errorf("failed to delete comments by post_id: %w", err)
	}
	return nil
}
