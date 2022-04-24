package sqlstore

import (
	"context"
	"database/sql"

	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
)

type PerformanceRepository struct {
	store *sql.DB
}

func (r *PerformanceRepository) Create(ctx context.Context, m *model.Performance) error {
	return r.store.QueryRowContext(
		ctx,
		`
		INSERT INTO performance(groupid, groupname, nextjoke)
		SELECT $1, $2, $3
		RETURNING id, groupid, groupname, nextjoke, created_at
		`,
		m.GroupId,
		m.GroupName,
		m.NextJoke,
	).Scan(
		&m.Id,
		&m.GroupId,
		&m.GroupName,
		&m.NextJoke,
		&m.CreatedAt,
	)
}

func (r *PerformanceRepository) GetByGroupId(ctx context.Context, m *model.Performance) error {
	return r.store.QueryRowContext(
		ctx,
		`
		SELECT id, groupid, groupname, nextjoke, created_at
		FROM performance
		WHERE groupid=$1
		`,
		m.GroupId,
	).Scan(
		&m.Id,
		&m.GroupId,
		&m.GroupName,
		&m.NextJoke,
		&m.CreatedAt,
	)
}

func (r *PerformanceRepository) Update(ctx context.Context, m *model.Performance) error {
	return r.store.QueryRowContext(
		ctx,
		`
		UPDATE performance
		SET nextjoke=$1
		WHERE groupid=$2
		RETURNING id, groupid, groupname, nextjoke, created_at
		`,
		m.NextJoke,
		m.GroupId,
	).Scan(
		&m.Id,
		&m.GroupId,
		&m.GroupName,
		&m.NextJoke,
		&m.CreatedAt,
	)
}

func (r *PerformanceRepository) Delete(ctx context.Context, m *model.Performance) error {
	return r.store.QueryRowContext(
		ctx,
		`
		DELETE FROM performance
		WHERE groupid=$1
		RETURNING id, groupid, groupname, nextjoke, created_at
		`,
		m.GroupId,
	).Scan(
		&m.Id,
		&m.GroupId,
		&m.GroupName,
		&m.NextJoke,
		&m.CreatedAt,
	)
}

func (r *PerformanceRepository) List(ctx context.Context) ([]*model.Performance, error) {
	arr := []*model.Performance{}

	rows, err := r.store.QueryContext(
		ctx,
		`
		SELECT id, groupid, groupname, nextjoke, created_at
		FROM performance
		ORDER BY id DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := model.Performance{}
		if err := rows.Scan(
			&p.Id,
			&p.GroupId,
			&p.GroupName,
			&p.NextJoke,
			&p.CreatedAt,
		); err != nil {
			continue
		}

		arr = append(arr, &p)
	}

	return arr, nil
}
