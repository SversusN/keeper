package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/SversusN/keeper/internal/client/models"
)

// Cache – потокобезопасный Map для хранения данных внутри клиента
type Cache struct {
	mem *sync.Map
}

var (
	ErrNotFound = errors.New(`not found in cache`)
)

// NewCache – функция инициализации кеша.
func NewCache() *Cache {
	return &Cache{
		mem: &sync.Map{},
	}
}

// Append – метод добавления данных в кеш.
func (c *Cache) Append(model *models.UserData) {
	c.mem.Store(model.ID, model)
}

// GetUserData – метод получения данных из кеша.
func (c *Cache) GetUserData(model models.UserDataModel) (*models.UserData, error) {
	d, ok := c.mem.Load(model.ID)
	if !ok {
		return nil, ErrNotFound
	}

	return d.(*models.UserData), nil
}

// GetUserDataList – метод получения информации о данных пользователя из кеша.
func (c *Cache) GetUserDataList() []models.UserDataList {
	records := make([]models.UserDataList, 0)
	c.mem.Range(func(k, v interface{}) bool {
		rec := models.UserDataList{
			Name:     v.(*models.UserData).Name,
			DataType: v.(*models.UserData).DataType,
			ID:       k.(int64),
			Version:  v.(*models.UserData).Version,
		}
		records = append(records, rec)

		return true
	})

	return records
}

// GetMaxTS получает максимальную дату информации пользователя
func (c *Cache) GetMaxTS() (int64, error) {
	var maxDate int64
	maxDate = 0
	c.mem.Range(func(k, v interface{}) bool {
		temp, err := time.Parse(time.DateTime, v.(*models.UserData).CreatedAt)
		if err != nil {
			return false
		}
		if temp.Unix() > maxDate {
			maxDate = temp.Unix()
		}
		return true
	})
	return maxDate, nil
}
