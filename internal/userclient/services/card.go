package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
)

type CardStorageProviver interface {
	CreateCard(ctx context.Context, record models.CardRecord) error
	UpdateCard(ctx context.Context, record models.CardRecord) error
	DeleteCard(ctx context.Context, record models.CardRecord) error
	SearchCard(ctx context.Context, searchInput string) ([]models.CardRecord, error)
	MarkCardSent(ctx context.Context, record models.CardRecord) error
}

// Services структура конструктора бизнес логики.
type CardServices struct {
	cfg        *config.ServiceConfig
	sl         CardStorageProviver
	clientGRPC ClientGRPCProvider
	Crypt
}

// New.
func NewCard(s CardStorageProviver, clientGRPC ClientGRPCProvider, cfg *config.ServiceConfig) *CardServices {
	return &CardServices{
		sl:         s,
		cfg:        cfg,
		clientGRPC: clientGRPC,
	}
}

// CardRec.
func (sr *CardServices) ProcessingCard(ctx context.Context, record models.CardRecord) error {
	var err error
	record.Brand, err = sr.EncryptAES(sr.cfg.Key, record.Brand)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	record.Number, err = sr.EncryptAES(sr.cfg.Key, record.Number)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	record.ValidDate, err = sr.EncryptAES(sr.cfg.Key, record.ValidDate)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	record.Code, err = sr.EncryptAES(sr.cfg.Key, record.Code)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	record.Holder, err = sr.EncryptAES(sr.cfg.Key, record.Holder)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	cardRecord := &pbpub.PublishCardRequest{
		ExchName:   sr.cfg.ExchName,
		RoutingKey: sr.cfg.RoutingKey + ".card",
		CardRecord: &pbpub.CardRecord{
			RecordID:  record.RecordID,
			ChngTime:  timestamppb.New(record.ChngTime),
			UID:       record.UID,
			AppID:     record.AppID,
			Brand:     record.Brand,
			Number:    record.Number,
			ValidDate: record.ValidDate,
			Code:      record.Code,
			Holder:    record.Holder,
			Metadata:  record.Metadata,
			Operation: int64(record.Operation),
		}}
	switch record.Operation {
	case models.Create:
		err := sr.sl.CreateCard(ctx, record)
		if err != nil {
			log.Print("create card record error: ", err)
			return err
		}
	case models.Update:
		err := sr.sl.UpdateCard(ctx, record)
		if err != nil {
			log.Print("update card record error: ", err)
			return err
		}
	case models.Delete:
		err := sr.sl.DeleteCard(ctx, record)
		if err != nil {
			log.Print("delete card record error: ", err)
			return err
		}
	}
	err = sr.clientGRPC.PublishCard(ctx, cardRecord)
	if err != nil {
		log.Print("publishing text record error: ", err)
		return err
	}
	err = sr.sl.MarkCardSent(ctx, record)
	if err != nil {
		log.Print("marking text record error: ", err)
		return err
	}
	return err
}

// SearchCard.
func (sr *CardServices) SearchCard(ctx context.Context, searchInput string) ([]models.CardRecord, error) {
	cardRecords, err := sr.sl.SearchCard(ctx, searchInput)
	if err != nil {
		log.Print("search binary record error: ", err)
	}
	for i := range cardRecords {
		cardRecords[i].Brand, err = sr.DecryptAES(sr.cfg.Key, cardRecords[i].Brand)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
		cardRecords[i].Number, err = sr.DecryptAES(sr.cfg.Key, cardRecords[i].Number)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
		cardRecords[i].ValidDate, err = sr.DecryptAES(sr.cfg.Key, cardRecords[i].ValidDate)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
		cardRecords[i].Code, err = sr.DecryptAES(sr.cfg.Key, cardRecords[i].Code)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
		cardRecords[i].Holder, err = sr.DecryptAES(sr.cfg.Key, cardRecords[i].Holder)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
	}
	return cardRecords, err
}
