package storage

import "errors"

var (
	ErrConflict          = errors.New(`already exists`)
	ErrNowRows           = errors.New(`missing data`)
	ErrConnectionRefused = errors.New(`connection refused`)
	ErrMigrationsFailed  = errors.New(`migrations failed`)
	ErrCreateUser        = errors.New(`create user error`)
	ErrFindUser          = errors.New(`find user error`)
	ErrSaveUserData      = errors.New(`save user data error`)
	ErrGetUserData       = errors.New(`get user data error`)
	ErrFindUserRecord    = errors.New(`find user record error`)
	ErrUpdateUserRecord  = errors.New(`update user record error`)
)
