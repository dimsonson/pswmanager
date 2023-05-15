package initstart

import (
	"context"
	"sync"

	"github.com/dimsonson/pswmanager/internal/masterserver/config"
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


	// инициализируем и стартуем clientRMQ, serverRMQ, clientGRPC,  clientGRPC, serverGRPC :
    
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







	// noSQLstorage := nosql.New(init.cfg.Redis)

	// clientRMQ, err := clientrmq.NewClientRMQ(init.cfg.Rabbitmq)
	// if err != nil {
	// 	log.Printf("new client error: %s", err)
	// 	return
	// }

	// init.cfg.Rabbitmq.ClientRMQ.Conn = clientRMQ.Conn
	// init.cfg.Rabbitmq.ClientRMQ.Ch = clientRMQ.Ch

	// cfgUser := services.NewUserData(noSQLstorage, clientRMQ, init.cfg.Rabbitmq)

	// SQLstorage := sql.New(init.cfg.Postgree.Dsn)
	// init.cfg.Postgree.Conn = SQLstorage.PostgreConn

	// servLoginRec := services.NewLogin(SQLstorage)
	// servTextRec := services.NewText(SQLstorage)
	// servCardRec := services.NewCard(SQLstorage)
	// servBinaryRec := services.NewBinary(SQLstorage)

	// cfgReadUsers := services.NewReadUser(SQLstorage)
	
	// grpcSrv := grpc.NewServer(ctx, stop, init.cfg.GRPC, wg)
	// grpcSrv.InitGRPCservice(cfgReadUsers, cfgUser)
	// wg.Add(1)
	// grpcSrv.GrpcGracefullShotdown()
	// wg.Add(1)
	// go grpcSrv.StartGRPC()
	

	// handlers := rmq.New(servTextRec, servLoginRec, servBinaryRec, servCardRec)
	// rmqRouter := router.New(ctx, init.cfg.Rabbitmq, *handlers)
	// rmqSrv := rabbitmq.NewServer(ctx, stop, init.cfg.Rabbitmq, wg)
	// rmqSrv.Init()
	// wg.Add(1)
	// rmqSrv.Shutdown()
	// wg.Add(1)
	// rmqSrv.Start(ctx, rmqRouter)

}

// ConnClose закрываем соединения при завершении работы.
func (init *Init) ConnClose(ctx context.Context) {
	// if cfg.Postgree.Conn != nil {
	// 	cfg.Postgree.Conn.Close()
	// }
	// if cfg.Rabbitmq.ClientRMQ.Conn != nil {
	// 	cfg.Rabbitmq.ClientRMQ.Conn.Close()
	// }
	// if cfg.Rabbitmq.ClientRMQ.Ch != nil {
	// 	cfg.Rabbitmq.ClientRMQ.Ch.Close()
	// }
}