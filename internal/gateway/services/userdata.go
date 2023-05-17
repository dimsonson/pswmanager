package services

import (
	"context"

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
	NewUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	NewApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error)
	ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error)
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	cfg        *config.ServiceConfig
	clientGRPC ClientGRPCProvider
}

// New.
func NewUserData(cfg *config.ServiceConfig, clientGRPC ClientGRPCProvider) *UserServices {
	return &UserServices{
		cfg:        cfg,
		clientGRPC: clientGRPC,
	}
}

// CreateUser.
func (s *UserServices) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	out, err := s.clientGRPC.NewUser(ctx, in)
	return out, err
}

// CreateApp .
func (s *UserServices) CreateApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	out, err := s.clientGRPC.NewApp(ctx, in)
	return out, err
}

// ReadUser.
func (s *UserServices) ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	out, err := s.clientGRPC.ReadUser(ctx, in)
	return out, err
}
