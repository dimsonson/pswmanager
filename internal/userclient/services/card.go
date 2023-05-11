package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type CardStorageProviver interface {
	CreateCard(ctx context.Context, record models.CardRecord) error
	UpdateCard(ctx context.Context, record models.CardRecord) error
	DeleteCard(ctx context.Context, record models.CardRecord) error
	SearchCard(ctx context.Context, searchInput string) ([]models.CardRecord, error)
}

// Services структура конструктора бизнес логики.
type CardServices struct {
	storage StorageProvider
}

// New.
func NewCard(s StorageProvider) *CardServices {
	return &CardServices{
		s,
	}
}

// CardRec.
func (sr *CardServices) ProcessingCard(ctx context.Context, record models.CardRecord) error {
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

func (sr *CardServices) SearchCard(ctx context.Context, searchInput string) ([]models.CardRecord, error) {
	cardRecords, err := sr.storage.SearchCard(ctx, searchInput)
	if err != nil {
		log.Print("search binary record error: ", err)
	}
	return cardRecords, err
}
