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

type ChangeDescRepositoryI interface {
	Create(ctx context.Context, m *model.ChangeDesc) (err error)
	Delete(ctx context.Context, m *model.ChangeDesc) (err error)
	List(ctx context.Context) ([]*model.ChangeDesc, error)
}

type PerformanceRepositoryI interface {
	Create(ctx context.Context, m *model.Performance) (err error)
	Delete(ctx context.Context, m *model.Performance) (err error)
	List(ctx context.Context) ([]*model.Performance, error)
	GetByGroupId(ctx context.Context, m *model.Performance) error
	Update(ctx context.Context, m *model.Performance) error
}
