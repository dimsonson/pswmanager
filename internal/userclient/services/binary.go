package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type BinaryStorageProviver interface {
	CreateBinary(ctx context.Context, record models.BinaryRecord) error
	UpdateBinary(ctx context.Context, record models.BinaryRecord) error
	DeleteBinary(ctx context.Context, record models.BinaryRecord) error
	SearchBinary(ctx context.Context, searchInput string) ([]models.BinaryRecord, error)
}

// Services структура конструктора бизнес логики.
type BinaryServices struct {
	storage BinaryStorageProviver
}

// New.
func NewBinary(s BinaryStorageProviver) *BinaryServices {
	return &BinaryServices{
		s,
	}
}

// BinaryRec.
func (sr *BinaryServices) ProcessingBinary(ctx context.Context, record models.BinaryRecord) error {
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

func (sr *BinaryServices) SearchBinary(ctx context.Context, searchInput string) ([]models.BinaryRecord, error) {
	binaryRecords, err := sr.storage.SearchBinary(ctx, searchInput)
	if err != nil {
		log.Print("search binary record error: ", err)
	}
	return binaryRecords, err
}
