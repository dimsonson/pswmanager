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
	Conn       *grpc.ClientConn
	UserPBconn pb.UserServicesClient
}

func NewClientGRPC(cfg config.GRPC) (*ClientGRPC, error) {
	connGRPC, err := grpc.Dial(cfg.MasterAddres, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
	}
	log.Print(connGRPC.GetState().String())
	c := pb.NewUserServicesClient(connGRPC)
	return &ClientGRPC{
		Cfg:  cfg,
		Conn: connGRPC,
		UserPBconn: c,
	}, err
}

// func NewUserApp(ctx context.Context, c pb.UserServicesClient) (*pb.CreateUserResponse, *pb.CreateAppResponse, error) {
// 	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
// 	newuser := &pb.CreateUserRequest{
// 		Login: uuid.NewString(),
// 		Psw:   "passw123test",
// 	}
// 	newUserCfg, err := c.CreateUser(ctx, newuser)
// 	if err != nil {
// 		log.Print("create user error: ", err)
// 	}
// 	newapp := &pb.CreateAppRequest{
// 		Uid: newUserCfg.UserID,
// 		Psw: "passw123test",
// 	}
// 	newAppCfg, err := c.CreateApp(ctx, newapp)
// 	if err != nil {
// 		log.Print("create app error: ", err)
// 	}
// 	return newUserCfg, newAppCfg, err
// }

// func (cl *ClientGRPC) NewUser(ctx context.Context, login string, psw string) (*pb.CreateUserResponse, error) {
// 	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
// 	newuser := &pb.CreateUserRequest{
// 		Login: login,
// 		Psw:   psw,
// 	}
// 	newUserCfg, err := cl.UserPBconn.CreateUser(ctx, newuser)
// 	if err != nil {
// 		log.Print("create user error: ", err)
// 	}
// 	return newUserCfg, err
// }

func (cl *ClientGRPC) NewUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	
	newUserCfg, err := cl.UserPBconn.CreateUser(ctx, in)
	if err != nil {
		log.Print("create user error: ", err)
	}
	return newUserCfg, err
}

func (cl *ClientGRPC) NewApp(ctx context.Context, uid string, psw string) (*pb.CreateAppResponse, error) {
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	newapp := &pb.CreateAppRequest{
		Uid: uid,
		Psw: psw,
	}
	newAppCfg, err := cl.UserPBconn.CreateApp(ctx, newapp)
	if err != nil {
		log.Print("create app error: ", err)
	}
	return newAppCfg, err
}

func (cl *ClientGRPC) ReadUser(ctx context.Context, newAppCfg *pb.CreateAppResponse) (*pb.ReadUserResponse, error) {
	// переменная запроса всех записей пользователя
	newread := &pb.ReadUserRequest{
		Uid: newAppCfg.UserID,
	}
	// запрос всех записей пользователя
	newRead, err := cl.UserPBconn.ReadUser(ctx, newread)
	if err != nil {
		log.Print("read records error: ", err)
	}
	return newRead, err
}
