package suite

import (
	"context"
	ssov1 "github.com/Korjick/service-go-protos/gen/go/sso"
	"github.com/Korjick/sso-service-go/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const (
	configPath = "../config/local.yaml"
	grpcHost   = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()   // Функция будет восприниматься как вспомогательная для тестов
	t.Parallel() // Разрешаем параллельный запуск тестов

	// Читаем конфиг из файла
	cfg := config.MustLoadByPath(configPath)

	// Основной родительский контекст
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	// Когда тесты пройдут, закрываем контекст
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// Адрес нашего gRPC-сервера
	grpcAddress := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))

	// Создаем клиент
	cc, err := grpc.NewClient(grpcAddress,
		// Используем insecure-коннект для тестов
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	// gRPC-клиент сервера Auth
	authClient := ssov1.NewAuthClient(cc)

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authClient,
	}
}
