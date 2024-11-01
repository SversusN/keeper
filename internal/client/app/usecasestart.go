package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/SversusN/keeper/internal/client/internalerrors"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// Start – функция для начала диалога.
// Пользователь авторизуется.
// Далее работает сценарий диалога во возможным веткам. //TODO покрыть больше ошибочных сценариев
func (c *Client) Start() error {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer cancelCtx()
	c.printer.PrintLogo()
	c.printer.Print("Здравствуйте и dont panic! Я ваш менеджер паролей.")
	c.printer.Print(fmt.Sprintf("Версия %s от %s", c.buildVersion, c.buildDate))

	if err := c.UserAuth(); err != nil {
		c.Logger.Log.Error(err)
		return err
	}
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return c.run(ctx)
	})

	return grp.Wait()
}

func (c *Client) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			c.Logger.Log.Info("client has been shutdown")
			return nil
		default:
			c.printer.Print("Выберите что хотели бы сделать (в консоли нужно ввести цифру)")
			fmt.Println("1. Получить информацию о всех ваших файлах")
			fmt.Println("2. Получить конкретную запись")
			fmt.Println("3. Сохранить новую запись")
			fmt.Println("4. Изменить запись")

			var commandNumber int
			_, err := c.printer.Scan(&commandNumber)
			if err != nil {
				fmt.Println("Введите целую цифру от 1 до 4")
			}

			switch commandNumber {
			case getUserDataList:
				err := c.GetUserDataList()
				if err != nil {
					if errors.Is(err, internalerrors.ErrNoData) {
						c.printer.Print("У вас нет записей. Создайте хотя бы одну")
						continue
					}
					c.Logger.Log.Error(err)
					continue
				}
			case getUserData:
				err := c.GetUserData()
				if err != nil {
					c.Logger.Log.Error(err)
					continue
				}
			case saveUserData:
				err := c.SaveData()
				if err != nil {
					c.Logger.Log.Error(err)
					continue
				}
			case editUserData:
				err := c.EditData()
				if err != nil {
					c.Logger.Log.Error(err)
					continue
				}
			default:
				fmt.Println("Вы сбились с курса.")
			}
			fmt.Printf("\n*******************\n")
		}
	}
}
