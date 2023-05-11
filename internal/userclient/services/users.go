package services

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/pkg/log"
)

type UsersStorageProviver interface {
	CreateUser(ctx context.Context, ucfg config.UserConfig) error
	ReadUser(ctx context.Context) (config.UserConfig, error)
	CheckUser(ctx context.Context, login string, passwHex string) error
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	sl StorageProvider
}

// New.
func NewUsers(s StorageProvider) *UserServices {
	return &UserServices{
		s,
	}
}

// TextRec.
func (sr *UserServices) CreateUser(ctx context.Context, ucfg config.UserConfig) error {
	var err error
	// switch record.Operation {
	// case models.Create:
	// 	err := sr.sl.CreateText(ctx, record)
	// 	if err != nil {
	// 		log.Print("create text record error: ", err)
	// 	}
	// 	return err
	// case models.Update:
	// 	err := sr.sl.UpdateText(ctx, record)
	// 	if err != nil {
	// 		log.Print("update text record error: ", err)
	// 	}
	// 	return err
	// case models.Delete:
	// 	err := sr.sl.DeleteText(ctx, record)
	// 	if err != nil {
	// 		log.Print("delete text record error: ", err)
	// 	}
	// 	return err
	// }
	return err
}

func (sr *UserServices) ReadUser(ctx context.Context) (config.UserConfig, error) {
	ucfg, err := sr.sl.ReadUser(ctx)
	if err != nil {
		log.Print("rearch text record error: ", err)
	}
	return ucfg, err
}

func (sr *UserServices) CheckUser(ctx context.Context, ulogin string, passw string) error {
	textRecords, err := sr.sl.SearchText(ctx, ulogin)
	if err != nil {
		log.Print("rearch text record error: ", err)
	}
	_ = textRecords
	return err
}
