package clientgrpc

import (
	"context"
	"time"

	"github.com/dimsonson/pswmanager/pkg/log"

	pbconsume "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protoconsume"
	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientGRPC struct {
	Cfg           *config.GRPC
	UserPBconn    pb.UserServicesClient
	ConsumePBconn pbconsume.ServerRMQhandlersClient
	PublishPBconn pbpub.ClientRMQhandlersClient
}

func NewClientGRPC(cfg *config.GRPC) (*ClientGRPC, error) {
	connGRPC, err := grpc.Dial(cfg.GatewayAddres, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
	}
	cfg.ClientConn = connGRPC
	log.Print("clientGRPC status: ", connGRPC.GetState().String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	connGRPC.WaitForStateChange(ctx, connectivity.Connecting)

	log.Print("clientGRPC status: ", connGRPC.GetState().String())

	users := pb.NewUserServicesClient(connGRPC)
	pub := pbpub.NewClientRMQhandlersClient(connGRPC)
	consume := pbconsume.NewServerRMQhandlersClient(connGRPC)

	return &ClientGRPC{
		Cfg:           cfg,
		UserPBconn:    users,
		ConsumePBconn: consume,
		PublishPBconn: pub,
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

func (cl *ClientGRPC) NewApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
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

// IsOnline возвращает true	если gRPC клиент онлайн и готов обслуживать запросы.
func (cl *ClientGRPC) IsOnline() bool {
	// проверка соединения gRPC
	return cl.Cfg.ClientConn.GetState() == connectivity.Ready
}

func (cl *ClientGRPC) PublishText(ctx context.Context, in *pbpub.PublishTextRequest) error {
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	_, err := cl.PublishPBconn.PublishText(ctx, in)
	if err != nil {
		log.Print("publish text error: ", err)
	}
	return err
}

func (cl *ClientGRPC) PublishLogins(ctx context.Context, in *pbpub.PublishLoginsRequest) error {
	log.Print("in", in)
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	_, err := cl.PublishPBconn.PublishLogins(ctx, in)
	if err != nil {
		log.Print("publish logins error: ", err)
	}
	return err
}

func (cl *ClientGRPC) PublishBinary(ctx context.Context, in *pbpub.PublishBinaryRequest) error {
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	_, err := cl.PublishPBconn.PublishBinary(ctx, in)
	if err != nil {
		log.Print("publish binary error: ", err)
	}
	return err
}

func (cl *ClientGRPC) PublishCard(ctx context.Context, in *pbpub.PublishCardRequest) error {
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	_, err := cl.PublishPBconn.PublishCard(ctx, in)
	if err != nil {
		log.Print("publish card error: ", err)
	}
	return err
}
