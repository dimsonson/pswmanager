package router

import (
	"context"

	"github.com/MashinIvan/rabbitmq"
	"github.com/dimsonson/pswmanager/internal/masterserver/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/handlers/rmq"
)

func New(ctx context.Context, cfg config.RabbitmqSrv, handlers rmq.Handlers) *rabbitmq.Router {
	router := rabbitmq.NewRouter()
	// f:= func(ctx *rmqs.DeliveryContext) {
	// 	// process delivery
	// 	create := models.TextRecord{}

	// 	err := ctx.BindJSON(&create)
	// 	if err != nil {
	// 		log.Print(err)
	// 	}
	// //	err = hnd.servText.ProcessingText(ctx, create)
	// 	if err != nil {
	// 		return
	// 	}
	// 	err = ctx.Delivery.Ack(true)
	// 	if err != nil {
	// 		log.Print("received Ask error:", settings.ColorRed, err, settings.ColorReset)
	// 	}

	// }
	routerGroup := router.Group(
		rabbitmq.ExchangeParams{
			Name:       cfg.Exchange.Name,
			Kind:       cfg.Exchange.Kind,
			AutoDelete: false,
			Durable:    true},
		rabbitmq.QueueParams{
			Name:    cfg.Consumer.ConsumerName,
			Durable: true},
		rabbitmq.QualityOfService{},
		rabbitmq.ConsumerParams{},
		rabbitmq.WithRouterEngine(rabbitmq.NewTopicRouterEngine()), // .NewDirectRouterEngine()), // use direct to speed up routing
		rabbitmq.WithNumWorkers(cfg.RoutingWorkers),
	)

	routerGroup.
		Route(cfg.Controllers[0].RoutingKey, handlers.TextRec(ctx, cfg)).
		Route(cfg.Controllers[1].RoutingKey, handlers.LoginRec(ctx, cfg)).
		Route(cfg.Controllers[2].RoutingKey, handlers.BinaryRec(ctx, cfg)).
		Route(cfg.Controllers[3].RoutingKey, handlers.CardRec(ctx, cfg))

	return router
}
