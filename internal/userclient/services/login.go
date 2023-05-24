package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
)

type LoginStorageProviver interface {
	CreateLogin(ctx context.Context, record models.LoginRecord) error
	UpdateLogin(ctx context.Context, record models.LoginRecord) error
	DeleteLogin(ctx context.Context, record models.LoginRecord) error
	SearchLogin(ctx context.Context, searchInput string) ([]models.LoginRecord, error)
	MarkLoginSent(ctx context.Context, record models.LoginRecord) error
}

// Services структура конструктора бизнес логики.
type LoginServices struct {
	cfg        *config.ServiceConfig
	sl         LoginStorageProviver
	clientGRPC ClientGRPCProvider
	Crypt
}

// New.
func NewLogin(s LoginStorageProviver, clientGRPC ClientGRPCProvider, cfg *config.ServiceConfig) *LoginServices {
	return &LoginServices{
		sl:         s,
		cfg:        cfg,
		clientGRPC: clientGRPC,
	}
}

// LoginRec.
func (sr *LoginServices) ProcessingLogin(ctx context.Context, record models.LoginRecord) error {
	var err error
	record.Login, err = sr.EncryptAES(sr.cfg.Key, record.Login)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	record.Psw, err = sr.EncryptAES(sr.cfg.Key, record.Psw)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	loginRecord := &pbpub.PublishLoginsRequest{
		ExchName:   sr.cfg.ExchName,
		RoutingKey: sr.cfg.RoutingKey + ".login",
		LoginsRecord: &pbpub.LoginRecord{
			RecordID:  record.RecordID,
			ChngTime:  timestamppb.New(record.ChngTime),
			UID:       record.UID,
			AppID:     record.AppID,
			Login:     record.Login,
			Psw:       record.Psw,
			Metadata:  record.Metadata,
			Operation: int64(record.Operation),
		}}

	switch record.Operation {
	case models.Create:
		err := sr.sl.CreateLogin(ctx, record)
		if err != nil {
			log.Print("create login record error: ", err)
			return err
		}
	case models.Update:
		err := sr.sl.UpdateLogin(ctx, record)
		if err != nil {
			log.Print("update login record error: ", err)
			return err
		}
	case models.Delete:
		err := sr.sl.DeleteLogin(ctx, record)
		if err != nil {
			log.Print("delete login record error: ", err)
			return err
		}
	}
	err = sr.clientGRPC.PublishLogins(ctx, loginRecord)
	if err != nil {
		log.Print("publishing login record error: ", err)
		return err
	}
	err = sr.sl.MarkLoginSent(ctx, record)
	if err != nil {
		log.Print("marking login record error: ", err)
		return err
	}
	return err
}

func (sr *LoginServices) SearchLogin(ctx context.Context, searchInput string) ([]models.LoginRecord, error) {
	loginRecords, err := sr.sl.SearchLogin(ctx, searchInput)
	if err != nil {
		log.Print("rearch login record error: ", err)
	}
	for i := range loginRecords {
		loginRecords[i].Login, err = sr.DecryptAES(sr.cfg.Key, loginRecords[i].Login)
		if err != nil {
			log.Print("encrypt error: ", err)
			return nil, err
		}
		loginRecords[i].Psw, err = sr.DecryptAES(sr.cfg.Key, loginRecords[i].Psw)
		if err != nil {
			log.Print("encrypt error: ", err)
			return nil, err
		}
	}
	return loginRecords, err
}
