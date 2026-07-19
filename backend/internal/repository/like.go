package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type LikeRepository struct {
	pool *pgxpool.Pool
}

func NewLikeRepository(pool *pgxpool.Pool) *LikeRepository {
	return &LikeRepository{pool: pool}
}

func (r *LikeRepository) Create(ctx context.Context, like *domain.Like) error {
	like.CreatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO likes (user_id, poem_id, created_at) VALUES ($1, $2, $3)`,
		like.UserID, like.PoemID, like.CreatedAt,
	)
	return err
}

func (r *LikeRepository) Delete(ctx context.Context, userID, poemID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM likes WHERE user_id = $1 AND poem_id = $2`, userID, poemID)
	return err
}

// ToggleLike flips the like and adjusts poems.like_count in one transaction so
// concurrent requests can't drift the count. Returns the resulting liked state.
func (r *LikeRepository) ToggleLike(ctx context.Context, userID, poemID uuid.UUID) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	del, err := tx.Exec(ctx, `DELETE FROM likes WHERE user_id = $1 AND poem_id = $2`, userID, poemID)
	if err != nil {
		return false, err
	}
	if del.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE poems SET like_count = like_count - 1 WHERE id = $1`, poemID); err != nil {
			return false, err
		}
		return false, tx.Commit(ctx)
	}

	ins, err := tx.Exec(ctx,
		`INSERT INTO likes (user_id, poem_id, created_at) VALUES ($1, $2, now())
		 ON CONFLICT (user_id, poem_id) DO NOTHING`,
		userID, poemID,
	)
	if err != nil {
		return false, err
	}
	if ins.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE poems SET like_count = like_count + 1 WHERE id = $1`, poemID); err != nil {
			return false, err
		}
	}
	return true, tx.Commit(ctx)
}

func (r *LikeRepository) Exists(ctx context.Context, userID, poemID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND poem_id = $2)`,
		userID, poemID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
