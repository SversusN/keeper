package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SversusN/keeper/internal/client/internalerrors"

	"github.com/SversusN/keeper/internal/client/models"
)

// GetUserData – получение сохранённых данных.
// Сценарий:
//  1. Приложение забирает с сервера все сохранённые данные пользователя (мета данные).
//  2. Приложение предлагает пользователю выбрать из сохранённых данных те, которые пользователь хочет получить.
//  3. Приложение достаёт нужные данные и отдаёт пользователю.
func (c *Client) GetUserData() error {
	err := c.GetUserDataList()
	if err != nil {
		if errors.Is(err, internalerrors.ErrNoData) {
			c.printer.Print("У вас пока нет записей")
			return nil
		}
		return err
	}

	c.printer.Print("Введите идентификатор записи")
	var (
		data   *models.UserData
		dataID int64
	)
	_, err = c.printer.Scan(&dataID)
	if err != nil {
		return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	m := models.UserDataModel{ID: dataID}
	data, err = c.GetDataFromCache(m) // сначала пытаемся достать из кеша
	if err != nil {
		c.Logger.Log.Warnf("get data from cache error: %w", err)
		data, err = c.gRPCClient.GetUserData(m) // если в кеше нет, идём на сервер
		if err != nil {
			return err
		}
		c.cache.Append(data) // складываем в кеш
	}

	err = printData(data)
	if err != nil {
		return err
	}

	return nil
}

func printData(data *models.UserData) error {
	var pretty []byte
	dataType := data.DataType
	switch dataType {
	case passwordDataType:
		passStruct := &PasswordData{}
		err := json.Unmarshal(data.Data, passStruct)
		if err != nil {
			return fmt.Errorf("json unmarshal error for struct PasswordData: %w", err)
		}
		pretty, err = json.MarshalIndent(passStruct, "", "  ")
		if err != nil {
			return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
		}
	case cardDataType:
		cardStruct := &CardData{}
		err := json.Unmarshal(data.Data, cardStruct)
		if err != nil {
			return fmt.Errorf("json unmarshal error for struct CardData: %w", err)
		}
		pretty, err = json.MarshalIndent(cardStruct, "", "  ")
		if err != nil {
			return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
		}
	case fileDataType:
		fileStruct := &FileData{}
		err := json.Unmarshal(data.Data, fileStruct)
		if err != nil {
			return fmt.Errorf("json unmarshal error for struct FileData: %w", err)
		}
		pretty, err = json.MarshalIndent(fileStruct, "", "  ")
		if err != nil {
			return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
		}
	case textDataType:
		textStruct := &TextData{}
		err := json.Unmarshal(data.Data, textStruct)
		if err != nil {
			return fmt.Errorf("json unmarshal error for struct TextData: %w", err)
		}
		pretty, err = json.MarshalIndent(textStruct, "", "  ")
		if err != nil {
			return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
		}
	default:
		return nil
	}
	fmt.Printf("\nВаша запись:\n%s", pretty)

	return nil
}

// GetUserDataList – получение информации о всех данных пользователя.
func (c *Client) GetUserDataList() error {
	records, err := c.gRPCClient.GetUserDataList()
	if err != nil { // что-то с сервером
		c.Logger.Log.Warnf("get user data list error: %v", err)
		records = c.GetDataListFromCache() // пытаемся достать из кеша
	} else {
		c.SyncCache(records)
	}
	if len(records) == 0 {
		return internalerrors.ErrNoData
	}

	c.printer.Print("Ваши сохраненные записи:")
	for _, el := range records {
		fmt.Printf("id: %d, name: %s, type: %s, version: %d\n", el.ID, el.Name, el.DataType, el.Version)
	}

	return nil
}
