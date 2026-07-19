package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/repository"
)

type likeStore interface {
	ToggleLike(ctx context.Context, userID, poemID uuid.UUID) (bool, error)
	Exists(ctx context.Context, userID, poemID uuid.UUID) (bool, error)
}

type stanzaLikeStore interface {
	ToggleLike(ctx context.Context, userID, stanzaID uuid.UUID) (bool, error)
}

type commentStore interface {
	Create(ctx context.Context, comment *domain.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByPoem(ctx context.Context, poemID uuid.UUID, page domain.PaginationParams) ([]domain.Comment, int, error)
}

type followStore interface {
	Create(ctx context.Context, follow *domain.Follow) error
	Delete(ctx context.Context, followerID, followedID uuid.UUID) error
	Exists(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)
	ListFollowers(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error)
	ListFollowing(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error)
}

type socialNotifStore interface {
	Create(ctx context.Context, notif *domain.Notification) error
	ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Notification, int, error)
	MarkRead(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
}

type socialPoemStore interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Poem, error)
	IncrementCounter(ctx context.Context, id uuid.UUID, column string, delta int) error
}

type SocialService struct {
	likes       likeStore
	stanzaLikes stanzaLikeStore
	comments    commentStore
	follows     followStore
	notifs      socialNotifStore
	poems       socialPoemStore
}

func NewSocialService(
	likes *repository.LikeRepository,
	stanzaLikes *repository.StanzaLikeRepository,
	comments *repository.CommentRepository,
	follows *repository.FollowRepository,
	notifs *repository.NotificationRepository,
	poems *repository.PoemRepository,
) *SocialService {
	return &SocialService{
		likes:       likes,
		stanzaLikes: stanzaLikes,
		comments:    comments,
		follows:     follows,
		notifs:      notifs,
		poems:       poems,
	}
}

func (s *SocialService) ToggleStanzaLike(ctx context.Context, userID, stanzaID uuid.UUID) (bool, error) {
	liked, err := s.stanzaLikes.ToggleLike(ctx, userID, stanzaID)
	if err != nil {
		return false, fmt.Errorf("toggling stanza like: %w", err)
	}
	return liked, nil
}

func (s *SocialService) ToggleLike(ctx context.Context, userID, poemID uuid.UUID) (bool, error) {
	liked, err := s.likes.ToggleLike(ctx, userID, poemID)
	if err != nil {
		return false, fmt.Errorf("toggling like: %w", err)
	}
	if !liked {
		return false, nil
	}

	poem, err := s.poems.GetByID(ctx, poemID)
	if err == nil && poem.AuthorID != userID {
		_ = s.notifs.Create(ctx, &domain.Notification{
			RecipientID: poem.AuthorID,
			ActorID:     &userID,
			Type:        domain.NotifLike,
			PoemID:      &poemID,
		})
	}

	return true, nil
}

func (s *SocialService) AddComment(ctx context.Context, userID, poemID uuid.UUID, parentID *uuid.UUID, text string) (*domain.Comment, error) {
	var comment = &domain.Comment{
		PoemID:   poemID,
		AuthorID: userID,
		ParentID: parentID,
		Text:     text,
	}

	if err := s.comments.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("creating comment: %w", err)
	}

	poem, err := s.poems.GetByID(ctx, poemID)
	if err == nil && poem.AuthorID != userID {
		_ = s.notifs.Create(ctx, &domain.Notification{
			RecipientID: poem.AuthorID,
			ActorID:     &userID,
			Type:        domain.NotifComment,
			PoemID:      &poemID,
		})
	}

	return comment, nil
}

func (s *SocialService) DeleteComment(ctx context.Context, userID, commentID uuid.UUID) error {
	comment, err := s.comments.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}
	if comment.AuthorID != userID {
		return fmt.Errorf("not the comment author")
	}

	if err := s.comments.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}

	return nil
}

func (s *SocialService) ListComments(ctx context.Context, poemID uuid.UUID, page domain.PaginationParams) ([]domain.Comment, int, error) {
	return s.comments.ListByPoem(ctx, poemID, page)
}

func (s *SocialService) ToggleFollow(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	if followerID == followedID {
		return false, fmt.Errorf("cannot follow yourself")
	}

	exists, err := s.follows.Exists(ctx, followerID, followedID)
	if err != nil {
		return false, fmt.Errorf("checking follow: %w", err)
	}

	if exists {
		if err := s.follows.Delete(ctx, followerID, followedID); err != nil {
			return false, fmt.Errorf("unfollowing: %w", err)
		}
		return false, nil
	}

	var follow = &domain.Follow{
		FollowerID: followerID,
		FollowedID: followedID,
	}
	if err := s.follows.Create(ctx, follow); err != nil {
		return false, fmt.Errorf("following: %w", err)
	}

	_ = s.notifs.Create(ctx, &domain.Notification{
		RecipientID: followedID,
		ActorID:     &followerID,
		Type:        domain.NotifFollow,
	})

	return true, nil
}

func (s *SocialService) ListFollowers(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	return s.follows.ListFollowers(ctx, userID, page)
}

func (s *SocialService) ListFollowing(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	return s.follows.ListFollowing(ctx, userID, page)
}

func (s *SocialService) ListNotifications(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Notification, int, error) {
	return s.notifs.ListByUser(ctx, userID, page)
}

func (s *SocialService) MarkNotificationsRead(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	return s.notifs.MarkRead(ctx, userID, ids)
}
