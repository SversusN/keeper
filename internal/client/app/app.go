package app

import (
	"sync/atomic"
	"time"

	"github.com/SversusN/keeper/internal/client/cache"
	"github.com/SversusN/keeper/internal/client/config"
	"github.com/SversusN/keeper/internal/client/grpcclient"
	"github.com/SversusN/keeper/internal/client/models"
	"github.com/SversusN/keeper/internal/client/utils"
	"github.com/SversusN/keeper/pkg/logger"
)

type printer interface {
	Print(s string)
	Scan(a ...interface{}) (int, error)
	PrintLogo()
}

type grpcClient interface {
	Register(model models.AuthModel) (models.AuthToken, error)
	SignIn(model models.AuthModel) (models.AuthToken, error)
	GetUserData(model models.UserDataModel) (*models.UserData, error)
	GetUserDataList() ([]models.UserDataList, error)
	SaveUserData(model *models.UserData) error
	UpdateUserData(model *models.UserData) error
	SyncUserData(ts int64) ([]models.UserDataList, error)
}

type clientCache interface {
	Append(data *models.UserData)
	GetUserData(model models.UserDataModel) (*models.UserData, error)
	GetUserDataList() []models.UserDataList
	GetMaxTS() (int64, error)
}

// isSignIn - глобальная переменная состояния чтобы не запускать синхронизацию до старта
var isSignIn atomic.Bool
var isSyncProces atomic.Bool

// Client – структура консольного клиента. Отвечает за сценарий приложения.
type Client struct {
	gRPCClient   grpcClient
	printer      printer
	cache        clientCache
	Config       *config.Config
	Logger       *logger.Logger
	dataSyncChan chan int64
	buildVersion string
	buildDate    string
}

// NewClient – функция для инициализации клиента.
func NewClient(l *logger.Logger, c *config.Config, bv string, bd string) (*Client, error) {
	gRPCClient, err := grpcclient.NewGRPCClient(c)
	if err != nil {
		return nil, err
	}
	localCache := cache.NewCache()
	dataSyncChan := make(chan int64, c.ChanSize)

	client := &Client{
		gRPCClient:   gRPCClient,
		printer:      &utils.Printer{},
		cache:        localCache,
		dataSyncChan: dataSyncChan,
		Config:       c,
		Logger:       l,
		buildVersion: bv,
		buildDate:    bd,
	}

	go client.startSync()

	return client, nil
}

func (c *Client) startSync() {
	ticker := time.NewTicker(time.Duration(c.Config.CashTimeRefresh) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case ID := <-c.dataSyncChan:
			model, err := c.gRPCClient.GetUserData(models.UserDataModel{ID: ID})
			if err != nil {
				c.Logger.Log.Errorf("get user data error: %v", err)
				continue
			}
			c.UpdateDataInCache(model)
			ticker.Reset(time.Duration(c.Config.CashTimeRefresh) * time.Second)
		case <-ticker.C:
			if !isSignIn.Load() || isSyncProces.Load() {
				continue //не запрашиваем если не залогинены или eщe идет синхронизация
			}
			isSyncProces.Store(true)
			maxTs, err := c.cache.GetMaxTS()
			if err != nil {
				c.Logger.Log.Errorf("get max ts error: %v", err)
				isSyncProces.Store(false)
				continue
			}
			result, err := c.gRPCClient.SyncUserData(maxTs)
			if err != nil {
				c.Logger.Log.Errorf("get user data error: %v", err)
				isSyncProces.Store(false)
				continue
			}
			for _, v := range result {
				data, err := c.gRPCClient.GetUserData(models.UserDataModel{ID: v.ID})
				if err != nil {
					c.Logger.Log.Errorf("get user data error: %v", err)
					isSyncProces.Store(false)
					continue
				}
				c.UpdateDataInCache(data)
				isSyncProces.Store(false)
			}
		default:
			continue
		}
	}
}
