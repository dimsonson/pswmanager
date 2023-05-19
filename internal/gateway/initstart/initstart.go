package initstart

import (
	"context"
	"log"
	"sync"


	"github.com/dimsonson/pswmanager/internal/gateway/servers/rmqsrv"
	"github.com/dimsonson/pswmanager/internal/gateway/client/clientgrpc"
	clientrmq "github.com/dimsonson/pswmanager/internal/gateway/client/rmq"
	"github.com/dimsonson/pswmanager/internal/gateway/config"
	"github.com/dimsonson/pswmanager/internal/gateway/servers/grpc"
	"github.com/dimsonson/pswmanager/internal/gateway/services"

	grpchandlers "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers"
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

	clientRMQsvs := services.NewPub(clientRMQ)

	hndlRMQpub := grpchandlers.NewClientRMQhandlers(ctx, clientRMQsvs)

	clientGRPC, err := clientgrpc.NewClientGRPC(init.cfg.GRPC)
	if err != nil {
		log.Printf("new gRPC client error: %s", err)
		return
	}
	svsUsers := services.NewUserData(init.cfg, clientGRPC) //clientRMQ, clientGRPC)

	srvRMQ := rmqsrv.NewServer(ctx, stop, init.cfg.Rabbitmq, wg)

	hndlRMQconsume :=  grpchandlers.NewServerRMQhandlers(ctx, init.cfg.Rabbitmq, srvRMQ, wg)

	grpcSrv := grpc.NewServer(ctx, stop, init.cfg.GRPC, wg)
	grpcSrv.InitGRPCservice(svsUsers, hndlRMQpub, hndlRMQconsume)
	wg.Add(1)
	grpcSrv.GrpcGracefullShotdown()
	wg.Add(1)
	//go
	grpcSrv.StartGRPC()

	//handlers := rmq.New(servTextRec, servLoginRec, servBinaryRec, servCardRec)
	//rmqRouter := router.New(ctx, init.cfg.Rabbitmq, *handlers)

	
	
	
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
