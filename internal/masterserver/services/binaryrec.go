package services

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type BinaryStorageProviver interface {
	CreateBinary(ctx context.Context, record models.BinaryRec) error
	UpdateBinary(ctx context.Context, record models.BinaryRec) error
	DeleteBinary(ctx context.Context, record models.BinaryRec) error
}

// Services структура конструктора бизнес логики.
type BinaryServices struct {
	storage StorageProvider
}

// New.
func NewBinaryRec(s StorageProvider) *BinaryServices {
	return &BinaryServices{
		s,
	}
}

// BinaryRec.
func (sr *BinaryServices) BinaryRec(ctx context.Context, record models.BinaryRec) error {
	var err error
	switch record.Operation {
	case models.Create:
		err := sr.storage.CreateBinary(ctx, record)
		if err != nil {
			log.Print("create binary record error: ", err)
		}
		return err
	case models.Update:
		err := sr.storage.UpdateBinary(ctx, record)
		if err != nil {
			log.Print("update binary record error: ", err)
		}
		return err
	case models.Delete:
		err := sr.storage.DeleteBinary(ctx, record)
		if err != nil {
			log.Print("delete binary record error: ", err)
		}
		return err
	}
	return err
}
