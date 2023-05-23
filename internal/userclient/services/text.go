package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TextStorageProviver interface {
	CreateText(ctx context.Context, record models.TextRecord) error
	UpdateText(ctx context.Context, record models.TextRecord) error
	DeleteText(ctx context.Context, record models.TextRecord) error
	SearchText(ctx context.Context, searchInput string) ([]models.TextRecord, error)
	MarkTextSent(ctx context.Context, record models.TextRecord) error
}

// TextServices структура конструктора бизнес логики.
type TextServices struct {
	cfg        *config.ServiceConfig
	sl         TextStorageProviver
	clientGRPC ClientGRPCProvider
	Crypt
}

// NewText конструктор сервиса текстовых записей.
func NewText(s TextStorageProviver, clientGRPC ClientGRPCProvider, cfg *config.ServiceConfig) *TextServices {
	return &TextServices{
		sl:         s,
		cfg:        cfg,
		clientGRPC: clientGRPC,
	}
}

// ProcessingText метод обратботки данных в хранилище в зависимости от типа операции.
func (sr *TextServices) ProcessingText(ctx context.Context, record models.TextRecord, key string) error {
	var err error

	log.Print("record", record)

	log.Print("sr :", sr) // nil

	//log.Print("ProcessingText ucfg.Key", key)

	record.Text, err = sr.EncryptAES(sr.cfg.Key, record.Text)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}

	log.Print("ProcessingText record.Text", record.Text)
	switch record.Operation {
	case models.Create:
		err = sr.sl.CreateText(ctx, record)
		if err != nil {
			log.Print("create text record error: ", err)
			return err
		}
		err := sr.clientGRPC.CreateText(ctx, &pbpub.PublishTextRequest{
			ExchName:   sr.cfg.ExchName,
			RoutingKey: sr.cfg.RoutingKey + ".text",
			TextRecord: &pbpub.TextRecord{
				RecordID:  record.RecordID,
				ChngTime:  timestamppb.New(record.ChngTime),
				UID:       record.UID,
				AppID:     record.AppID,
				Text:      record.Text,
				Metadata:  record.Metadata,
				Operation: int64(record.Operation),
			},
		})
		if err != nil {
			log.Print("publishing text record error: ", err)
			return err
		}
	case models.Update:
		err := sr.sl.UpdateText(ctx, record)
		if err != nil {
			log.Print("update text record error: ", err)
		}
		return err
	case models.Delete:
		err := sr.sl.DeleteText(ctx, record)
		if err != nil {
			log.Print("delete text record error: ", err)
		}
		return err
	}
	return err
}

// SearchText метод поиск в хранилицще текстовых данных.
func (sr *TextServices) SearchText(ctx context.Context, searchInput string) ([]models.TextRecord, error) {
	textRecords, err := sr.sl.SearchText(ctx, searchInput)
	if err != nil {
		log.Print("rearch text record error: ", err)
	}
	for i := range textRecords {
		textRecords[i].Text, err = sr.DecryptAES(sr.cfg.Key, textRecords[i].Text)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
	}
	return textRecords, err
}
