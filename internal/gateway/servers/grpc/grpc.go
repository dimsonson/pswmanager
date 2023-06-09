package grpc

import (
	"context"

	"net"
	"sync"

	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	pbconsume "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protoconsume"
	grpchandlers "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers"
	"github.com/dimsonson/pswmanager/internal/gateway/config"
	"github.com/dimsonson/pswmanager/internal/gateway/services"
	"github.com/dimsonson/pswmanager/pkg/log"

	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

)

// Server структура для хранения серверов.
type Server struct {
	GRPCserver        *grpc.Server
	Wg                *sync.WaitGroup
	Ctx               context.Context
	Stop              context.CancelFunc
	Cfg               config.GRPC
	UserService       *UserServices
	ClientRMQhandlers *ClientRMQhandlers
}

type UserServices struct {
	user *services.UserServices
	pb.UnimplementedUserServicesServer
	Server
}

type ClientRMQhandlers struct {
	pub *grpchandlers.ClientRMQhandlers
	pbpub.UnimplementedClientRMQhandlersServer
	Server
}

type ServerRMQhandlers struct {
	consume *grpchandlers.ServerRMQhandlers
	pbconsume.UnimplementedServerRMQhandlersServer
	Server
}

// NewServer конструктор создания нового сервера в соответствии с существующей конфигурацией.
func NewServer(ctx context.Context, stop context.CancelFunc, cfg config.GRPC, wg *sync.WaitGroup) *Server {
	return &Server{
		Ctx:  ctx,
		Stop: stop,
		Wg:   wg,
		Cfg:  cfg,
	}
}

// InitGRPC инциализация GRPC сервера.
func (srv *Server) InitGRPCservice(user *services.UserServices, clientRMQ *grpchandlers.ClientRMQhandlers, serverRMQ *grpchandlers.ServerRMQhandlers) {
	srv.UserService = &UserServices{user: user}
	srv.ClientRMQhandlers = &ClientRMQhandlers{pub: clientRMQ}
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
			logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(log.Logg)),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
	)
	pb.RegisterUserServicesServer(srv.GRPCserver, srv.UserService)
	pbpub.RegisterClientRMQhandlersServer(srv.GRPCserver, clientRMQ)
	pbconsume.RegisterServerRMQhandlersServer(srv.GRPCserver, serverRMQ)
}

// StartGRPC запуск GRPC сервера.
func (srv *Server) StartGRPC() {
	listen, err := net.Listen(srv.Cfg.ServerNetwork, srv.Cfg.ServerPort)
	if err != nil {
		log.Printf("gRPC listener error: %v", err)
	}
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

func (svs *UserServices) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	out, err := svs.user.CreateUser(ctx, in)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}

func (svs *UserServices) CreateApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	out, err := svs.user.CreateApp(ctx, in)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}

func (svs *UserServices) ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	out, err := svs.user.ReadUser(ctx, in)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}
