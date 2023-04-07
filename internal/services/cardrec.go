package services

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/models"
)

type CardStorageProviver interface {
	CreateCard(ctx context.Context, record models.CardRec) error
	UpdateCard(ctx context.Context, record models.CardRec) error
	DeleteCard(ctx context.Context, record models.CardRec) error
}

// Services структура конструктора бизнес логики.
type CardServices struct {
	storage StorageProvider
}

// New.
func NewCardRec(s StorageProvider) *CardServices {
	return &CardServices{
		s,
	}
}

// CardRec.
func (sr *CardServices) CardRec(ctx context.Context, record models.CardRec) error {
	var err error
	switch record.Operation {
	case models.Create:
		err := sr.storage.CreateCard(ctx, record)
		if err != nil {
			log.Print("create card record error: ", err)
		}
		return err
	case models.Update:
		err := sr.storage.UpdateCard(ctx, record)
		if err != nil {
			log.Print("update card record error: ", err)
		}
		return err
	case models.Delete:
		err := sr.storage.DeleteCard(ctx, record)
		if err != nil {
			log.Print("delete card record error: ", err)
		}
		return err

	}
	return err
}