package services

import (
	"context"

	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
)

type ClientGRPCProvider interface {
	NewUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	NewApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error)
	ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error)
	IsOnline() bool 
}

// Services структура конструктора бизнес логики.
// type UserServices struct {
// 	cfg        *config.ServiceConfig
// 	clientGRPC ClientGRPCProvider
// }

// New.
func NewUserData1(cfg *config.ServiceConfig, clientGRPC ClientGRPCProvider) *UserServices {
	return &UserServices{
		cfg:        cfg,
		clientGRPC: clientGRPC,
	}
}

// CreateUser.
func (s *UserServices) CreateUser1(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	out, err := s.clientGRPC.NewUser(ctx, in)
	return out, err
}

// CreateApp .
func (s *UserServices) CreateApp1(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	out, err := s.clientGRPC.NewApp(ctx, in)
	return out, err
}

// ReadUser.
func (s *UserServices) ReadUser1(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	out, err := s.clientGRPC.ReadUser(ctx, in)
	return out, err
}
