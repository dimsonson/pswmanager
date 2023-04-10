package grpc

import (
	"context"

	"net"
	"sync"

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
	GRPCserver  *grpc.Server
	Wg          *sync.WaitGroup
	Ctx         context.Context
	Stop        context.CancelFunc
	Cfg         models.GRPC
	UserService *UserServices
}

type ReadUserServicesProvider interface {
	ReadUser(ctx context.Context, uid string) (models.SetRecords, error)
}

type UserServicesProvider interface {
	CreateUser(ctx context.Context, login string, psw string) (models.UserConfig, error)
	CreateApp(ctx context.Context, uid string, psw string) (string, models.UserConfig, error)
}

type UserServices struct {
	user *services.UserServices
	read *services.ReadUserServices
	pb.UnimplementedUserServicesServer
	Server
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

// InitGRPC инциализация GRPC сервера.
func (srv *Server) InitGRPCservice(readUser *services.ReadUserServices, user *services.UserServices) {
	srv.UserService = &UserServices{}
	srv.UserService.read = readUser
	srv.UserService.user = user
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

	pb.RegisterUserServicesServer(srv.GRPCserver, srv.UserService)
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

func (s *UserServices) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var out pb.CreateUserResponse
	var err error
	usercfg, err := s.user.CreateUser(ctx, in.Login, in.Psw)
	out.UserID = usercfg.UserID
	out.RmqHost = usercfg.RmqHost
	out.RmqPort = usercfg.RmqPort
	out.RmqUID = usercfg.RmqUID
	out.ExchangeName = usercfg.ExchangeName
	out.ExchangeKind = usercfg.ExchangeKind

	for _, v := range usercfg.Apps {
		strout := &pb.App{
			AppID:        v.AppID,
			RoutingKey:   v.RoutingKey,
			ConsumeQueue: v.ConsumeQueue,
			ConsumerName: v.ConsumerName,
		}
		out.Apps = append(out.Apps, strout)
	}

	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}

	return &out, err
}

func (s *UserServices) CreateApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	var out pb.CreateAppResponse
	var err error
	appid, usercfg, err := s.user.CreateApp(ctx, in.Uid, in.Psw)
	out.Appid = appid
	out.UserID = usercfg.UserID
	out.RmqHost = usercfg.RmqHost
	out.RmqPort = usercfg.RmqPort
	out.RmqUID = usercfg.RmqUID
	out.ExchangeName = usercfg.ExchangeName
	out.ExchangeKind = usercfg.ExchangeKind

	for _, v := range usercfg.Apps {
		strout := &pb.App{
			AppID:        v.AppID,
			RoutingKey:   v.RoutingKey,
			ConsumeQueue: v.ConsumeQueue,
			ConsumerName: v.ConsumerName,
		}
		out.Apps = append(out.Apps, strout)
	}

	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}

	return &out, err
}

func (s *UserServices) ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	var out pb.ReadUserResponse
	var err error
	setrecords, err := s.read.ReadUser(ctx, in.Uid)

	for _, v := range setrecords.SetTextRec {
		strout := &pb.TextRec{
			RecordID: v.RecordID,
			//ChngTime: v.ChngTime,
			UID:      v.UID,
			AppID:    v.AppID,
			Text:     v.Text,
			Metadata: v.Metadata,
		}
		out.SetTextRec = append(out.SetTextRec, strout)
	}

	for _, v := range setrecords.SetBinaryRec {
		strout := &pb.BinaryRec{
			RecordID: v.RecordID,
			//ChngTime: v.ChngTime,
			UID:      v.UID,
			AppID:    v.AppID,
			Binary:   v.Binary,
			Metadata: v.Metadata,
		}
		out.SetBinaryRec = append(out.SetBinaryRec, strout)
	}

	for _, v := range setrecords.SetLoginRec {
		strout := &pb.LoginRec{
			RecordID: v.RecordID,
			//ChngTime: v.ChngTime,
			UID:      v.UID,
			AppID:    v.AppID,
			Login:    v.Login,
			Psw:      v.Psw,
			Metadata: v.Metadata,
		}
		out.SetLoginRec = append(out.SetLoginRec, strout)
	}

	for _, v := range setrecords.SetCardRec {
		strout := &pb.CardRec{
			RecordID: v.RecordID,
			//ChngTime: v.ChngTime,
			UID:       v.UID,
			AppID:     v.AppID,
			Brand:     int64(v.Brand),
			Number:    v.Number,
			ValidDate: v.ValidDate,
			Code:      int64(v.Code),
			Holder:    v.Holder,
			Metadata:  v.Metadata,
		}
		out.SetCardRec = append(out.SetCardRec, strout)
	}

	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}

	return &out, err
}
