package repository

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepository struct {
	pool *pgxpool.Pool
}

func NewTagRepository(pool *pgxpool.Pool) *TagRepository {
	return &TagRepository{pool: pool}
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonSlugChars.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// SetForPoem replaces a poem's tags with the given names, upserting the tags.
func (r *TagRepository) SetForPoem(ctx context.Context, poemID uuid.UUID, names []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM poem_tags WHERE poem_id = $1`, poemID); err != nil {
		return err
	}

	seen := map[string]bool{}
	for _, name := range names {
		slug := slugify(name)
		if slug == "" || seen[slug] {
			continue
		}
		seen[slug] = true

		var tagID uuid.UUID
		err := tx.QueryRow(ctx,
			`INSERT INTO tags (id, name, slug) VALUES ($1, $2, $3)
			 ON CONFLICT (slug) DO UPDATE SET slug = EXCLUDED.slug
			 RETURNING id`,
			uuid.New(), strings.TrimSpace(name), slug,
		).Scan(&tagID)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO poem_tags (poem_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			poemID, tagID,
		); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// ListForPoem returns a poem's tag slugs.
func (r *TagRepository) ListForPoem(ctx context.Context, poemID uuid.UUID) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.slug FROM tags t JOIN poem_tags pt ON pt.tag_id = t.id
		 WHERE pt.poem_id = $1 ORDER BY t.slug`, poemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		tags = append(tags, s)
	}
	return tags, rows.Err()
}

// ListForPoems batch-loads tag slugs for many poems (avoids N+1 on list views).
func (r *TagRepository) ListForPoems(ctx context.Context, poemIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	result := make(map[uuid.UUID][]string)
	if len(poemIDs) == 0 {
		return result, nil
	}
	rows, err := r.pool.Query(ctx,
		`SELECT pt.poem_id, t.slug FROM poem_tags pt JOIN tags t ON t.id = pt.tag_id
		 WHERE pt.poem_id = ANY($1) ORDER BY t.slug`, poemIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pid uuid.UUID
		var slug string
		if err := rows.Scan(&pid, &slug); err != nil {
			return nil, err
		}
		result[pid] = append(result[pid], slug)
	}
	return result, rows.Err()
}
