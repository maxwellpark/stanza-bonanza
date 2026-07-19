package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StanzaLikeRepository struct {
	pool *pgxpool.Pool
}

func NewStanzaLikeRepository(pool *pgxpool.Pool) *StanzaLikeRepository {
	return &StanzaLikeRepository{pool: pool}
}

// ToggleLike flips the like and adjusts stanzas.like_count in one transaction.
// Returns the resulting liked state.
func (r *StanzaLikeRepository) ToggleLike(ctx context.Context, userID, stanzaID uuid.UUID) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	del, err := tx.Exec(ctx, `DELETE FROM stanza_likes WHERE user_id = $1 AND stanza_id = $2`, userID, stanzaID)
	if err != nil {
		return false, err
	}
	if del.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE stanzas SET like_count = like_count - 1 WHERE id = $1`, stanzaID); err != nil {
			return false, err
		}
		return false, tx.Commit(ctx)
	}

	ins, err := tx.Exec(ctx,
		`INSERT INTO stanza_likes (user_id, stanza_id, created_at) VALUES ($1, $2, now())
		 ON CONFLICT (user_id, stanza_id) DO NOTHING`,
		userID, stanzaID,
	)
	if err != nil {
		return false, err
	}
	if ins.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE stanzas SET like_count = like_count + 1 WHERE id = $1`, stanzaID); err != nil {
			return false, err
		}
	}
	return true, tx.Commit(ctx)
}

// LikedStanzaIDs returns which of the given stanzas the user has liked.
func (r *StanzaLikeRepository) LikedStanzaIDs(ctx context.Context, userID uuid.UUID, stanzaIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	liked := make(map[uuid.UUID]bool)
	if userID == uuid.Nil || len(stanzaIDs) == 0 {
		return liked, nil
	}
	rows, err := r.pool.Query(ctx,
		`SELECT stanza_id FROM stanza_likes WHERE user_id = $1 AND stanza_id = ANY($2)`,
		userID, stanzaIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		liked[id] = true
	}
	return liked, rows.Err()
}
