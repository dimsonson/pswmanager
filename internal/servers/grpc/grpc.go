package grpc

import (
	"context"

	"net"
	"sync"

	"github.com/dimsonson/pswmanager/internal/models"
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
	//ShortServicePrivider
	//ShortService *ShortServices
}

// NewServer конструктор создания нового сервера в соответствии с существующей конфигурацией.
func NewServer(ctx context.Context, stop context.CancelFunc, cfg models.GRPC, wg *sync.WaitGroup) *Server {
	return &Server{
		Ctx:  ctx,
		Stop: stop,
		Wg:   wg,
		Cfg:  cfg,
	}
}

// Start метод запуска сервара, вид запвсукаемого сервера зависит от EnableGRPC в структуре Config.
func (srv *Server) Start() {

	//if srv.EnableGRPC {
	//srv.InitGRPC()
	srv.InitGRPCservice()
	//	srv.InitGRPC()
	srv.GrpcGracefullShotdown()
	srv.StartGRPC()
	//return
	//}
	//srv.InitHTTP()
	//srv.httpGracefullShotdown()
	//srv.StartHTTP()
}

// InitGRPC инциализация GRPC сервера.
func (srv *Server) InitGRPCservice() {
	// Инициализируем конструкторы.
	//srv.ShortService = &ShortServices{}
	// Конструктор хранилища.
	//s := newStrorageProvider(srv.DatabaseDsn, srv.FileStoragePath)
	// Конструкторы.
	//svcRand := &service.Rand{}
	//srv.ShortService.svsPut = service.NewPutService(s, srv.BaseURL, svcRand)
	// Конструктор Get слоя.
	//srv.ShortService.svsGet = service.NewGetService(s, srv.BaseURL)
	// Конструктор Delete слоя.
	//srv.ShortService.svsDel = service.NewDeleteService(s, srv.BaseURL)
	// Констуктор Ping слоя.
	//srv.ShortService.svsPing = service.NewPingService(s)
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
	//pb.RegisterShortServiceServer(srv.GRPCserver, srv.ShortService)
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
