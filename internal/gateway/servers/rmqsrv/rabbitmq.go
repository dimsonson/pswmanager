package rmqsrv

import (
	"context"
	"fmt"
	"sync"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/MashinIvan/rabbitmq"
	"github.com/MashinIvan/rabbitmq/pkg/backoff"

	"github.com/dimsonson/pswmanager/internal/gateway/config"
	"github.com/dimsonson/pswmanager/internal/gateway/settings"
	"github.com/streadway/amqp"
)

// Client структура для хранения RabbitMQ client.
type Server struct {
	RabbitConn *rabbitmq.Connection
	RabbitSrv  *rabbitmq.Server
	Wg         *sync.WaitGroup
	Ctx        context.Context
	Stop       context.CancelFunc
	Cfg        config.RabbitmqSrv
}

func NewServer(ctx context.Context, stop context.CancelFunc, cfg config.RabbitmqSrv, wg *sync.WaitGroup) *Server {
	return &Server{
		Ctx:  ctx,
		Stop: stop,
		Wg:   wg,
		Cfg:  cfg,
	}
}

func (rmqs *Server) connFactory() (*amqp.Connection, error) {
	connUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		rmqs.Cfg.User,
		rmqs.Cfg.Psw,
		rmqs.Cfg.Host,
		rmqs.Cfg.Port,
	)
	return amqp.Dial(connUrl)
}

func (rmqs *Server) Init() {
	var err error
	rmqs.RabbitConn, err = rabbitmq.NewConnection(rmqs.connFactory, backoff.NewDefaultSigmoidBackoff())
	if err != nil {
		log.Print("rabbitmq connection starting error:", settings.ColorRed, err, settings.ColorReset)
	}
}

func (rmqs *Server) Start(ctx context.Context, router *rabbitmq.Router) {
	defer rmqs.Wg.Done()
	rmqs.RabbitSrv = rabbitmq.NewServer(rmqs.RabbitConn, router)
	err := rmqs.RabbitSrv.ListenAndServe(rmqs.Ctx)
	if err != nil {
		log.Print("rabbitmq server starting error:", settings.ColorRed, err, settings.ColorReset)
	}
	log.Print("rabbitmq server shutting down...")

}

func (rmqs *Server) Shutdown(ctx context.Context) {
	go func() {
		defer rmqs.Wg.Done()
		// получаем сигнал о завершении приложения
		//<-rmqs.Ctx.Done()
		<-ctx.Done()
		log.Print("rmqSrv got signal, attempting graceful shutdown")
		err := rmqs.RabbitSrv.Shutdown(ctx)
		if err != nil {
			log.Print("rabbitmq server shutdown error: ", settings.ColorRed, err, settings.ColorReset)
		}
		rmqs.RabbitConn.Close()
	}()
}
