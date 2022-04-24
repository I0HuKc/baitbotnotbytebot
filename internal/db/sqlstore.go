package db

type SqlStore interface {
	Desc() DescRepositoryI
	ChangeDesc() ChangeDescRepositoryI
	Performance() PerformanceRepositoryI
}
