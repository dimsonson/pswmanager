package services

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type ReadUserStorageProviver interface {
	ReadUserRecords(ctx context.Context, userID string) (*models.SetRecords, error)
}

// Services структура конструктора бизнес логики.
type ReadUserServices struct {
	storage StorageProvider
}

// New.
func NewReadUser(s StorageProvider) *ReadUserServices {
	return &ReadUserServices{
		s,
	}
}

func (sr *ReadUserServices) ReadUser(ctx context.Context, uid string) (models.SetRecords, error) {
	setRecords, err := sr.storage.ReadUserRecords(ctx, uid)
	return *setRecords, err
}
