package interceptors

import (
	"context"
	"github.com/SversusN/keeper/internal/utils"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/SversusN/keeper/pkg/grpc"
	"github.com/SversusN/keeper/pkg/logger"
)

type crypt interface {
	GetUserID(string, string) (int64, error)
}

// AuthInterceptor – перехватчик сервера для проверки авторизации пользователя.
func AuthInterceptor(l *logger.Logger, secret string, cr crypt) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, r interface{}, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
		if i.FullMethod == keeper.Keeper_Register_FullMethodName ||
			i.FullMethod == keeper.Keeper_SignIn_FullMethodName ||
			i.FullMethod == keeper.Keeper_Ping_FullMethodName {
			return h(ctx, r)
		}

		var token string
		if meta, ok := metadata.FromIncomingContext(ctx); ok {
			values := meta.Get("token")
			if len(values) > 0 {
				token = values[0]
			}
		}
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, http.StatusText(http.StatusForbidden))
		}

		userID, err := cr.GetUserID(token, secret)
		if err != nil {
			l.Log.Debugf("invalid token: %v", token)
			return nil, status.Error(codes.Unauthenticated, http.StatusText(http.StatusForbidden))
		}

		ctx = context.WithValue(ctx, utils.UserIDContextKey, userID)

		return h(ctx, r)
	}
}
