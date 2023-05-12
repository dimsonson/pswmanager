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
	ucfg config.UserConfig
	sl   StorageProvider
	c    CryptProvider
}

// New.
func NewText(s StorageProvider, ucfg config.UserConfig) *TextServices {
	return &TextServices{
		sl:   s,
		ucfg: ucfg,
		c: &Crypt{},
	}
}

// TextRec.
func (sr *TextServices) ProcessingText(ctx context.Context, record models.TextRecord) error {
	var err error

	log.Print(sr.ucfg.UserPsw)
	log.Print(sr.ucfg.Key)

	//key:= string([]byte(sr.ucfg.UserPsw)[1:])   +"00000"
	record.Text, err = sr.c.EncryptAES(sr.ucfg.Key, record.Text)
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
	return textRecords, err
}
