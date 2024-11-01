package app

import (
	"fmt"
	"github.com/SversusN/keeper/internal/client/internalerrors"

	"github.com/SversusN/keeper/internal/client/models"
)

// UserAuth – функция авторизации пользователя.
func (c *Client) UserAuth() error {
	var answer string
	c.printer.Print("Вы уже зарегистрированы? (y/n)")
	_, err := c.printer.Scan(&answer)
	if err != nil {
		return fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	switch answer {
	case "y":
		authM, err := buildAuthData(c.printer)
		if err != nil {
			return err

		}
		return c.signIn(*authM)
	case "n":
		authM, err := buildAuthData(c.printer)
		if err != nil {
			return err
		}
		return c.register(*authM)
	default:
		return c.UserAuth()
	}
}

func (c *Client) signIn(authM models.AuthModel) error {
	_, err := c.gRPCClient.SignIn(authM)
	if err != nil {
		fmt.Printf("%w: SignIn error: %w \n", internalerrors.ErrUserNotAuthorized, err)
		c.printer.Print("Ошибка! попробуйте еще раз...")
		return c.UserAuth()
	}
	return nil
}

func (c *Client) register(authM models.AuthModel) error {
	_, err := c.gRPCClient.Register(authM)
	if err != nil {
		fmt.Printf("%w: Register error: %w \n", internalerrors.ErrUserNotAuthorized, err)
		c.printer.Print("Ошибка! попробуйте еще раз...")
		return c.UserAuth()
	}
	return nil
}

func buildAuthData(p printer) (*models.AuthModel, error) {
	var (
		login, password string
		err             error
	)

	p.Print("Введите Ваши логин\\пароль")
	fmt.Print(loginInput)
	_, err = p.Scan(&login)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}
	fmt.Print(passwordInput)
	_, err = p.Scan(&password)
	if err != nil {
		return nil, fmt.Errorf(InternalErrTemplate, internalerrors.ErrInternal, err)
	}

	return &models.AuthModel{
		Login:    login,
		Password: password,
	}, nil
}
