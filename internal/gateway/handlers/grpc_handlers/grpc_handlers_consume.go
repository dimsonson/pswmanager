package grpchandlers

import (
	"context"
	"sync"

	"github.com/MashinIvan/rabbitmq"
	"github.com/dimsonson/pswmanager/internal/gateway/config"
	pbconsume "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protoconsume"
	"github.com/dimsonson/pswmanager/internal/gateway/servers/rmqsrv"
	"github.com/dimsonson/pswmanager/internal/gateway/settings"
	"github.com/dimsonson/pswmanager/pkg/log"
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

	//hc.Cfg.Consumer.ConsumerName = "all.02e0094f-667e-4958-9392-049b0aeea125.4e54529d-0eb4-42b9-b046-7d8b6e4ea84b"
	hc.Cfg.Consumer.ConsumerName = in.ConsumerQname
	hc.Cfg.Controllers[0].RoutingKey = in.RoutingKey
	log.Print("CONSUME")

	// out.Record = []byte{1}
	// err = stream.Send(&out)
	// if err != nil {
	// 	log.Print("sending to stream error")
	// 	return err
	// }

	router := rabbitmq.NewRouter()
	f := func(ctx *rabbitmq.DeliveryContext) {

		log.Print(" TO STREAM")
		
		out.Record = ctx.Delivery.Body
		// err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		// if err != nil {
		// 	log.Print("status.Errorf error")
		// }
		// out.Error = codes.Internal.String()

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
		Route(hc.Cfg.Controllers[0].RoutingKey, f).
		Route(hc.Cfg.Controllers[1].RoutingKey, f).
		Route(hc.Cfg.Controllers[2].RoutingKey, f).
		Route(hc.Cfg.Controllers[3].RoutingKey, f)

	hc.ServerRMQ.Init()

	hc.Wg.Add(1)
	hc.ServerRMQ.Shutdown()
	hc.Wg.Add(1)
	hc.ServerRMQ.Start(hc.Ctx, router)

	// var txtRecord models.TextRecord
	// txtRecord.RecordID = in.TextRecord.RecordID
	// txtRecord.ChngTime = in.TextRecord.ChngTime.AsTime()
	// txtRecord.UID = in.TextRecord.UID
	// txtRecord.AppID = in.TextRecord.AppID
	// txtRecord.Text = in.TextRecord.Text
	// txtRecord.Metadata = in.TextRecord.Metadata
	// txtRecord.Operation = models.MsgType(in.TextRecord.Operation)

	// err := hs.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, txtRecord)
	// if err != nil {
	// 	log.Printf("call PublishRecord error: %v", err)
	// 	err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
	// 	if err != nil {
	// 		log.Print("status.Errorf error")
	// 	}
	// 	out.Error = codes.Internal.String()
	// }
	return err
}

// func (hs *ServerRMQhandlers) PublishLogins(ctx context.Context, in *pbpub.PublishLoginsRequest) (*pbpub.PublishLoginsResponse, error) {
// 	var out pbpub.PublishLoginsResponse
// 	var loginsRecord models.LoginRecord
// 	loginsRecord.RecordID = in.LoginsRecord.RecordID
// 	loginsRecord.ChngTime = in.LoginsRecord.ChngTime.AsTime()
// 	loginsRecord.UID = in.LoginsRecord.UID
// 	loginsRecord.AppID = in.LoginsRecord.AppID
// 	loginsRecord.Login = in.LoginsRecord.Login
// 	loginsRecord.Psw = in.LoginsRecord.Psw
// 	loginsRecord.Metadata = in.LoginsRecord.Metadata
// 	loginsRecord.Operation = models.MsgType(in.LoginsRecord.Operation)

// 	err := hs.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, loginsRecord)
// 	if err != nil {
// 		log.Printf("call PublishRecord error: %v", err)
// 		err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
// 		if err != nil {
// 			log.Print("status.Errorf error")
// 		}
// 		out.Error = codes.Internal.String()
// 	}
// 	return &out, err
// }

// func (hs *ServerRMQhandlers) PublishBinary(ctx context.Context, in *pbpub.PublishBinaryRequest) (*pbpub.PublishBinaryResponse, error) {
// 	var out pbpub.PublishBinaryResponse
// 	var binaryRecord models.BinaryRecord
// 	binaryRecord.RecordID = in.BinaryRecord.RecordID
// 	binaryRecord.ChngTime = in.BinaryRecord.ChngTime.AsTime()
// 	binaryRecord.UID = in.BinaryRecord.UID
// 	binaryRecord.AppID = in.BinaryRecord.AppID
// 	binaryRecord.Binary = in.BinaryRecord.Binary
// 	binaryRecord.Metadata = in.BinaryRecord.Metadata
// 	binaryRecord.Operation = models.MsgType(in.BinaryRecord.Operation)

// 	err := hs.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, binaryRecord)
// 	if err != nil {
// 		log.Printf("call PublishRecord error: %v", err)
// 		err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
// 		if err != nil {
// 			log.Print("status.Errorf error")
// 		}
// 		out.Error = codes.Internal.String()
// 	}
// 	return &out, err
// }

// func (hs *ServerRMQhandlers) PublishCard(ctx context.Context, in *pbpub.PublishCardRequest) (*pbpub.PublishCardResponse, error) {
// 	var out pbpub.PublishCardResponse
// 	var cardRecord models.CardRecord
// 	cardRecord.RecordID = in.CardRecord.RecordID
// 	cardRecord.ChngTime = in.CardRecord.ChngTime.AsTime()
// 	cardRecord.UID = in.CardRecord.UID
// 	cardRecord.AppID = in.CardRecord.AppID
// 	cardRecord.Brand = in.CardRecord.Brand
// 	cardRecord.Number = in.CardRecord.Number
// 	cardRecord.ValidDate = in.CardRecord.ValidDate
// 	cardRecord.Code = in.CardRecord.Code
// 	cardRecord.Holder = in.CardRecord.Holder
// 	cardRecord.Metadata = in.CardRecord.Metadata
// 	cardRecord.Operation = models.MsgType(in.CardRecord.Operation)

// 	err := hs.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, cardRecord)
// 	if err != nil {
// 		log.Printf("call PublishRecord error: %v", err)
// 		err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
// 		if err != nil {
// 			log.Print("status.Errorf error")
// 		}
// 		out.Error = codes.Internal.String()
// 	}
// 	return &out, err
// }
