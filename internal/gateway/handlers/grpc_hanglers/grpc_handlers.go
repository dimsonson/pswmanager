package grpchanglers

import (
	"context"

	//pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_hanglers/proto"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientRMQServicesProvider interface {
	PublishRecord(ctx context.Context, exchName string, routingKey string, record interface{}) error
}

type ClientRMQhandlers struct {
	Ctx          context.Context
	ClientRMQSvs ClientRMQServicesProvider
	pbpub.UnimplementedClientRMQhandlersServer
}

func NewClientRMQhandlers(ctx context.Context, clientRMQsvs ClientRMQServicesProvider) *ClientRMQhandlers {
	return &ClientRMQhandlers{
		Ctx:          ctx,
		ClientRMQSvs: clientRMQsvs,
	}
}

func (h *ClientRMQhandlers) PublishText(ctx context.Context, in *pbpub.PublishTextRequest) (*pbpub.PublishTextResponse, error) {
	var out pbpub.PublishTextResponse
	var txtRecord models.TextRecord
	txtRecord.RecordID = in.TextRecord.RecordID
	// txtRecord.ChngTime = time.Time(in.TextRecord.ChngTime)
	txtRecord.UID = in.TextRecord.UID
	txtRecord.AppID = in.TextRecord.AppID
	txtRecord.Text = in.TextRecord.Text
	txtRecord.Metadata = in.TextRecord.Metadata
	txtRecord.Operation = models.MsgType(in.TextRecord.Operation)

	err := h.ClientRMQSvs.PublishRecord(ctx, in.ExchName, in.RoutingKey, txtRecord)
	if err != nil {
		log.Printf("call PublishRecord error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}
	return &out, err
}

// func (h *ClientRMQhandlers) CreateApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
// 	out, err := h.UserSvs.CreateApp(ctx, in)
// 	if err != nil {
// 		log.Printf("call Put error: %v", err)
// 		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
// 		out.Error = codes.Internal.String()
// 	}
// 	return out, err
// }

// func (h *ClientRMQhandlers) ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
// 	out, err := h.UserSvs.ReadUser(ctx, in)
// 	if err != nil {
// 		log.Printf("call Put error: %v", err)
// 		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
// 		out.Error = codes.Internal.String()
// 	}
// 	return out, err
// }
