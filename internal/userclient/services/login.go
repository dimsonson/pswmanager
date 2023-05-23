package services

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
)

type LoginStorageProviver interface {
	CreateLogin(ctx context.Context, record models.LoginRecord) error
	UpdateLogin(ctx context.Context, record models.LoginRecord) error
	DeleteLogin(ctx context.Context, record models.LoginRecord) error
	SearchLogin(ctx context.Context, searchInput string) ([]models.LoginRecord, error)
	MarkLoginSent(ctx context.Context, record models.LoginRecord) error
}

// Services структура конструктора бизнес логики.
type LoginServices struct {
	cfg *config.ServiceConfig
	sl  LoginStorageProviver
	c   CryptProvider
}

// New.
func NewLogin(s LoginStorageProviver, cfg *config.ServiceConfig) *LoginServices {
	return &LoginServices{
		sl:  s,
		cfg: cfg,
		c:   &Crypt{},
	}
}

// LoginRec.
func (sr *LoginServices) ProcessingLogin(ctx context.Context, record models.LoginRecord) error {
	var err error
	record.Login, err = sr.c.EncryptAES(sr.cfg.Key, record.Login)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	record.Psw, err = sr.c.EncryptAES(sr.cfg.Key, record.Psw)
	if err != nil {
		log.Print("encrypt error: ", err)
		return err
	}
	switch record.Operation {
	case models.Create:
		err := sr.sl.CreateLogin(ctx, record)
		if err != nil {
			log.Print("create login record error: ", err)
		}
		return err
	case models.Update:
		err := sr.sl.UpdateLogin(ctx, record)
		if err != nil {
			log.Print("update login record error: ", err)
		}
		return err
	case models.Delete:
		err := sr.sl.DeleteLogin(ctx, record)
		if err != nil {
			log.Print("delete login record error: ", err)
		}
		return err
	}
	return err
}

func (sr *LoginServices) SearchLogin(ctx context.Context, searchInput string) ([]models.LoginRecord, error) {
	loginRecords, err := sr.sl.SearchLogin(ctx, searchInput)
	if err != nil {
		log.Print("rearch login record error: ", err)
	}
	for i := range loginRecords {
		loginRecords[i].Login, err = sr.c.DecryptAES(sr.cfg.Key, loginRecords[i].Login)
		if err != nil {
			log.Print("encrypt error: ", err)
			return nil, err
		}
		loginRecords[i].Psw, err = sr.c.DecryptAES(sr.cfg.Key, loginRecords[i].Psw)
		if err != nil {
			log.Print("encrypt error: ", err)
			return nil, err
		}
	}
	return loginRecords, err
}
