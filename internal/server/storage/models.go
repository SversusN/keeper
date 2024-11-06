package storage

// User – пользователь.
type User struct {
	// Login – login.
	Login string
	// Password – пароль.
	Password string
	// ID – идентификатор пользователя.
	ID int64
}

// InfoRecord – модель для информации о данных пользователя.
type InfoRecord struct {
	// Name – название для данных.
	Name string
	// DataType – тип данных.
	DataType string
	// CreatedAt – дата создания.
	CreatedAt string
	// ID – идентификатор данных.
	ID int64
	// Version – версия данных.
	Version int64
}

// Record – модель данных пользователя.
type Record struct {
	// Name – название для данных.
	Name string
	// DataType – тип данных.
	DataType string
	// CreatedAt – дата создания.
	CreatedAt string
	// Data – данные пользователя в байтах.
	Data []byte
	// ID – идентификатор данных.
	ID int64
	// Version – версия данных.
	Version int64
}
