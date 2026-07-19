package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/repository"
)

type poemStore interface {
	Create(ctx context.Context, poem *domain.Poem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Poem, error)
	List(ctx context.Context, page domain.PaginationParams, format, sort, q, tag string) ([]domain.Poem, int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error)
	ListFeed(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error)
	ListExplore(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error)
	ListHallOfFame(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error)
	Update(ctx context.Context, poem *domain.Poem) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementCounter(ctx context.Context, id uuid.UUID, column string, delta int) error
}

type stanzaStore interface {
	Create(ctx context.Context, stanza *domain.Stanza) error
	ListByPoem(ctx context.Context, poemID uuid.UUID) ([]domain.Stanza, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.StanzaStatus) error
	GetNextPosition(ctx context.Context, poemID uuid.UUID) (int, error)
}

type poemNotifStore interface {
	Create(ctx context.Context, notif *domain.Notification) error
}

type likeChecker interface {
	Exists(ctx context.Context, userID, poemID uuid.UUID) (bool, error)
}

type tagStore interface {
	SetForPoem(ctx context.Context, poemID uuid.UUID, names []string) error
	ListForPoem(ctx context.Context, poemID uuid.UUID) ([]string, error)
	ListForPoems(ctx context.Context, poemIDs []uuid.UUID) (map[uuid.UUID][]string, error)
}

type stanzaLikeChecker interface {
	LikedStanzaIDs(ctx context.Context, userID uuid.UUID, stanzaIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

type PoemService struct {
	poems       poemStore
	stanzas     stanzaStore
	notifs      poemNotifStore
	likes       likeChecker
	tags        tagStore
	stanzaLikes stanzaLikeChecker
}

func NewPoemService(poems *repository.PoemRepository, stanzas *repository.StanzaRepository, notifs *repository.NotificationRepository, likes *repository.LikeRepository, tags *repository.TagRepository, stanzaLikes *repository.StanzaLikeRepository) *PoemService {
	return &PoemService{
		poems:       poems,
		stanzas:     stanzas,
		notifs:      notifs,
		likes:       likes,
		tags:        tags,
		stanzaLikes: stanzaLikes,
	}
}

// attachTags fills Tags on a batch of poems with one query.
func (s *PoemService) attachTags(ctx context.Context, poems []domain.Poem) {
	if s.tags == nil || len(poems) == 0 {
		return
	}
	ids := make([]uuid.UUID, len(poems))
	for i := range poems {
		ids[i] = poems[i].ID
	}
	byID, err := s.tags.ListForPoems(ctx, ids)
	if err != nil {
		return
	}
	for i := range poems {
		poems[i].Tags = byID[poems[i].ID]
	}
}

func (s *PoemService) Create(ctx context.Context, userID uuid.UUID, title, description string, format domain.PoemFormat, approvalMode domain.ApprovalMode, maxStanzas *int, tags []string) (*domain.Poem, error) {
	var poem = &domain.Poem{
		AuthorID:     userID,
		Title:        title,
		Description:  description,
		Format:       format,
		ApprovalMode: approvalMode,
		MaxStanzas:   maxStanzas,
	}

	if err := s.poems.Create(ctx, poem); err != nil {
		return nil, fmt.Errorf("creating poem: %w", err)
	}

	if len(tags) > 0 && s.tags != nil {
		if err := s.tags.SetForPoem(ctx, poem.ID, tags); err != nil {
			return nil, fmt.Errorf("setting tags: %w", err)
		}
		poem.Tags, _ = s.tags.ListForPoem(ctx, poem.ID)
	}

	return poem, nil
}

func (s *PoemService) Get(ctx context.Context, id, viewerID uuid.UUID) (*domain.Poem, error) {
	poem, err := s.poems.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("poem not found: %w", err)
	}

	stanzas, err := s.stanzas.ListByPoem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("loading stanzas: %w", err)
	}
	poem.Stanzas = stanzas

	if viewerID != uuid.Nil && s.likes != nil {
		if liked, err := s.likes.Exists(ctx, viewerID, id); err == nil {
			poem.LikedByMe = liked
		}
	}

	if viewerID != uuid.Nil && s.stanzaLikes != nil && len(poem.Stanzas) > 0 {
		ids := make([]uuid.UUID, len(poem.Stanzas))
		for i := range poem.Stanzas {
			ids[i] = poem.Stanzas[i].ID
		}
		if liked, err := s.stanzaLikes.LikedStanzaIDs(ctx, viewerID, ids); err == nil {
			for i := range poem.Stanzas {
				poem.Stanzas[i].LikedByMe = liked[poem.Stanzas[i].ID]
			}
		}
	}

	if s.tags != nil {
		poem.Tags, _ = s.tags.ListForPoem(ctx, id)
	}

	return poem, nil
}

func (s *PoemService) List(ctx context.Context, page domain.PaginationParams, format, sort, q, tag string) ([]domain.Poem, int, error) {
	poems, total, err := s.poems.List(ctx, page, format, sort, q, tag)
	if err != nil {
		return nil, 0, err
	}
	s.attachTags(ctx, poems)
	return poems, total, nil
}

func (s *PoemService) ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	poems, total, err := s.poems.ListByUser(ctx, userID, page)
	if err != nil {
		return nil, 0, err
	}
	s.attachTags(ctx, poems)
	return poems, total, nil
}

