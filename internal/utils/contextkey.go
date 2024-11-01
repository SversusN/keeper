package utils

type contextKey string

// Для реализации интерфейса строки по требованию линтера использовать типы
func (c contextKey) String() string {
	return "context key " + string(c)
}

// Требование линтера использовать типы
var UserIDContextKey = contextKey("user_id")
