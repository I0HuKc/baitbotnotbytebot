package sqlstore

import (
	"context"
	"database/sql"

	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
)

type DescRepository struct {
	store *sql.DB
}

func (r *DescRepository) Create(ctx context.Context, m *model.Desc) (err error) {
	return r.store.QueryRowContext(
		ctx,
		`
		INSERT INTO desctb(text)
		SELECT $1
		RETURNING id, text, created_at
		`,
		m.Text,
	).Scan(
		&m.Id,
		&m.Text,
		&m.CreatedAt,
	)
}

func (r *DescRepository) Get(ctx context.Context, m *model.Desc) (err error) {
	return r.store.QueryRowContext(
		ctx,
		`
		SELECT id, text, created_at
		FROM desctb
		WHERE id=$1
		`,
		m.Id,
	).Scan(
		&m.Id,
		&m.Text,
		&m.CreatedAt,
	)
}

func (r *DescRepository) Delete(ctx context.Context, m *model.Desc) (err error) {
	return r.store.QueryRowContext(
		ctx,
		`
		DELETE FROM desctb
		WHERE id=$1
		RETURNING id, text, created_at
		`,
		m.Id,
	).Scan(
		&m.Id,
		&m.Text,
		&m.CreatedAt,
	)
}

func (r *DescRepository) Count(ctx context.Context) (int, error) {
	var c int
	return c, r.store.QueryRowContext(
		ctx,
		`
		SELECT count(*)
		FROM desctb		
		`,
	).Scan(&c)
}

func (r *DescRepository) List(ctx context.Context) ([]*model.Desc, error) {
	arr := []*model.Desc{}

	rows, err := r.store.QueryContext(
		ctx,
		`
		SELECT id, text, created_at
		FROM desctb
		ORDER BY id DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := model.Desc{}
		if err := rows.Scan(
			&d.Id,
			&d.Text,
			&d.CreatedAt,
		); err != nil {
			continue
		}

		arr = append(arr, &d)
	}

	return arr, nil
}

func (r *DescRepository) FistLast(ctx context.Context) (int, int, error) {
	arr := []int{}

	rows, err := r.store.QueryContext(
		ctx,
		`
		(SELECT id
		FROM desctb
		ORDER BY id ASC
		LIMIT 1
		)
		
		UNION ALL

		(SELECT id
		FROM desctb
		ORDER BY id DESC
		LIMIT 1)
		`,
	)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			continue
		}

		arr = append(arr, v)
	}

	switch len(arr) {
	case 0:
		return 0, 0, nil
	case 1:
		return arr[0], arr[0], nil
	default:
		return arr[0], arr[1], nil
	}
}
