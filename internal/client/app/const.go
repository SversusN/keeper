package app

const (
	passwordData = iota + 1
	cardData
	fileData
	textData
)

const (
	passwordDataType = "password"
	cardDataType     = "card"
	fileDataType     = "file"
	textDataType     = "text"
	loginInput       = "Пользователь: "
	passwordInput    = "Пароль: "
	siteInput        = "Ресурс: "
	cardNumberInput  = "Номер карты: "
	cardHolderInput  = "Держатель карты: "
	cardExpDateInput = "Дата истечения срока (мм/гг): "
)

const (
	getUserDataList = iota + 1
	getUserData
	saveUserData
	editUserData
)

const InternalErrTemplate = "%w: something went wrong: %w"
