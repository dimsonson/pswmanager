package grpchandle

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/config"
	pb "github.com/dimsonson/pswmanager/internal/handlers/protobuf"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPChandlers struct {
	config.ServicesGRPC
	pb.UnimplementedUserServiceServer
}

func (s *GRPChandlers) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var out pb.CreateUserResponse
	var err error

	//s. . . .CreateUser(ctx, in.Login, in.Psw)
	//out.Value, out.Del, err = s.svsGet.Get(ctx, in.Key)

	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}

	// switch out{
	// case true:
	// 	// сообщаем что ссылка удалена
	// 	err = status.Errorf(codes.NotFound, `this link already deleted: %s`, in.Key)
	// 	out.Error = codes.NotFound.String()
	// case false:
	// 	// отправляем сокращенную сылку
	// 	out.Error = codes.OK.String()
	// 	return &out, err
	// }

	return &out, err
}
