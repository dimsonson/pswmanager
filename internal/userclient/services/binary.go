package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
)

type BinaryStorageProviver interface {
	CreateBinary(ctx context.Context, record models.BinaryRecord) error
	UpdateBinary(ctx context.Context, record models.BinaryRecord) error
	DeleteBinary(ctx context.Context, record models.BinaryRecord) error
	SearchBinary(ctx context.Context, searchInput string) ([]models.BinaryRecord, error)
}

// Services структура конструктора бизнес логики.
type BinaryServices struct {
	cfg *config.ServiceConfig
	sl  BinaryStorageProviver
	c   CryptProvider
}

// New.
func NewBinary(s BinaryStorageProviver, cfg *config.ServiceConfig) *BinaryServices {
	return &BinaryServices{
		sl:  s,
		cfg: cfg,
		c:   &Crypt{},
	}
}

// BinaryRec.
func (sr *BinaryServices) ProcessingBinary(ctx context.Context, record models.BinaryRecord) error {
	var err error
	record.Binary, err = sr.c.EncryptAES(sr.cfg.Key, record.Binary)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	switch record.Operation {
	case models.Create:
		err := sr.sl.CreateBinary(ctx, record)
		if err != nil {
			log.Print("create binary record error: ", err)
		}
		return err
	case models.Update:
		err := sr.sl.UpdateBinary(ctx, record)
		if err != nil {
			log.Print("update binary record error: ", err)
		}
		return err
	case models.Delete:
		err := sr.sl.DeleteBinary(ctx, record)
		if err != nil {
			log.Print("delete binary record error: ", err)
		}
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
		binaryRecords[i].Binary, err = sr.c.DecryptAES(sr.cfg.Key, binaryRecords[i].Binary)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
	}
	return binaryRecords, err
}
