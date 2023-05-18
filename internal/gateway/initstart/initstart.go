package initstart

import (
	"context"
	"log"
	"sync"

	"github.com/dimsonson/pswmanager/internal/gateway/client/clientgrpc"
	clientrmq "github.com/dimsonson/pswmanager/internal/gateway/client/rmq"
	"github.com/dimsonson/pswmanager/internal/gateway/config"
	"github.com/dimsonson/pswmanager/internal/gateway/servers/grpc"
	"github.com/dimsonson/pswmanager/internal/gateway/services"
	//pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_hanglers"
)

type Init struct {
	cfg *config.ServiceConfig
}

func New(cfg *config.ServiceConfig) *Init {
	return &Init{
		cfg: cfg,
	}
}

// InitAndStart инициализация и старт clientRMQ, serverRMQ, serverGRPC, clientGRPC.
func (init *Init) InitAndStart(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup) {

	// //инициализируем и стартуем clientRMQ, serverRMQ, clientGRPC,  clientGRPC, serverGRPC :

	// clientRMQ нициализация конфигурацией и старт

	// инициализация сервисов

	// handlersRMQ инициализация сервисами
	// routerRMQ  инициализация конфигом, хендлерами
	// serverRMQ инициализация конфигом, роутером
	// serverRMQ graceful Shutdown
	// serverRMQ старт

	// serverGRPC инциализация конфигурацией
	// serverGRPC graceful Shutdown
	// serverGRPC старт

	clientRMQ, err := clientrmq.NewClientRMQ(init.cfg.Rabbitmq)
	if err != nil {
		log.Printf("new client error: %s", err)
		return
	}

	init.cfg.Rabbitmq.ClientRMQ.Conn = clientRMQ.Conn
	init.cfg.Rabbitmq.ClientRMQ.Ch = clientRMQ.Ch

	clientRMQsvs := services.NewTextPub(clientRMQ)

	handlersRMQpub := grpchanglers.NewClientRMQhandlers(ctx, clientRMQsvs)

	clientGRPC, err := clientgrpc.NewClientGRPC(init.cfg.GRPC)
	if err != nil {
		log.Printf("new gRPC client error: %s", err)
		return
	}
	//c := pb.NewUserServicesClient(clientGRPC.Conn)

	// newuser := &pb.CreateUserRequest{
	// 	Login: uuid.NewString(),
	// 	Psw:   "passw123test",
	// }

	// newUserCfg, err := clientGRPC.UserPBconn.CreateUser(ctx, newuser)
	// if err != nil {
	// 	log.Print("create user error: ", err)
	// }
	// log.Print(newUserCfg)

	// newapp := &pb.CreateAppRequest{
	// 	Uid: newUserCfg.UserID,
	// 	Psw: "passw123test",
	// }

	// newUserApp, err := clientGRPC.UserPBconn.CreateApp(ctx, newapp)
	// if err != nil {
	// 	log.Print("create user error: ", err)
	// }

	// log.Print(newUserApp)

	// msgText := models.TextRecord{
	// 	RecordID:  uuid.NewString(),
	// 	ChngTime:  time.Now(),
	// 	UID:       "aad84e70-7b9d-4dea-8a49-771ac186771b",
	// 	AppID:     "4834f853-8cf0-4814-a21a-82bd8a9cbb34",
	// 	Text:      "secured text sending",
	// 	Metadata:  "meta data description sample",
	// 	Operation: models.Create,
	// }
	// msgTextJSON, err := json.Marshal(msgText)
	// if err != nil {
	// 	log.Print("marshall error", err)
	// }

	// log.Print(msgTextJSON)

	// [123 34 82 101 99 111 114 100 73 68 34 58 34 56 101 49 97 50 97 50 100 45 54 49 101 50 45 52 100 101 57 45 98 57 55 99 45 49 100 101 57 101 54 50 56 54 55 102 54 34 44 34 67 104 110 103 84 105 109 101 34 58 34 50 48 50 51 45 48 53 45 49 56 84 49 49 58 53 55 58 52 51 46 52 52 50 55 53 53 43 48 51 58 48 48 34 44 34 85 73 68 34 58 34 97 97 100 56 52 101 55 48 45 55 98 57 100 45 52 100 101 97 45 56 97 52 57 45 55 55 49 97 99 49 56 54 55 55 49 98 34 44 34 65 112 112 73 68 34 58 34 52 56 51 52 102 56 53 51 45 56 99 102 48 45 52 56 49 52 45 97 50 49 97 45 56 50 98 100 56 97 57 99 98 98 51 52 34 44 34 84 101 120 116 34 58 34 115 101 99 117 114 101 100 32 116 101 120 116 32 115 101 110 100 105 110 103 34 44 34 77 101 116 97 100 97 116 97 34 58 34 109 101 116 97 32 100 97 116 97 32 100 101 115 99 114 105 112 116 105 111 110 32 115 97 109 112 108 101 34 44 34 79 112 101 114 97 116 105 111 110 34 58 49 125]

	// clientRMQ.PublishRecord(init.cfg.Rabbitmq.Exchange.Name,
	// 	"all.aad84e70-7b9d-4dea-8a49-771ac186771b.4834f853-8cf0-4814-a21a-82bd8a9cbb34.text",
	// 	msgTextJSON)

	// clientGRPC.NewUser(ctx, )

	svsUser := services.NewUserData(init.cfg, clientGRPC) //clientRMQ, clientGRPC)

	// servLoginRec := services.NewLogin(SQLstorage)
	// servTextRec := services.NewText(SQLstorage)
	// servCardRec := services.NewCard(SQLstorage)
	// servBinaryRec := services.NewBinary(SQLstorage)

	//cfgReadUsers := services.NewReadUser(SQLstorage)

	grpcSrv := grpc.NewServer(ctx, stop, init.cfg.GRPC, wg)
	grpcSrv.InitGRPCservice(svsUser, handlersRMQpub)
	wg.Add(1)
	grpcSrv.GrpcGracefullShotdown()
	wg.Add(1)
	//go
	grpcSrv.StartGRPC()

	//handlers := rmq.New(servTextRec, servLoginRec, servBinaryRec, servCardRec)
	//rmqRouter := router.New(ctx, init.cfg.Rabbitmq, *handlers)
	//rmqSrv := rabbitmq.NewServer(ctx, stop, init.cfg.Rabbitmq, wg)
	//rmqSrv.Init()
	// wg.Add(1)
	// rmqSrv.Shutdown()
	// wg.Add(1)
	// rmqSrv.Start(ctx, rmqRouter)

}

// ConnClose закрываем соединения при завершении работы.
func (init *Init) ConnClose(ctx context.Context) {
	if init.cfg.Rabbitmq.ClientRMQ.Conn != nil {
		init.cfg.Rabbitmq.ClientRMQ.Conn.Close()
	}
	if init.cfg.Rabbitmq.ClientRMQ.Ch != nil {
		init.cfg.Rabbitmq.ClientRMQ.Ch.Close()
	}
	if init.cfg.GRPC.ClientConn != nil {
		init.cfg.GRPC.ClientConn.Close()
	}

}
