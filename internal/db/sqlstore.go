package db

type SqlStore interface {
	Desc() DescRepositoryI
}
