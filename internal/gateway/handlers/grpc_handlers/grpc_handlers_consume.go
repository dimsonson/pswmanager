package grpchandlers

import (
	"context"

	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerRMQServicesProvider interface {
	PublishRecord(ctx context.Context, exchName string, routingKey string, record interface{}) error
}

type ServerRMQhandlers struct {
	Ctx          context.Context
	ClientRMQSvs ClientRMQServicesProvider
	pbpub.UnimplementedClientRMQhandlersServer
}

func NewServerRMQhandlers(ctx context.Context, clientRMQsvs ClientRMQServicesProvider) *ClientRMQhandlers {
	return &ClientRMQhandlers{
		Ctx:          ctx,
		ClientRMQSvs: clientRMQsvs,
	}
}

func (hs *ServerRMQhandlers) PublishText(ctx context.Context, in *pbpub.PublishTextRequest) (*pbpub.PublishTextResponse, error) {
	var out pbpub.PublishTextResponse
	var txtRecord models.TextRecord
	txtRecord.RecordID = in.TextRecord.RecordID
	txtRecord.ChngTime = in.TextRecord.ChngTime.AsTime()
	txtRecord.UID = in.TextRecord.UID
	txtRecord.AppID = in.TextRecord.AppID
	txtRecord.Text = in.TextRecord.Text
	txtRecord.Metadata = in.TextRecord.Metadata
	txtRecord.Operation = models.MsgType(in.TextRecord.Operation)

	err := hs.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, txtRecord)
	if err != nil {
		log.Printf("call PublishRecord error: %v", err)
		err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		if err != nil {
			log.Print("status.Errorf error")
		}
		out.Error = codes.Internal.String()
	}
	return &out, err
}

func (hs *ServerRMQhandlers) PublishLogins(ctx context.Context, in *pbpub.PublishLoginsRequest) (*pbpub.PublishLoginsResponse, error) {
	var out pbpub.PublishLoginsResponse
	var loginsRecord models.LoginRecord
	loginsRecord.RecordID = in.LoginsRecord.RecordID
	loginsRecord.ChngTime = in.LoginsRecord.ChngTime.AsTime()
	loginsRecord.UID = in.LoginsRecord.UID
	loginsRecord.AppID = in.LoginsRecord.AppID
	loginsRecord.Login = in.LoginsRecord.Login
	loginsRecord.Psw = in.LoginsRecord.Psw
	loginsRecord.Metadata = in.LoginsRecord.Metadata
	loginsRecord.Operation = models.MsgType(in.LoginsRecord.Operation)

	err := hs.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, loginsRecord)
	if err != nil {
		log.Printf("call PublishRecord error: %v", err)
		err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		if err != nil {
			log.Print("status.Errorf error")
		}
		out.Error = codes.Internal.String()
	}
	return &out, err
}

func (hs *ServerRMQhandlers) PublishBinary(ctx context.Context, in *pbpub.PublishBinaryRequest) (*pbpub.PublishBinaryResponse, error) {
	var out pbpub.PublishBinaryResponse
	var binaryRecord models.BinaryRecord
	binaryRecord.RecordID = in.BinaryRecord.RecordID
	binaryRecord.ChngTime = in.BinaryRecord.ChngTime.AsTime()
	binaryRecord.UID = in.BinaryRecord.UID
	binaryRecord.AppID = in.BinaryRecord.AppID
	binaryRecord.Binary = in.BinaryRecord.Binary
	binaryRecord.Metadata = in.BinaryRecord.Metadata
	binaryRecord.Operation = models.MsgType(in.BinaryRecord.Operation)

	err := hs.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, binaryRecord)
	if err != nil {
		log.Printf("call PublishRecord error: %v", err)
		err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		if err != nil {
			log.Print("status.Errorf error")
		}
		out.Error = codes.Internal.String()
	}
	return &out, err
}

func (hs *ServerRMQhandlers) PublishCard(ctx context.Context, in *pbpub.PublishCardRequest) (*pbpub.PublishCardResponse, error) {
	var out pbpub.PublishCardResponse
	var cardRecord models.CardRecord
	cardRecord.RecordID = in.CardRecord.RecordID
	cardRecord.ChngTime = in.CardRecord.ChngTime.AsTime()
	cardRecord.UID = in.CardRecord.UID
	cardRecord.AppID = in.CardRecord.AppID
	cardRecord.Brand = in.CardRecord.Brand
	cardRecord.Number = in.CardRecord.Number
	cardRecord.ValidDate = in.CardRecord.ValidDate
	cardRecord.Code = in.CardRecord.Code
	cardRecord.Holder = in.CardRecord.Holder
	cardRecord.Metadata = in.CardRecord.Metadata
	cardRecord.Operation = models.MsgType(in.CardRecord.Operation)

	err := hs.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, cardRecord)
	if err != nil {
		log.Printf("call PublishRecord error: %v", err)
		err = status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		if err != nil {
			log.Print("status.Errorf error")
		}
		out.Error = codes.Internal.String()
	}
	return &out, err
}
