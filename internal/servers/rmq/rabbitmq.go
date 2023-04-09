package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/MashinIvan/rabbitmq"
	"github.com/MashinIvan/rabbitmq/pkg/backoff"

	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/dimsonson/pswmanager/internal/settings"
	"github.com/streadway/amqp"
)

// Client структура для хранения RabbitMQ client.
type Server struct {
	RabbitConn *rabbitmq.Connection
	RabbitSrv  *rabbitmq.Server
	Wg         *sync.WaitGroup
	Ctx        context.Context
	Stop       context.CancelFunc
	Cfg        models.RabbitmqSrv
}

func NewServer(ctx context.Context, stop context.CancelFunc, cfg models.RabbitmqSrv, wg *sync.WaitGroup) *Server {
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

	rmqs.RabbitSrv = rabbitmq.NewServer(rmqs.RabbitConn, router)

	// log.Printf("rmqs.Cfg.Consumer: %v\n", rmqs.Cfg.Consumer)
	// log.Printf("rmqs.Cfg.Controllers: %v\n", rmqs.Cfg.Controllers)
	// log.Printf("rmqs.Cfg.Exchange: %v\n", rmqs.Cfg.Exchange)
	// log.Printf("rmqs.Cfg.Queue: %v\n", rmqs.Cfg.Queue)

	err := rmqs.RabbitSrv.ListenAndServe(rmqs.Ctx)
	if err != nil {
		log.Print("rabbitmq server starting error:", settings.ColorRed, err, settings.ColorReset)
	}
	log.Print("rabbitmq server shutting down...")
	rmqs.Wg.Done()

}

func (rmqs *Server) Shutdown() {
	go func() {
		// получаем сигнал о завершении приложения
		<-rmqs.Ctx.Done()
		log.Print("rmqSrv got signal, attempting graceful shutdown")
		err := rmqs.RabbitSrv.Shutdown(rmqs.Ctx)
		if err != nil {
			log.Print("rabbitmq server shutdown error: ", settings.ColorRed, err, settings.ColorReset)
		}
		rmqs.Wg.Done()
	}()
}
