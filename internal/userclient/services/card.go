package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
)

type CardStorageProviver interface {
	CreateCard(ctx context.Context, record models.CardRecord) error
	UpdateCard(ctx context.Context, record models.CardRecord) error
	DeleteCard(ctx context.Context, record models.CardRecord) error
	SearchCard(ctx context.Context, searchInput string) ([]models.CardRecord, error)
}

// Services структура конструктора бизнес логики.
type CardServices struct {
	cfg *config.ServiceConfig
	sl  CardStorageProviver
	c   CryptProvider
}

// New.
func NewCard(s CardStorageProviver, cfg *config.ServiceConfig) *CardServices {
	return &CardServices{
		sl:  s,
		cfg: cfg,
		c:   &Crypt{},
	}
}

// CardRec.
func (sr *CardServices) ProcessingCard(ctx context.Context, record models.CardRecord) error {
	var err error

	record.Brand, err = sr.c.EncryptAES(sr.cfg.Key, record.Brand)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}

	record.Number, err = sr.c.EncryptAES(sr.cfg.Key, record.Number)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}

	record.ValidDate, err = sr.c.EncryptAES(sr.cfg.Key, record.ValidDate)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}

	record.Code, err = sr.c.EncryptAES(sr.cfg.Key, record.Code)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}

	record.Holder, err = sr.c.EncryptAES(sr.cfg.Key, record.Holder)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}

	switch record.Operation {
	case models.Create:
		err := sr.sl.CreateCard(ctx, record)
		if err != nil {
			log.Print("create card record error: ", err)
		}
		return err
	case models.Update:
		err := sr.sl.UpdateCard(ctx, record)
		if err != nil {
			log.Print("update card record error: ", err)
		}
		return err
	case models.Delete:
		err := sr.sl.DeleteCard(ctx, record)
		if err != nil {
			log.Print("delete card record error: ", err)
		}
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
		cardRecords[i].Brand, err = sr.c.DecryptAES(sr.cfg.Key, cardRecords[i].Brand)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
		cardRecords[i].Number, err = sr.c.DecryptAES(sr.cfg.Key, cardRecords[i].Number)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
		cardRecords[i].ValidDate, err = sr.c.DecryptAES(sr.cfg.Key, cardRecords[i].ValidDate)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
		cardRecords[i].Code, err = sr.c.DecryptAES(sr.cfg.Key, cardRecords[i].Code)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
		cardRecords[i].Holder, err = sr.c.DecryptAES(sr.cfg.Key, cardRecords[i].Holder)
		if err != nil {
			log.Print("decrypt error: ", err)
			return nil, err
		}
	}
	return cardRecords, err
}
