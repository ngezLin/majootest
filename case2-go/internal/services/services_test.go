package services

import (
	"context"
	"errors"
	"testing"

	"majootest/case2-go/internal/config"
	"majootest/case2-go/internal/models"
	"majootest/case2-go/internal/repositories"
)

func setupTestServices() (AuthService, PostService, CommentService, *repositories.MockUserRepository, *repositories.MockPostRepository, *repositories.MockCommentRepository) {
	cfg := &config.Config{
		JWTSecret:      "test-secret-key-12345",
		JWTExpiryHours: 24,
	}

	userRepo := repositories.NewMockUserRepository()
	postRepo := repositories.NewMockPostRepository()
	commentRepo := repositories.NewMockCommentRepository()
	txManager := &repositories.MockTxManager{}

	authSvc := NewAuthService(userRepo, cfg)
	postSvc := NewPostService(postRepo, commentRepo, txManager)
	commentSvc := NewCommentService(commentRepo, postRepo, txManager)

	return authSvc, postSvc, commentSvc, userRepo, postRepo, commentRepo
}

func TestAuthService(t *testing.T) {
	authSvc, _, _, _, _, _ := setupTestServices()
	ctx := context.Background()

	// 1. Register Success
	regReq := &models.RegisterRequest{
		Name:     "Alice Johnson",
		Email:    "alice@example.com",
		Password: "password123",
	}
	resp, err := authSvc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("expected successful register, got error: %v", err)
	}
	if resp.Token == "" || resp.User.ID == 0 {
		t.Errorf("expected non-empty token and user ID, got %v", resp)
	}

	// 2. Register Duplicate Email
	_, err = authSvc.Register(ctx, regReq)
	if !errors.Is(err, models.ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists on duplicate email, got: %v", err)
	}

	// 3. Login Success
	loginReq := &models.LoginRequest{
		Email:    "alice@example.com",
		Password: "password123",
	}
	loginResp, err := authSvc.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("expected successful login, got error: %v", err)
	}
	if loginResp.Token == "" {
		t.Errorf("expected non-empty login token")
	}

	// 4. Login Invalid Password
	badLoginReq := &models.LoginRequest{
		Email:    "alice@example.com",
		Password: "wrongpassword",
	}
	_, err = authSvc.Login(ctx, badLoginReq)
	if !errors.Is(err, models.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}

	// 5. GetProfile
	profile, err := authSvc.GetProfile(ctx, resp.User.ID)
	if err != nil || profile.Email != "alice@example.com" {
		t.Errorf("expected profile retrieval, got: %v, err: %v", profile, err)
	}
}

func TestPostService(t *testing.T) {
	_, postSvc, _, _, postRepo, _ := setupTestServices()
	ctx := context.Background()

	userID := int64(1)
	otherUserID := int64(2)

	// 1. Create Post
	createReq := &models.CreatePostRequest{
		Title:   "My First Blog Post",
		Content: "This is the complete content of the blog post.",
	}
	post, err := postSvc.Create(ctx, userID, createReq)
	if err != nil {
		t.Fatalf("expected successful post creation, got error: %v", err)
	}
	if post.ID == 0 || post.Title != createReq.Title {
		t.Errorf("unexpected post details: %v", post)
	}

	// 2. Get Post by ID
	fetched, err := postSvc.GetByID(ctx, post.ID)
	if err != nil || fetched.ID != post.ID {
		t.Errorf("expected fetched post, got: %v, err: %v", fetched, err)
	}

	// 3. List Posts
	posts, meta, err := postSvc.List(ctx, 1, 10, "")
	if err != nil || len(posts) != 1 || meta.TotalItems != 1 {
		t.Errorf("expected 1 post in list, got %d, meta: %v", len(posts), meta)
	}

	// 4. Update Post - Non-Author (Forbidden)
	newTitle := "Updated Title"
	updateReq := &models.UpdatePostRequest{
		Title: &newTitle,
	}
	_, err = postSvc.Update(ctx, post.ID, otherUserID, updateReq)
	if !errors.Is(err, models.ErrForbidden) {
		t.Errorf("expected ErrForbidden when non-owner updates, got: %v", err)
	}

	// 5. Update Post - Author (Success)
	updatedPost, err := postSvc.Update(ctx, post.ID, userID, updateReq)
	if err != nil || updatedPost.Title != newTitle {
		t.Errorf("expected updated title %s, got %v, err: %v", newTitle, updatedPost, err)
	}

	// 6. Delete Post - Non-Author (Forbidden)
	err = postSvc.Delete(ctx, post.ID, otherUserID)
	if !errors.Is(err, models.ErrForbidden) {
		t.Errorf("expected ErrForbidden when non-owner deletes, got: %v", err)
	}

	// 7. Delete Post - Author (Success)
	err = postSvc.Delete(ctx, post.ID, userID)
	if err != nil {
		t.Errorf("expected successful deletion, got: %v", err)
	}

	// Verify post is deleted
	_, err = postRepo.GetByID(ctx, post.ID)
	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got: %v", err)
	}
}

