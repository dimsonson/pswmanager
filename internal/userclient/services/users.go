package services

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type UsersStorageProviver interface {
	CreateUser(ctx context.Context, ucfg config.UserConfig) error
	ReadUser(ctx context.Context) (config.UserConfig, error)
	CheckUser(ctx context.Context, login string) (string, error)
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
func (sr *UserServices) CreateUser(ctx context.Context, ucfg *config.UserConfig) error {
	passHex, err := bcrypt.GenerateFromPassword([]byte(ucfg.UserPsw), bcrypt.DefaultCost)
	if err != nil {
		log.Print("generate hex error: ", err)
	}
	ucfg.UserPsw = string(passHex)
	err = sr.sl.CreateUser(ctx, *ucfg)
	if err != nil {
		log.Print("create user error: ", err)
	}
	return err
}

func (sr *UserServices) ReadUser(ctx context.Context) (config.UserConfig, error) {
	ucfg, err := sr.sl.ReadUser(ctx)
	if err != nil {
		log.Print("read user cfg error: ", err)
	}
	return ucfg, err
}

func (sr *UserServices) CheckUser(ctx context.Context, ulogin, upsw string) error {
	passwDB, err := sr.sl.CheckUser(ctx, ulogin)
	if err != nil {
		log.Print("check user cfg error: ", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwDB), []byte(upsw))
	if err != nil {
		log.Print("check psw error: ", err)
	}
	return err
}
