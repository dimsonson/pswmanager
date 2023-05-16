package grpchandle

import (
	"context"

	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReadUserServicesProvider interface {
	ReadUser(ctx context.Context, uid string) (models.SetOfRecords, error)
}

type UserServicesProvider interface {
	CreateUser(ctx context.Context, login string, psw string) (*pb.CreateUserResponse, error)
	CreateApp(ctx context.Context, uid string, psw string) (*pb.CreateAppResponse, error)
	ReadUser(ctx context.Context, uid string) (*pb.ReadUserResponse, error)
}

type UserHandlers struct {
	Ctx context.Context
	ReadUserServicesProvider
	UserServicesProvider
	pb.UnimplementedUserServicesServer
}

func NewUserHandlers(ctx context.Context) *UserHandlers {
	return &UserHandlers{
		Ctx: ctx,
	}
}

func (h *UserHandlers) CreateUsers(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	out, err := h.UserServicesProvider.CreateUser(ctx, in.Login, in.Psw)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}

func (h *UserHandlers) CreateApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	out, err := h.UserServicesProvider.CreateApp(ctx, in.Uid, in.Psw)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}

func (h *UserHandlers) ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	out, err := h.UserServicesProvider.ReadUser(ctx, in.Uid)
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return out, err
}
