package clientgrpc

import (
	"context"
	"log"

	"github.com/dimsonson/pswmanager/internal/gateway/config"
	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientGRPC struct {
	Cfg        config.GRPC
	UserPBconn pb.UserServicesClient
}

func NewClientGRPC(cfg config.GRPC) (*ClientGRPC, error) {
	connGRPC, err := grpc.Dial(cfg.MasterAddres, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
	}
	cfg.ClientConn = connGRPC
	log.Print(connGRPC.GetState().String())
	c := pb.NewUserServicesClient(connGRPC)
	return &ClientGRPC{
		Cfg:  cfg,
		UserPBconn: c,
	}, err
}

func (cl *ClientGRPC) NewUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	newUserCfg, err := cl.UserPBconn.CreateUser(ctx, in)
	if newUserCfg == nil {
		newUserCfg = &pb.CreateUserResponse{}
	}
	if err != nil {
		log.Print("create user error: ", err)
	}
	return newUserCfg, err
}

func (cl *ClientGRPC) NewApp(ctx context.Context,in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	newAppCfg, err := cl.UserPBconn.CreateApp(ctx, in)
	if newAppCfg == nil {
		newAppCfg = &pb.CreateAppResponse{}
	}
	if err != nil {
		log.Print("create app error: ", err)
	}
	return newAppCfg, err
}

func (cl *ClientGRPC) ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	// запрос всех записей пользователя
	newRead, err := cl.UserPBconn.ReadUser(ctx, in)
	if newRead == nil {
		newRead = &pb.ReadUserResponse{}
	}
	if err != nil {
		log.Print("read records error: ", err)
	}
	return newRead, err
}

func (cl *ClientGRPC) Close() {
	cl.Cfg.ClientConn.Close()
}
