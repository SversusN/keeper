package app

import "github.com/SversusN/keeper/internal/client/models"

// SyncCache - запускается при получении информации о пользователе и пишет в буферизированный канал ID записи
// saveDataToCache запускающий параллельную горутину скачки
func (c *Client) SyncCache(records []models.UserDataList) {
	for _, rec := range records {
		cachedData, err := c.cache.GetUserData(models.UserDataModel{ID: rec.ID})
		if err != nil || cachedData.Version != rec.Version {
			c.saveDataToCache(rec.ID)
		}
	}
}

func (c *Client) saveDataToCache(dataID int64) {
	c.dataSyncChan <- dataID
}

// Добавляем в кэш
func (c *Client) UpdateDataInCache(data *models.UserData) {
	c.cache.Append(data)
}

// GetDataFromCache - получение из кэш
func (c *Client) GetDataFromCache(m models.UserDataModel) (*models.UserData, error) {
	c.Logger.Log.Debug("load data from cache...")
	return c.cache.GetUserData(m)
}

func (c *Client) GetDataListFromCache() []models.UserDataList {
	c.Logger.Log.Debug("load data list from cache...")
	return c.cache.GetUserDataList()
}
