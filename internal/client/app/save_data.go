package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SversusN/keeper/internal/client/internalerrors"
	"io"
	"os"
	"strings"

	"github.com/SversusN/keeper/internal/client/models"
)

var dataTypes = [4]string{passwordDataType, cardDataType, fileDataType, textDataType}

// SaveData – сохранение данных пользователя.
func (c *Client) SaveData() error {
	var (
		dti  int
		name string
	)
	c.printer.Print("Что Вы хотели бы сохранить?")
	for i, dt := range dataTypes {
		fmt.Printf("%v. %v\n", i+1, dt)
	}
	_, err := c.printer.Scan(&dti)
	if err != nil {
		return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	model, err := buildData(dti, c.printer)
	if err != nil {
		return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	c.printer.Print("Введите строковой идентификатор записи?")
	_, err = c.printer.Scan(&name)
	if err != nil {
		return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	model.Name = name

	err = c.gRPCClient.SaveUserData(model)
	if err != nil {
		return err
	}
	c.UpdateDataInCache(model)

	c.printer.Print("Сохранено!")

	return nil
}

func buildData(dti int, p printer) (*models.UserData, error) {
	switch dti {
	case passwordData:
		return buildPassword(p)
	case cardData:
		return buildCardData(p)
	case fileData:
		return buildFileData(p)
	case textData:
		return buildTextData(p)
	default:
		return nil, internalerrors.ErrUnknownDataType
	}
}

//nolint:dupl // it's builder
func buildPassword(p printer) (*models.UserData, error) {
	p.Print("Введите ресурс логин и пароль")
	pass := &PasswordData{}
	fmt.Print(siteInput)
	_, err := p.Scan(&pass.Site)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(loginInput)
	_, err = p.Scan(&pass.Login)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(passwordInput)
	_, err = p.Scan(&pass.Password)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(metaPass)
	_, err = p.Scan(&pass.MetaInfo)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	bd, err := json.Marshal(pass)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	return &models.UserData{
		DataType: passwordDataType,
		Data:     bd,
	}, nil
}

//nolint:dupl // it's builder
func buildCardData(p printer) (*models.UserData, error) {
	p.Print("Внесите информацию о банковской карте")
	card := &CardData{}
	fmt.Print(cardNumberInput)
	_, err := p.Scan(&card.Number)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(cardExpDateInput)
	_, err = p.Scan(&card.ExpDate)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(cardHolderInput)
	_, err = p.Scan(&card.CardHolder)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(metaCard)
	_, err = p.Scan(&card.MetaInfo)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	//	byteData, err := easyjson.Marshal(card)
	bd, err := json.Marshal(card)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	return &models.UserData{
		DataType: cardDataType,
		Data:     bd,
	}, nil
}

func buildTextData(p printer) (*models.UserData, error) {
	text := &TextData{}
	p.Print("Введите текст заметки")
	//Нужно учесть пробелы
	in := bufio.NewReader(os.Stdin)
	t, err := in.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(metaText)
	_, err = p.Scan(&text.MetaInfo)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	text.Text = strings.TrimSpace(t)
	bd, err := json.Marshal(text)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	return &models.UserData{
		DataType: textDataType,
		Data:     bd,
	}, nil
}

func buildFileData(p printer) (*models.UserData, error) {
	file := &FileData{}
	p.Print("Введите полный путь до файла в Вашей операционной системе")
	_, err := p.Scan(&file.Path)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	openedFile, err := os.Open(file.Path)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	defer func() {
		err = openedFile.Close()
		if err != nil {
			fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
		}
	}()

	stat, err := openedFile.Stat()
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(openedFile).Read(bs)
	if err != nil && errors.Is(err, io.EOF) {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	file.Path = stat.Name()
	file.Data = bs

	bd, err := json.Marshal(file)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(metaFile)
	_, err = p.Scan(&file.MetaInfo)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	return &models.UserData{
		DataType: fileDataType,
		Data:     bd,
	}, nil
}