func TestCommentServiceAndTransactions(t *testing.T) {
	_, postSvc, commentSvc, _, postRepo, commentRepo := setupTestServices()
	ctx := context.Background()

	userID := int64(1)
	otherUserID := int64(2)

	// Create a post
	post, _ := postSvc.Create(ctx, userID, &models.CreatePostRequest{
		Title:   "Post for comments",
		Content: "Discussion post",
	})

	// 1. Create Comment on Non-Existent Post (Fails)
	_, err := commentSvc.Create(ctx, 9999, userID, &models.CreateCommentRequest{
		Content: "Hello on fake post",
	})
	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound for comment on fake post, got: %v", err)
	}

	// 2. Create Comment on Valid Post (Succeeds & Increments Comment Count Atomically)
	comment, err := commentSvc.Create(ctx, post.ID, otherUserID, &models.CreateCommentRequest{
		Content: "Great blog post! Really enjoyed reading it.",
	})
	if err != nil {
		t.Fatalf("expected successful comment create, got error: %v", err)
	}
	if comment.ID == 0 {
		t.Errorf("expected non-zero comment ID")
	}

	// Verify post comment count incremented
	p, _ := postRepo.GetByID(ctx, post.ID)
	if p.CommentCount != 1 {
		t.Errorf("expected comment_count to be 1, got %d", p.CommentCount)
	}

	// 3. List Comments
	comments, meta, err := commentSvc.ListByPostID(ctx, post.ID, 1, 10)
	if err != nil || len(comments) != 1 || meta.TotalItems != 1 {
		t.Errorf("expected 1 comment, got %d, meta: %v", len(comments), meta)
	}

	// 4. Update Comment - Non-Author (Forbidden)
	newCommentContent := "Updated comment text"
	_, err = commentSvc.Update(ctx, comment.ID, userID, &models.UpdateCommentRequest{
		Content: newCommentContent,
	})
	if !errors.Is(err, models.ErrForbidden) {
		t.Errorf("expected ErrForbidden for non-author comment update, got: %v", err)
	}

	// 5. Update Comment - Author (Success)
	updatedComment, err := commentSvc.Update(ctx, comment.ID, otherUserID, &models.UpdateCommentRequest{
		Content: newCommentContent,
	})
	if err != nil || updatedComment.Content != newCommentContent {
		t.Errorf("expected updated comment, got %v, err: %v", updatedComment, err)
	}

	// 6. Delete Comment - Author (Success & Decrements Comment Count Atomically)
	err = commentSvc.Delete(ctx, comment.ID, otherUserID)
	if err != nil {
		t.Errorf("expected successful comment delete, got: %v", err)
	}

	// Verify post comment count decremented
	p, _ = postRepo.GetByID(ctx, post.ID)
	if p.CommentCount != 0 {
		t.Errorf("expected comment_count to be 0 after delete, got %d", p.CommentCount)
	}

	// Verify comment deleted from repository
	_, err = commentRepo.GetByID(ctx, comment.ID)
	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound after comment delete, got: %v", err)
	}
}
