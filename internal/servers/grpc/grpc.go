package grpc

import (
	"context"

	"net"
	"sync"

	"github.com/dimsonson/pswmanager/internal/handlers/grpchandle"
	pb "github.com/dimsonson/pswmanager/internal/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/dimsonson/pswmanager/internal/services"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server структура для хранения серверов.
type Server struct {
	GRPCserver *grpc.Server
	Wg         *sync.WaitGroup
	Ctx        context.Context
	Stop       context.CancelFunc
	Cfg        models.GRPC
	Handler    *grpchandle.UserServices
	ServicesGRPC
	//ShortServicePrivider
	//ShortService *ShortServices
}

type ServicesGRPC struct {
	User      *services.UserServices
	ReadUsers *services.ReadUserServices
	pb.UnimplementedUserServicesServer
}

// NewServer конструктор создания нового сервера в соответствии с существующей конфигурацией.
func NewServer(ctx context.Context, stop context.CancelFunc, cfg models.GRPC, hn *grpchandle.UserServices, wg *sync.WaitGroup) *Server {
	return &Server{
		Ctx:  ctx,
		Stop: stop,
		Wg:   wg,
		Cfg:  cfg,
		Handler: hn,
	}
}

// Start метод запуска сервара, вид запвсукаемого сервера зависит от EnableGRPC в структуре Config.
func (srv *Server) Start() {
	srv.InitGRPCservice()
	srv.GrpcGracefullShotdown()
	srv.StartGRPC()
}

// InitGRPC инциализация GRPC сервера.
func (srv *Server) InitGRPCservice() {

	// Обявление customFunc для использования в обработке паники.
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	// Опции для логгера и восстановления после паники.
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	// создаём gRPC-сервер без зарегистрированной службы
	srv.GRPCserver = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(log.Logger)),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
	)
}

// StartGRPC запуск GRPC сервера.
func (srv *Server) StartGRPC() {
	listen, err := net.Listen(srv.Cfg.Network, srv.Cfg.Port)
	if err != nil {
		log.Printf("gRPC listener error: %v", err)
	}



	pb.RegisterUserServicesServer(srv.GRPCserver, srv.ServicesGRPC)
	log.Print("gRPCServer ListenAndServe starting listening")
	// запуск gRPC сервера
	if err := srv.GRPCserver.Serve(listen); err != nil {
		log.Printf("gRPC server error: %v", err)
	}
	log.Print("grpc server shutting down...")
	srv.Wg.Done()
}

// grpcGracefullShotdown метод благопроиятного для соединений и незавершенных запросов закрытия сервера.
func (srv *Server) GrpcGracefullShotdown() {
	go func() {
		// получаем сигнал о завершении приложения
		<-srv.Ctx.Done()
		log.Print("grpcSrv got signal, attempting graceful shutdown")
		srv.GRPCserver.GracefulStop()
		srv.Wg.Done()
	}()
}
