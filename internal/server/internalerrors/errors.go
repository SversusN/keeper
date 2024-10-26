package internalerrors

import "errors"

var (
	ErrConflict          = errors.New(`already exists`)
	ErrNoRows            = errors.New(`no data`)
	ErrConnectionRefused = errors.New(`connection refused`)
	ErrCreateUser        = errors.New(`create user error`)
	ErrFindUser          = errors.New(`find user error`)
	ErrSaveUserData      = errors.New(`save user data error`)
	ErrGetUserData       = errors.New(`get user data error`)
	ErrFindUserRecord    = errors.New(`find user record error`)
	ErrUpdateUserRecord  = errors.New(`update user record error`)
	ErrMigrationsFailed  = errors.New(`migrations failed`)
)
