package handlers

import (
	"context"
	"errors"
	errors2 "github.com/SversusN/keeper/internal/server/internalerrors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/SversusN/keeper/pkg/grpc"
)

// Register – метод регистрации пользователя на сервере.
func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	login := in.Login
	pass := in.Password
	encryptedPass, err := s.crypto.HashFunc(pass)
	if err != nil {
		s.Logger.Log.Error(err)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	userID, err := s.Storage.CreateUser(ctx, login, encryptedPass)
	if err != nil {
		if errors.Is(err, errors2.ErrConflict) {
			s.Logger.Log.Debug("oops, user already register")
			return nil, status.Error(codes.AlreadyExists, http.StatusText(http.StatusConflict))
		}
		s.Logger.Log.Error(err)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	authToken, err := s.crypto.BuildJWT(userID, s.Config.SecretKey)
	if err != nil {
		s.Logger.Log.Errorf("update user token error: %v", err)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return &pb.RegisterResponse{
		Token: authToken,
	}, nil
}

// SignIn – метод аутентификации на сервере.
func (s *Server) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	login := in.Login
	pass := in.Password
	user, err := s.Storage.FindUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, errors2.ErrNoRows) {
			s.Logger.Log.Debugf("user with login `%v` not found", login)
			return nil, status.Error(codes.NotFound, http.StatusText(http.StatusUnauthorized))
		}
		s.Logger.Log.Error(err)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	if err := s.crypto.CompareHash(pass, user.Password); err != nil {
		s.Logger.Log.Debugf("authorization failed, login: %v", login)
		return nil, status.Error(codes.Unauthenticated, http.StatusText(http.StatusUnauthorized))
	}

	authToken, err := s.crypto.BuildJWT(user.ID, s.Config.SecretKey)
	if err != nil {
		s.Logger.Log.Errorf("update user token error: %v", err)
		return nil, status.Error(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}
	return &pb.SignInResponse{
		Token: authToken,
	}, nil
}
