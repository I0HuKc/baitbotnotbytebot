package sqlstore

import (
	"context"
	"database/sql"

	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
)

type ChangeDescRepository struct {
	store *sql.DB
}

func (r *ChangeDescRepository) Create(ctx context.Context, m *model.ChangeDesc) error {
	return r.store.QueryRowContext(
		ctx,
		`
		INSERT INTO changedesc(groupid, groupname, nextdescchange)
		SELECT $1, $2, $3
		RETURNING id, groupid, groupname, nextdescchange, created_at
		`,
		m.GroupId,
		m.GroupName,
		m.NextDescChange,
	).Scan(
		&m.Id,
		&m.GroupId,
		&m.GroupName,
		&m.NextDescChange,
		&m.CreatedAt,
	)
}

func (r *ChangeDescRepository) Delete(ctx context.Context, m *model.ChangeDesc) error {
	return r.store.QueryRowContext(
		ctx,
		`
		DELETE FROM changedesc
		WHERE id=$1
		RETURNING id, groupid, groupname, nextdescchange, created_at
		`,
		m.Id,
	).Scan(
		&m.Id,
		&m.GroupId,
		&m.GroupName,
		&m.NextDescChange,
		&m.CreatedAt,
	)
}

func (r *ChangeDescRepository) List(ctx context.Context) ([]*model.ChangeDesc, error) {
	arr := []*model.ChangeDesc{}

	rows, err := r.store.QueryContext(
		ctx,
		`
		SELECT id, groupid, groupname, nextdescchange, created_at
		FROM changedesc
		ORDER BY id DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := model.ChangeDesc{}
		if err := rows.Scan(
			&p.Id,
			&p.GroupId,
			&p.GroupName,
			&p.NextDescChange,
			&p.CreatedAt,
		); err != nil {
			continue
		}

		arr = append(arr, &p)
	}

	return arr, nil
}
