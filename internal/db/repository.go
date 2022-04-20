package db

import (
	"context"

	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
)

type DescRepositoryI interface {
	Create(ctx context.Context, m *model.Desc) (err error)
	Get(ctx context.Context, m *model.Desc) (err error)
	Delete(ctx context.Context, m *model.Desc) (err error)
	Count(ctx context.Context) (int, error)
	FistLast(ctx context.Context) (int, int, error)
	List(ctx context.Context) ([]*model.Desc, error)
}
