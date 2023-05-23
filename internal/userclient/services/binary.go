package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
)

type BinaryStorageProviver interface {
	CreateBinary(ctx context.Context, record models.BinaryRecord) error
	UpdateBinary(ctx context.Context, record models.BinaryRecord) error
	DeleteBinary(ctx context.Context, record models.BinaryRecord) error
	SearchBinary(ctx context.Context, searchInput string) ([]models.BinaryRecord, error)
	MarkBinarySent(ctx context.Context, record models.BinaryRecord) error
}

// Services структура конструктора бизнес логики.
type BinaryServices struct {
	cfg        *config.ServiceConfig
	sl         BinaryStorageProviver
	clientGRPC ClientGRPCProvider
	Crypt
}

// New.
func NewBinary(s BinaryStorageProviver, clientGRPC ClientGRPCProvider, cfg *config.ServiceConfig) *BinaryServices {
	return &BinaryServices{
		sl:         s,
		cfg:        cfg,
		clientGRPC: clientGRPC,
	}
}

// BinaryRec.
func (sr *BinaryServices) ProcessingBinary(ctx context.Context, record models.BinaryRecord) error {
	var err error
	record.Binary, err = sr.EncryptAES(sr.cfg.Key, record.Binary)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	binaryRecord := &pbpub.PublishBinaryRequest{
		ExchName:   sr.cfg.ExchName,
		RoutingKey: sr.cfg.RoutingKey + ".binary",
		BinaryRecord: &pbpub.BinaryRecord{
			RecordID:  record.RecordID,
			ChngTime:  timestamppb.New(record.ChngTime),
			UID:       record.UID,
			AppID:     record.AppID,
			Binary:    record.Binary,
			Metadata:  record.Metadata,
			Operation: int64(record.Operation),
		}}
	switch record.Operation {
	case models.Create:
		err := sr.sl.CreateBinary(ctx, record)
		if err != nil {
			log.Print("create binary record error: ", err)
			return err
		}
	case models.Update:
		err := sr.sl.UpdateBinary(ctx, record)
		if err != nil {
			log.Print("update binary record error: ", err)
			return err
		}
	case models.Delete:
		err := sr.sl.DeleteBinary(ctx, record)
		if err != nil {
			log.Print("delete binary record error: ", err)
			return err
		}
	}
	err = sr.clientGRPC.PublishBinary(ctx, binaryRecord)
	if err != nil {
		log.Print("publishing binary record error: ", err)
		return err
	}
	err = sr.sl.MarkBinarySent(ctx, record)
	if err != nil {
		log.Print("marking binary record error: ", err)
		return err
	}
	return err
}

func (sr *BinaryServices) SearchBinary(ctx context.Context, searchInput string) ([]models.BinaryRecord, error) {
	binaryRecords, err := sr.sl.SearchBinary(ctx, searchInput)
	if err != nil {
		log.Print("search binary record error: ", err)
	}
	for i := range binaryRecords {
		binaryRecords[i].Binary, err = sr.DecryptAES(sr.cfg.Key, binaryRecords[i].Binary)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
	}
	return binaryRecords, err
}
