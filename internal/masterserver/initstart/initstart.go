package initstart

import (
	"context"
	"log"
	"sync"

	"github.com/dimsonson/pswmanager/internal/masterserver/clientrmq"
	"github.com/dimsonson/pswmanager/internal/masterserver/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/handlers/rmq"
	"github.com/dimsonson/pswmanager/internal/masterserver/router"
	"github.com/dimsonson/pswmanager/internal/masterserver/servers/grpc"
	rabbitmq "github.com/dimsonson/pswmanager/internal/masterserver/servers/rmq"
	"github.com/dimsonson/pswmanager/internal/masterserver/services"
	"github.com/dimsonson/pswmanager/internal/masterserver/storage/nosql"
	"github.com/dimsonson/pswmanager/internal/masterserver/storage/sql"
)

type Init struct {
	cfg *config.ServiceConfig
}

func New(cfg *config.ServiceConfig) *Init {
	return &Init{
		cfg: cfg,
	}
}

func (init *Init) InitAndStart(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup) {

	noSQLstorage := nosql.New(init.cfg.Redis)

	clientRMQ, err := clientrmq.NewClientRMQ(init.cfg.Rabbitmq)
	if err != nil {
		log.Printf("new client error: %s", err)
		return
	}
	init.cfg.Rabbitmq.ClientRMQ.Conn = clientRMQ.Conn
	init.cfg.Rabbitmq.ClientRMQ.Ch = clientRMQ.Ch
	cfgUser := services.NewUserData(noSQLstorage, clientRMQ, init.cfg.Rabbitmq)

	SQLstorage := sql.New(init.cfg.Postgree.Dsn)
	init.cfg.Postgree.Conn = SQLstorage.PostgreConn

	servLoginRec := services.NewLogin(SQLstorage)
	servTextRec := services.NewText(SQLstorage)
	servCardRec := services.NewCard(SQLstorage)
	servBinaryRec := services.NewBinary(SQLstorage)

	cfgReadUsers := services.NewReadUser(SQLstorage)

	handlers := rmq.New(servTextRec, servLoginRec, servBinaryRec, servCardRec)

	grpcSrv := grpc.NewServer(ctx, stop, init.cfg.GRPC, wg)
	grpcSrv.InitGRPCservice(cfgReadUsers, cfgUser)
	wg.Add(1)
	grpcSrv.GrpcGracefullShotdown()
	wg.Add(1)
	go grpcSrv.StartGRPC()

	rmqRouter := router.New(ctx, init.cfg.Rabbitmq, *handlers)
	rmqSrv := rabbitmq.NewServer(ctx, stop, init.cfg.Rabbitmq, wg)
	rmqSrv.Init()
	wg.Add(1)
	rmqSrv.Shutdown()
	wg.Add(1)
	rmqSrv.Start(ctx, rmqRouter)

}
