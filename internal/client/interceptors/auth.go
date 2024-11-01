package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type client interface {
	GetAuthToken() string
}

// AuthInterceptor – перехватчик запросов клиента, которые добавляет auth_token клиента в запрос.
func AuthInterceptor(c client) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{},
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "token", c.GetAuthToken())
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
