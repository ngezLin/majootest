package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"majootest/case2-go/internal/models"
)

// PostRepository defines database operations for posts.
type PostRepository interface {
	Create(ctx context.Context, post *models.Post) error
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	List(ctx context.Context, limit, offset int, search string) ([]*models.Post, int64, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, id int64) error
	DeleteWithTx(ctx context.Context, tx *sql.Tx, id int64) error
	IncrementCommentCount(ctx context.Context, tx *sql.Tx, postID int64, delta int) error
}

type mysqlPostRepository struct {
	db *sql.DB
}

// NewPostRepository returns a new MySQL post repository instance.
func NewPostRepository(db *sql.DB) PostRepository {
	return &mysqlPostRepository{db: db}
}

func (r *mysqlPostRepository) Create(ctx context.Context, post *models.Post) error {
	query := `
		INSERT INTO posts (user_id, title, content, comment_count, created_at, updated_at)
		VALUES (?, ?, ?, 0, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query, post.UserID, post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	post.ID = id
	return nil
}

func (r *mysqlPostRepository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.comment_count, p.created_at, p.updated_at,
			u.id, u.name, u.email, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`
	post := &models.Post{Author: &models.User{}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CommentCount,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Author.ID,
		&post.Author.Name,
		&post.Author.Email,
		&post.Author.CreatedAt,
		&post.Author.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("failed to query post by id: %w", err)
	}
	return post, nil
}

func (r *mysqlPostRepository) List(ctx context.Context, limit, offset int, search string) ([]*models.Post, int64, error) {
	countQuery := `SELECT COUNT(*) FROM posts`
	listQuery := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.comment_count, p.created_at, p.updated_at,
			u.id, u.name, u.email, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
	`
	args := []interface{}{}
	countArgs := []interface{}{}

	if search != "" {
		filter := " WHERE p.title LIKE ? OR p.content LIKE ?"
		searchParam := "%" + search + "%"
		countQuery += " WHERE title LIKE ? OR content LIKE ?"
		countArgs = append(countArgs, searchParam, searchParam)
		listQuery += filter
		args = append(args, searchParam, searchParam)
	}

	var totalItems int64
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	listQuery += " ORDER BY p.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		post := &models.Post{Author: &models.User{}}
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CommentCount,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author.ID,
			&post.Author.Name,
			&post.Author.Email,
			&post.Author.CreatedAt,
			&post.Author.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan post row: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, totalItems, nil
}

func (r *mysqlPostRepository) Update(ctx context.Context, post *models.Post) error {
	query := `
		UPDATE posts
		SET title = ?, content = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (r *mysqlPostRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

func (r *mysqlPostRepository) DeleteWithTx(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post within tx: %w", err)
	}
	return nil
}

func (r *mysqlPostRepository) IncrementCommentCount(ctx context.Context, tx *sql.Tx, postID int64, delta int) error {
	query := `
		UPDATE posts
		SET comment_count = GREATEST(0, comment_count + ?), updated_at = NOW()
		WHERE id = ?
	`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, delta, postID)
	} else {
		_, err = r.db.ExecContext(ctx, query, delta, postID)
	}
	if err != nil {
		return fmt.Errorf("failed to update comment count: %w", err)
	}
	return nil
}
