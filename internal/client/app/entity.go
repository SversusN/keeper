package app

// PasswordData – структура для типа данных Пароль.
type PasswordData struct {
	// Site – сайт, пароль от которого пользователь хочет сохранить.
	Site string `json:"site"`
	// Login – логин пользователь.
	Login string `json:"login"`
	// Password – пароль пользователя.
	Password string `json:"password"`
}

// CardData – структура для типа данных Карта.
type CardData struct {
	// Number – номер карты.
	Number string `json:"number"`
	// ExpDate – дата, до которой валидна карта.
	ExpDate string `json:"exp_date"`
	// CardHolder – держатель карты.
	CardHolder string `json:"card_holder"`
}

// FileData – структура для типа данных Файл.
type FileData struct {
	// Path – путь до файла.
	Path string `json:"path"`
	// Data – файл в бинарном представлении.
	Data []byte `json:"data"`
}

// TextData – структура для типа данных Текст.
type TextData struct {
	// Text – текст.
	Text string `json:"text"`
}
