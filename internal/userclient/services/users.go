package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type UsersStorageProviver interface {
	CreateUser(ctx context.Context, ucfg config.UserConfig, keyDb string) error
	ReadUser(ctx context.Context) (config.UserConfig, error)
	CheckUser(ctx context.Context, login string) (string, error)
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	ucfg config.UserConfig
	sl   StorageProvider
	c    CryptProvider
}

// New.
func NewUsers(s StorageProvider, ucfg config.UserConfig) *UserServices {
	return &UserServices{
		sl:   s,
		ucfg: ucfg,
		c:    &Crypt{},
	}
}

// TextRec.
func (sr *UserServices) CreateUser(ctx context.Context, ucfg *config.UserConfig) error {
	// key := make([]byte, 32)
	// if _, err := io.ReadFull(rand.Reader, key); err != nil {
	// 	panic(err.Error())
	// }
	psw256 := sha256.Sum256([]byte(ucfg.UserPsw))
	keyDB := hex.EncodeToString(psw256[:])
	//keyDB, err := sr.c.EncryptAES(p, string(key))
	// if err != nil {
	// 	log.Print("create keyDB error: ", err)
	// }

	passHex, err := bcrypt.GenerateFromPassword([]byte(ucfg.UserPsw), bcrypt.DefaultCost)
	if err != nil {
		log.Print("generate hex error: ", err)
	}

	ucfg.UserPsw = string(passHex)

	ucfg.Key = keyDB

	err = sr.sl.CreateUser(ctx, *ucfg, keyDB)
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
