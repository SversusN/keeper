package app

import (
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
}

type clientCache interface {
	Append(data *models.UserData)
	GetUserData(model models.UserDataModel) (*models.UserData, error)
	GetUserDataList() []models.UserDataList
}

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
	for {
		select {
		case ID := <-c.dataSyncChan:
			model, err := c.gRPCClient.GetUserData(models.UserDataModel{ID: ID})
			if err != nil {
				c.Logger.Log.Errorf("get user data error: %v", err)
				continue
			}
			c.UpdateDataInCache(model)
		default:
			continue
		}
	}
}
