package app

// PasswordData – структура для типа данных Пароль.
type PasswordData struct {
	// Site – сайт, пароль от которого пользователь хочет сохранить.
	Site string `json:"site"`
	// Login – логин пользователь.
	Login string `json:"login"`
	// Password – пароль пользователя.
	Password string `json:"password"`
	// MetaInfo - метаинформация для пароля
	MetaInfo string `json:"meta_info"`
}

// CardData – структура для типа данных Карта.
type CardData struct {
	// Number номер карточки .
	Number string `json:"number"`
	// ExpDate – дата, до которой валидна карта.
	ExpDate string `json:"exp_date"`
	// CardHolder – держатель карты.
	CardHolder string `json:"card_holder"`
	// MetaInfo - метаинформация для пароля
	MetaInfo string `json:"meta_info"`
}

// FileData – структура для типа данных Файл.
type FileData struct {
	// Path – путь до файла.
	Path string `json:"path"`
	// Data – файл в бинарном представлении.
	Data []byte `json:"data"`
	// MetaInfo - метаинформация для файла
	MetaInfo string `json:"meta_info"`
}

// TextData – структура для типа данных Текст.
type TextData struct {
	// Text – текст.
	Text string `json:"text"`
	// MetaInfo - метаинформация для записи
	MetaInfo string `json:"meta_info"`
}
