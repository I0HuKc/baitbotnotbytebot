package sqlstore

import (
	"database/sql"

	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
)

type Store struct {
	db *sql.DB

	descRepository *DescRepository
}

func (s *Store) Desc() db.DescRepositoryI {
	if s.descRepository != nil {
		return s.descRepository
	}

	s.descRepository = &DescRepository{store: s.db}
	return s.descRepository
}

func CreateSqlStore(db *sql.DB) db.SqlStoreI {
	return &Store{
		db: db,
	}
}
