package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/gateway/config"
	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type UserStorageProviver interface {
	// Close()
	// CreateUser(ctx context.Context, login string, psw string, uid string, usercfg config.UserConfig) error
	// ReadUserCfg(ctx context.Context, uid string) (config.UserConfig, error)
	// UpdateUser(ctx context.Context, uid string, usercfg config.UserConfig) error
	// CheckPsw(ctx context.Context, uid string, psw string) (bool, error)
	// IsUserLoginExist(ctx context.Context, login string) (bool, error)
}

type ClientRMQProvider interface {
	Close()
	ExchangeDeclare(exchName string) error
	QueueDeclare(queueName string) (models.Queue, error)
	QueueBind(queueName string, routingKey string) error
	UserInit() (config.UserConfig, *config.App)
	AppInit(usercfg config.UserConfig) config.App
}

type ClientGRPCProvider interface {
	NewUser(ctx context.Context, c pb.UserServicesClient, login string, psw string) (*pb.CreateUserResponse, error)
	NewApp(ctx context.Context, c pb.UserServicesClient, uid string, psw string) (*pb.CreateAppResponse, error)
	ReadUser(ctx context.Context, c pb.UserServicesClient, newAppCfg *pb.CreateAppResponse) (*pb.ReadUserResponse, error)
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	//storage    UserStorageProviver
	cfg        *config.ServiceConfig
	AppCfg     *pb.CreateAppResponse
	clientRMQ  ClientRMQProvider
	clientGRPC ClientGRPCProvider
}

// New.
func NewUserData(cfg *config.ServiceConfig, /*clientrmq ClientRMQProvider,*/ clientgrpc ClientGRPCProvider) *UserServices {
	return &UserServices{
		cfg:        cfg,
	//	clientRMQ:  clientrmq,
		clientGRPC: clientgrpc,
	}
}

// CreateUser.
func (s *UserServices) CreateUser(ctx context.Context, login string, psw string) (*pb.CreateUserResponse, error) {
	c := pb.NewUserServicesClient(s.cfg.GRPC.ClientConn)
	out, err := s.clientGRPC.NewUser(ctx, c, login, psw)
	if err != nil {
		log.Printf("gRPC create user service error: %v", err)
	}
	return out, err
}

// CreateUser получаем psw хешированный base64 и .
func (s *UserServices) CreateApp(ctx context.Context, uid string, psw string) (*pb.CreateAppResponse, error) {
	c := pb.NewUserServicesClient(s.cfg.GRPC.ClientConn)
	out, err := s.clientGRPC.NewApp(ctx, c, uid, psw)
	if err != nil {
		log.Printf("gRPC create app service error: %v", err)
	}
	s.AppCfg = out
	return out, err
}

func (s *UserServices) ReadUser(ctx context.Context, uid string) (*pb.ReadUserResponse, error) {
	c := pb.NewUserServicesClient(s.cfg.GRPC.ClientConn)
	out, err := s.clientGRPC.ReadUser(ctx, c, s.AppCfg)
	if err != nil {
		log.Printf("gRPC read user service error: %v", err)
	}
	return out, err
}
