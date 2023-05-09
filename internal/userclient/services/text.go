package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type TextStorageProviver interface {
	CreateText(ctx context.Context, record models.TextRecord) error
	UpdateText(ctx context.Context, record models.TextRecord) error
	DeleteText(ctx context.Context, record models.TextRecord) error
}

// Services структура конструктора бизнес логики.
type TextServices struct {
	sl StorageProvider
}

// New.
func NewText(s StorageProvider) *TextServices {
	return &TextServices{
		s,
	}
}

// TextRec.
func (sr *TextServices) ProcessingText(ctx context.Context, record models.TextRecord) error {
	var err error
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