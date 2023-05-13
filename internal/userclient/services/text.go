package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
)

type TextStorageProviver interface {
	CreateText(ctx context.Context, record models.TextRecord) error
	UpdateText(ctx context.Context, record models.TextRecord) error
	DeleteText(ctx context.Context, record models.TextRecord) error
	SearchText(ctx context.Context, searchInput string) ([]models.TextRecord, error)
}

// Services структура конструктора бизнес логики.
type TextServices struct {
	cfg *config.ServiceConfig
	sl  StorageProvider
	c   CryptProvider
}

// New.
func NewText(s StorageProvider, cfg *config.ServiceConfig) *TextServices {
	return &TextServices{
		sl:  s,
		cfg: cfg,
		c:   &Crypt{},
	}
}

// TextRec.
func (sr *TextServices) ProcessingText(ctx context.Context, record models.TextRecord) error {
	var err error

	log.Print(sr.cfg.UserPsw)
	log.Print(sr.cfg.Key)

	record.Text, err = sr.c.EncryptAES(sr.cfg.Key, record.Text)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}

	switch record.Operation {
	case models.Create:
		err := sr.sl.CreateText(ctx, record)
		if err != nil {
			log.Print("create text record error: ", err)
		}
		return err
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

func (sr *TextServices) SearchText(ctx context.Context, searchInput string) ([]models.TextRecord, error) {
	textRecords, err := sr.sl.SearchText(ctx, searchInput)
	if err != nil {
		log.Print("rearch text record error: ", err)
	}

	for i := range textRecords {

		textRecords[i].Text, err = sr.c.DecryptAES(sr.cfg.Key, textRecords[i].Text)
		if err != nil {
			log.Print("encrypt error: ", err)
			return nil, err
		}
	}

	return textRecords, err
}
