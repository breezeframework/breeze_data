package db

type DBConnectionPool interface {
	GetConnection() DBConnection
	Close() error
}
