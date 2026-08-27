package models

import "time"

// Post represents the blog post domain entity.
type Post struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CommentCount int       `json:"comment_count"`
	Author       *User     `json:"author,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreatePostRequest holds payload data for creating a new post.
type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required,min=5"`
}

// UpdatePostRequest holds payload data for updating an existing post.
type UpdatePostRequest struct {
	Title   *string `json:"title" validate:"omitempty,min=3,max=255"`
	Content *string `json:"content" validate:"omitempty,min=5"`
}
