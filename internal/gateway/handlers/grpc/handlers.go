package grpchandle

import (
	"context"

	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"github.com/dimsonson/pswmanager/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type UserServicesProvider interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	CreateApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) 
	ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) 
}

type UserHandlers struct {
	Ctx context.Context
	UserSvs UserServicesProvider
	pb.UnimplementedUserServicesServer
}

func NewUserHandlers(ctx context.Context, userSvs UserServicesProvider) *UserHandlers {
	return &UserHandlers{
		Ctx:     ctx,
		UserSvs: userSvs,
	}
}

func (h *UserHandlers) CreateUsers(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	out, err := h.UserSvs.CreateUser(ctx, in)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}

func (h *UserHandlers) CreateApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	out, err := h.UserSvs.CreateApp(ctx, in)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}

func (h *UserHandlers) ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	out, err := h.UserSvs.ReadUser(ctx, in)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}
