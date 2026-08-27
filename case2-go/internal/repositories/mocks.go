package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"majootest/case2-go/internal/models"
)

// MockTxManager is an in-memory transaction manager for testing
type MockTxManager struct {
	ShouldFail bool
}

func (m *MockTxManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	if m.ShouldFail {
		return fmt.Errorf("transaction failed")
	}
	return fn(nil)
}

// MockUserRepository is an in-memory mock for UserRepository
type MockUserRepository struct {
	Users       map[int64]*models.User
	UsersByMail map[string]*models.User
	NextID      int64
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:       make(map[int64]*models.User),
		UsersByMail: make(map[string]*models.User),
		NextID:      1,
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if _, exists := m.UsersByMail[user.Email]; exists {
		return models.ErrUserAlreadyExists
	}
	user.ID = m.NextID
	m.NextID++
	m.Users[user.ID] = user
	m.UsersByMail[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	u, ok := m.Users[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return u, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u, ok := m.UsersByMail[email]
	if !ok {
		return nil, models.ErrNotFound
	}
	return u, nil
}

// MockPostRepository is an in-memory mock for PostRepository
type MockPostRepository struct {
	Posts  map[int64]*models.Post
	NextID int64
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		Posts:  make(map[int64]*models.Post),
		NextID: 1,
	}
}

func (m *MockPostRepository) Create(ctx context.Context, post *models.Post) error {
	post.ID = m.NextID
	m.NextID++
	m.Posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	p, ok := m.Posts[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return p, nil
}

func (m *MockPostRepository) List(ctx context.Context, limit, offset int, search string) ([]*models.Post, int64, error) {
	var list []*models.Post
	for _, p := range m.Posts {
		list = append(list, p)
	}
	return list, int64(len(list)), nil
}

func (m *MockPostRepository) Update(ctx context.Context, post *models.Post) error {
	if _, ok := m.Posts[post.ID]; !ok {
		return models.ErrNotFound
	}
	m.Posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Delete(ctx context.Context, id int64) error {
	delete(m.Posts, id)
	return nil
}

func (m *MockPostRepository) DeleteWithTx(ctx context.Context, tx *sql.Tx, id int64) error {
	return m.Delete(ctx, id)
}

func (m *MockPostRepository) IncrementCommentCount(ctx context.Context, tx *sql.Tx, postID int64, delta int) error {
	p, ok := m.Posts[postID]
	if !ok {
		return models.ErrNotFound
	}
	p.CommentCount += delta
	return nil
}

// MockCommentRepository is an in-memory mock for CommentRepository
type MockCommentRepository struct {
	Comments map[int64]*models.Comment
	NextID   int64
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		Comments: make(map[int64]*models.Comment),
		NextID:   1,
	}
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	return m.CreateWithTx(ctx, nil, comment)
}

func (m *MockCommentRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, comment *models.Comment) error {
	comment.ID = m.NextID
	m.NextID++
	m.Comments[comment.ID] = comment
	return nil
}

func (m *MockCommentRepository) GetByID(ctx context.Context, id int64) (*models.Comment, error) {
	c, ok := m.Comments[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return c, nil
}

func (m *MockCommentRepository) ListByPostID(ctx context.Context, postID int64, limit, offset int) ([]*models.Comment, int64, error) {
	var list []*models.Comment
	for _, c := range m.Comments {
		if c.PostID == postID {
			list = append(list, c)
		}
	}
	return list, int64(len(list)), nil
}

func (m *MockCommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	if _, ok := m.Comments[comment.ID]; !ok {
		return models.ErrNotFound
	}
	m.Comments[comment.ID] = comment
	return nil
}

func (m *MockCommentRepository) DeleteWithTx(ctx context.Context, tx *sql.Tx, id int64) error {
	delete(m.Comments, id)
	return nil
}

func (m *MockCommentRepository) DeleteByPostIDWithTx(ctx context.Context, tx *sql.Tx, postID int64) error {
	for id, c := range m.Comments {
		if c.PostID == postID {
			delete(m.Comments, id)
		}
	}
	return nil
}