func (s *PoemService) Feed(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	poems, total, err := s.poems.ListFeed(ctx, userID, page)
	if err != nil {
		return nil, 0, err
	}
	s.attachTags(ctx, poems)
	return poems, total, nil
}

func (s *PoemService) Explore(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	poems, total, err := s.poems.ListExplore(ctx, page)
	if err != nil {
		return nil, 0, err
	}
	s.attachTags(ctx, poems)
	return poems, total, nil
}

func (s *PoemService) HallOfFame(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	return s.poems.ListHallOfFame(ctx, page)
}

func (s *PoemService) Update(ctx context.Context, userID, poemID uuid.UUID, title, description string, tags []string) error {
	poem, err := s.poems.GetByID(ctx, poemID)
	if err != nil {
		return fmt.Errorf("poem not found: %w", err)
	}
	if poem.AuthorID != userID {
		return fmt.Errorf("not the poem author")
	}

	poem.Title = title
	poem.Description = description
	if err := s.poems.Update(ctx, poem); err != nil {
		return err
	}

	if tags != nil && s.tags != nil {
		if err := s.tags.SetForPoem(ctx, poemID, tags); err != nil {
			return fmt.Errorf("setting tags: %w", err)
		}
	}
	return nil
}

func (s *PoemService) Delete(ctx context.Context, userID, poemID uuid.UUID) error {
	poem, err := s.poems.GetByID(ctx, poemID)
	if err != nil {
		return fmt.Errorf("poem not found: %w", err)
	}
	if poem.AuthorID != userID {
		return fmt.Errorf("not the poem author")
	}

	return s.poems.Delete(ctx, poemID)
}

func (s *PoemService) ListStanzas(ctx context.Context, poemID uuid.UUID) ([]domain.Stanza, error) {
	return s.stanzas.ListByPoem(ctx, poemID)
}

func (s *PoemService) SubmitStanza(ctx context.Context, userID, poemID uuid.UUID, text, literaryDevice string) (*domain.Stanza, error) {
	poem, err := s.poems.GetByID(ctx, poemID)
	if err != nil {
		return nil, fmt.Errorf("poem not found: %w", err)
	}

	if poem.ApprovalMode == domain.ApprovalClosed && poem.AuthorID != userID {
		return nil, fmt.Errorf("poem is closed for submissions")
	}

	if poem.MaxStanzas != nil && poem.StanzaCount >= *poem.MaxStanzas {
		return nil, fmt.Errorf("poem has reached the maximum number of stanzas")
	}

	nextPos, err := s.stanzas.GetNextPosition(ctx, poemID)
	if err != nil {
		return nil, fmt.Errorf("getting next position: %w", err)
	}

	var status domain.StanzaStatus
	if poem.ApprovalMode == domain.ApprovalRequired && poem.AuthorID != userID {
		status = domain.StanzaPending
	} else {
		status = domain.StanzaApproved
	}

	var stanza = &domain.Stanza{
		PoemID:         poemID,
		AuthorID:       userID,
		Text:           text,
		Position:       nextPos,
		LiteraryDevice: literaryDevice,
		Status:         status,
	}

	if err := s.stanzas.Create(ctx, stanza); err != nil {
		return nil, fmt.Errorf("creating stanza: %w", err)
	}

	if status == domain.StanzaApproved {
		_ = s.poems.IncrementCounter(ctx, poemID, "stanza_count", 1)
	}

	if poem.AuthorID != userID {
		_ = s.notifs.Create(ctx, &domain.Notification{
			RecipientID: poem.AuthorID,
			ActorID:     &userID,
			Type:        domain.NotifStanzaSubmit,
			PoemID:      &poemID,
		})
	}

	return stanza, nil
}

func (s *PoemService) ReviewStanza(ctx context.Context, userID, poemID, stanzaID uuid.UUID, approved bool) error {
	poem, err := s.poems.GetByID(ctx, poemID)
	if err != nil {
		return fmt.Errorf("poem not found: %w", err)
	}
	if poem.AuthorID != userID {
		return fmt.Errorf("not the poem author")
	}

	// The stanza must belong to this poem and still be pending, else an author
	// could review a foreign stanza and re-approving would double-count.
	stanzas, err := s.stanzas.ListByPoem(ctx, poemID)
	if err != nil {
		return fmt.Errorf("loading stanzas: %w", err)
	}
	var target *domain.Stanza
	for i := range stanzas {
		if stanzas[i].ID == stanzaID {
			target = &stanzas[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("stanza does not belong to this poem")
	}
	if target.Status != domain.StanzaPending {
		return fmt.Errorf("stanza is not pending review")
	}

	var status domain.StanzaStatus
	var notifType domain.NotificationType
	if approved {
		status = domain.StanzaApproved
		notifType = domain.NotifStanzaApproved
	} else {
		status = domain.StanzaRejected
		notifType = domain.NotifStanzaRejected
	}

	if err := s.stanzas.UpdateStatus(ctx, stanzaID, status); err != nil {
		return fmt.Errorf("updating stanza status: %w", err)
	}
	if approved {
		_ = s.poems.IncrementCounter(ctx, poemID, "stanza_count", 1)
	}

	if target.AuthorID != userID {
		_ = s.notifs.Create(ctx, &domain.Notification{
			RecipientID: target.AuthorID,
			ActorID:     &userID,
			Type:        notifType,
			PoemID:      &poemID,
		})
	}

	return nil
}
