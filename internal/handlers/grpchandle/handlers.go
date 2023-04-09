package grpchandle

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/config"
	pb "github.com/dimsonson/pswmanager/internal/handlers/protobuf"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPChandlers struct {
	config.ServicesGRPC
	pb.UnimplementedUserServiceServer
}

func (s *GRPChandlers) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var out pb.CreateUserResponse
	var err error
	usercfg, err := s.User.CreateUser(ctx, in.Login, in.Psw)
	out.UserID = usercfg.UserID
	out.RmqHost = usercfg.RmqHost
	out.RmqPort = usercfg.RmqPort
	out.RmqUID = usercfg.RmqUID
	out.ExchangeName = usercfg.ExchangeName
	out.ExchangeKind = usercfg.ExchangeKind

	for _, v := range usercfg.Apps {
		strout := &pb.App{
			AppID:        v.AppID,
			RoutingKey:   v.RoutingKey,
			ConsumeQueue: v.ConsumeQueue,
			ConsumerName: v.ConsumerName,
		}
		out.Apps = append(out.Apps, strout)
	}

	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}

	return &out, err
}


func (s *GRPChandlers) CreateApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	var out pb.CreateAppResponse
	var err error
	appid, usercfg, err := s.User.CreateApp(ctx, in.Uid, in.Psw)
	out.Appid = appid
	out.UserID = usercfg.UserID
	out.RmqHost = usercfg.RmqHost
	out.RmqPort = usercfg.RmqPort
	out.RmqUID = usercfg.RmqUID
	out.ExchangeName = usercfg.ExchangeName
	out.ExchangeKind = usercfg.ExchangeKind

	for _, v := range usercfg.Apps {
		strout := &pb.App{
			AppID:        v.AppID,
			RoutingKey:   v.RoutingKey,
			ConsumeQueue: v.ConsumeQueue,
			ConsumerName: v.ConsumerName,
		}
		out.Apps = append(out.Apps, strout)
	}

	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}

	return &out, err
}

func (s *GRPChandlers) ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	var out pb.ReadUserResponse
	var err error
	setrecords, err := s.ReadUsers.ReadUser(ctx, in.Uid)

	for _, v := range setrecords.SetTextRec  {
		strout := &pb.TextRec{
			RecordID: v.RecordID,
			//ChngTime: v.ChngTime,
			UID: v.UID,
			AppID: v.AppID,
			Text: v.Text,
			Metadata: v.Metadata,
		}
		out.SetTextRec = append(out.SetTextRec, strout)
	}
	
	for _, v := range setrecords.SetBinaryRec  {
		strout := &pb.BinaryRec{
			RecordID: v.RecordID,
			//ChngTime: v.ChngTime,
			UID: v.UID,
			AppID: v.AppID,
			Binary: v.Binary,
			Metadata: v.Metadata,
		}
		out.SetBinaryRec = append(out.SetBinaryRec, strout)
	}

	for _, v := range setrecords.SetLoginRec  {
		strout := &pb.LoginRec{
			RecordID: v.RecordID,
			//ChngTime: v.ChngTime,
			UID: v.UID,
			AppID: v.AppID,
			Login: v.Login,
			Psw: v.Psw,
			Metadata: v.Metadata,
		}
		out.SetLoginRec = append(out.SetLoginRec, strout)
	}

	for _, v := range setrecords.SetCardRec  {
		strout := &pb.CardRec{
			RecordID: v.RecordID,
			//ChngTime: v.ChngTime,
			UID: v.UID,
			AppID: v.AppID,
			Brand: int64(v.Brand),
			Number: v.Number,
			ValidDate: v.ValidDate,
			Code: int64(v.Code),
			Holder: v.Holder,
			Metadata: v.Metadata,
		}
		out.SetCardRec = append(out.SetCardRec, strout)
	}
	
	if err != nil {
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	}

	return &out, err
}