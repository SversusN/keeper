package handlers

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"

	"github.com/SversusN/keeper/internal/server/config"
	"github.com/SversusN/keeper/internal/server/interceptors"
	"github.com/SversusN/keeper/internal/server/storage"
	"github.com/SversusN/keeper/internal/utils/encrypter"
	pb "github.com/SversusN/keeper/pkg/grpc"
	"github.com/SversusN/keeper/pkg/logger"
)

const (
	missingKeyErrText = "missing key in context"
)

type crypto interface {
	HashFunc(src string) (string, error)
	CompareHash(src, hash string) error
	BuildJWT(userID int64, secret string) (string, error)
	GetUserID(tokenString, secret string) (int64, error)
}

type repository interface {
	HealthCheck() error
	CreateUser(ctx context.Context, login, password string) (int64, error)
	FindUserByLogin(ctx context.Context, login string) (*storage.User, error)
	SaveUserData(ctx context.Context, userID int64, name, dataType string, data []byte) error
	GetUserData(ctx context.Context, userID int64) ([]storage.InfoRecord, error)
	FindUserRecord(ctx context.Context, id, userID int64) (*storage.Record, error)
	UpdateUserRecord(ctx context.Context, record *storage.Record) error
	GetUserDataForSync(ctx context.Context, userID int64, ts int64) ([]storage.InfoRecord, error)
}

// Server – сервер приложения, который отвечает за хранение и обработку приватных данных пользователя.
type Server struct {
	pb.UnimplementedKeeperServer
	Storage repository
	crypto  crypto
	Config  *config.Config
	Logger  *logger.Logger
}

// NewServer – создает объект сервера
// Функция принимает репозиторий, конфигуратор и логгер.
func NewServer(r repository, c *config.Config, l *logger.Logger) *Server {
	return &Server{
		Storage: r,
		Config:  c,
		Logger:  l,
		crypto:  &encrypter.Token{},
	}
}

// Start – метод для запуска сервера приложения.
func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.Config.Host)
	if err != nil {
		s.Logger.Log.Error(err)
		return fmt.Errorf("tcp connection failed")
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.AuthInterceptor(s.Logger, s.Config.SecretKey, s.crypto),
		logging.UnaryServerInterceptor(interceptors.LoggerInterceptor()),
	))

	pb.RegisterKeeperServer(gRPCServer, s)

	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer cancelCtx()

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer s.Logger.Log.Info("server has been shutdown")
		defer wg.Done()
		<-ctx.Done()

		s.Logger.Log.Info("app got a signal")
		gRPCServer.GracefulStop()
	}()

	s.Logger.Log.Info("gRPC server is running")

	return gRPCServer.Serve(listen)
}
