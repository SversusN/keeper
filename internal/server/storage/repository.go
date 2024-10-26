package storage

import "context"

type Repository interface {
	HealthCheck() error
	CreateUser(ctx context.Context, login, password string) (int64, error)
	FindUserByLogin(ctx context.Context, login string) (*User, error)
	SaveUserData(ctx context.Context, userID int64, name, dataType string, data []byte) error
	GetUserData(ctx context.Context, userID int64) ([]InfoRecord, error)
	FindUserRecord(ctx context.Context, id, userID int64) (*Record, error)
	UpdateUserRecord(ctx context.Context, record *Record) error
}
