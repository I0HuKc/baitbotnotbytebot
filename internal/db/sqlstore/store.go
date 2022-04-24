package sqlstore

import (
	"database/sql"

	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
)

type Store struct {
	db *sql.DB

	descRepository   *DescRepository
	perfRepository   *PerformanceRepository
	chDescRepository *ChangeDescRepository
}

func (s *Store) Desc() db.DescRepositoryI {
	if s.descRepository != nil {
		return s.descRepository
	}

	s.descRepository = &DescRepository{store: s.db}
	return s.descRepository
}

func (s *Store) ChangeDesc() db.ChangeDescRepositoryI {
	if s.chDescRepository != nil {
		return s.chDescRepository
	}

	s.chDescRepository = &ChangeDescRepository{store: s.db}
	return s.chDescRepository
}

func (s *Store) Performance() db.PerformanceRepositoryI {
	if s.perfRepository != nil {
		return s.perfRepository
	}

	s.perfRepository = &PerformanceRepository{store: s.db}
	return s.perfRepository
}

func CreateSqlStore(db *sql.DB) db.SqlStore {
	return &Store{
		db: db,
	}
}
