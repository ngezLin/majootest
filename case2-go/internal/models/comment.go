package models

import "time"

// Comment represents a blog post comment domain entity.
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	Author    *User     `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCommentRequest holds payload data for creating a new comment.
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// UpdateCommentRequest holds payload data for updating an existing comment.
type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}
