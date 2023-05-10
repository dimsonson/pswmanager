package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type LoginStorageProviver interface {
	CreateLogin(ctx context.Context, record models.LoginRecord) error
	UpdateLogin(ctx context.Context, record models.LoginRecord) error
	DeleteLogin(ctx context.Context, record models.LoginRecord) error
	SearchLogin(ctx context.Context, searchInput string) ([]models.LoginRecord, error)
}

// Services структура конструктора бизнес логики.
type LoginServices struct {
	storage LoginStorageProviver
}

// New.
func NewLogin(s LoginStorageProviver) *LoginServices {
	return &LoginServices{
		s,
	}
}

// LoginRec.
func (sr *LoginServices) ProcessingLogin(ctx context.Context, record models.LoginRecord) error {
	var err error
	switch record.Operation {
	case models.Create:
		err := sr.storage.CreateLogin(ctx, record)
		if err != nil {
			log.Print("create login record error: ", err)
		}
		return err
	case models.Update:
		err := sr.storage.UpdateLogin(ctx, record)
		if err != nil {
			log.Print("update login record error: ", err)
		}
		return err
	case models.Delete:
		err := sr.storage.DeleteLogin(ctx, record)
		if err != nil {
			log.Print("delete login record error: ", err)
		}
		return err
	}
	return err
}

func (sr *LoginServices) SearchLogin(ctx context.Context, searchInput string) ([]models.LoginRecord, error) {
	loginRecords, err := sr.storage.SearchLogin(ctx, searchInput)
	if err != nil {
		log.Print("rearch login record error: ", err)
	}
	return loginRecords, err
}
