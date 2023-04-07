package services

import (
	"context"

	"github.com/rs/zerolog/log"
)

type ReadUserRecStorageProviver interface {
	ReadUserRecords(ctx context.Context, userID string) error
}

// Services структура конструктора бизнес логики.
type ReadUserRecServices struct {
	storage StorageProvider
}

// New.
func NewReadUserRec(s StorageProvider) *BinaryServices {
	return &BinaryServices{
		s,
	}
}

// BinaryRec.
func (sr *BinaryServices) ReadUserRec(ctx context.Context, userID string) error {
	err := sr.storage.ReadUserRecords(ctx, userID)
	if err != nil {
		log.Print("read user records error: ", err)
	}
	return err
}
