package grpchandlers

import (
	"context"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/MashinIvan/rabbitmq"
	"github.com/dimsonson/pswmanager/internal/gateway/config"
	pbconsume "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protoconsume"
	"github.com/dimsonson/pswmanager/internal/gateway/servers/rmqsrv"
	"github.com/dimsonson/pswmanager/internal/gateway/settings"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerRMQServicesProvider interface {
	PublishRecord(ctx context.Context, exchName string, routingKey string, record interface{}) error
}

type ServerRMQhandlers struct {
	Ctx       context.Context
	Wg        *sync.WaitGroup
	Cfg       config.RabbitmqSrv
	ServerRMQ *rmqsrv.Server
	pbconsume.UnimplementedServerRMQhandlersServer
}

func NewServerRMQhandlers(ctx context.Context, cfg config.RabbitmqSrv, serverRMQ *rmqsrv.Server, wg *sync.WaitGroup) *ServerRMQhandlers {
	return &ServerRMQhandlers{
		Ctx:       ctx,
		Wg:        wg,
		Cfg:       cfg,
		ServerRMQ: serverRMQ,
	}
}

func (hc *ServerRMQhandlers) Consume(in *pbconsume.ConsumeRequest, stream pbconsume.ServerRMQhandlers_ConsumeServer) error {
	var out pbconsume.ConsumeResponse
	var err error
	hc.Cfg.Consumer.ConsumerName = in.ConsumerQname
	hc.Cfg.Controllers[0].RoutingKey = in.RoutingKey
	router := rabbitmq.NewRouter()
	f := func(ctx *rabbitmq.DeliveryContext) {
		// создание поля тип записи для ответа в зависимости от последней части ключа роутинга
		routing := strings.Split(ctx.Delivery.RoutingKey, ".")
		switch routing[len(routing)-1] {
		case "text":
			out.RecordType = int64(models.TextType)
		case "login":
			out.RecordType = int64(models.LoginsType)
		case "bunary":
			out.RecordType = int64(models.BinaryType)
		case "card":
			out.RecordType = int64(models.CardType)
		}
		out.Record = ctx.Delivery.Body

		err = ctx.Err()
		if err != nil {
			log.Printf("stream error: %v", err)
			err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
			if err != nil {
				log.Print("status.Errorf error")
			}
			out.Error = codes.Internal.String()
		}

		err = stream.Send(&out)
		if err != nil {
			log.Print("sending to stream error")
			return
		}

		err = ctx.Delivery.Ack(true)
		if err != nil {
			log.Print("received Ask error:", settings.ColorRed, err, settings.ColorReset)
		}

	}
	routerGroup := router.Group(
		rabbitmq.ExchangeParams{
			Name:       hc.Cfg.Exchange.Name,
			Kind:       hc.Cfg.Exchange.Kind,
			AutoDelete: false,
			Durable:    true},
		rabbitmq.QueueParams{
			Name:    hc.Cfg.Consumer.ConsumerName,
			Durable: true},
		rabbitmq.QualityOfService{},
		rabbitmq.ConsumerParams{},
		rabbitmq.WithRouterEngine(rabbitmq.NewTopicRouterEngine()), // .NewDirectRouterEngine()), // use direct to speed up routing
		rabbitmq.WithNumWorkers(hc.Cfg.RoutingWorkers),
	)
	routerGroup.
		Route(hc.Cfg.Controllers[0].RoutingKey, f)

	hc.ServerRMQ.Init()
	ctxStream, _ := signal.NotifyContext(stream.Context(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	hc.Wg.Add(1)
	hc.ServerRMQ.Shutdown(ctxStream)
	hc.Wg.Add(1)
	hc.ServerRMQ.Start(ctxStream, router)
	return err
}

