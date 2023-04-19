package services

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type TextStorageProviver interface {
	CreateText(ctx context.Context, record models.TextRec) error
	UpdateText(ctx context.Context, record models.TextRec) error
	DeleteText(ctx context.Context, record models.TextRec) error
}

// Services структура конструктора бизнес логики.
type TextServices struct {
	storage StorageProvider
}

// New.
func NewTextRec(s StorageProvider) *TextServices {
	return &TextServices{
		s,
	}
}

// TextRec.
func (sr *TextServices) TextRec(ctx context.Context, record models.TextRec) error {
	var err error
	switch record.Operation {
	case models.Create:
		err := sr.storage.CreateText(ctx, record)
		if err != nil {
			log.Print("create text record error: ", err)
		}
		return err
	case models.Update:
		err := sr.storage.UpdateText(ctx, record)
		if err != nil {
			log.Print("update text record error: ", err)
		}
		return err
	case models.Delete:
		err := sr.storage.DeleteText(ctx, record)
		if err != nil {
			log.Print("delete text record error: ", err)
		}
		return err
	}
	return err
}
