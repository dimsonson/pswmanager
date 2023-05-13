package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"

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
	cfg *config.ServiceConfig
	sl  StorageProvider
	c   CryptProvider
}

// New.
func NewUsers(s StorageProvider, cfg *config.ServiceConfig) *UserServices {
	return &UserServices{
		sl:  s,
		cfg: cfg,
		c:   &Crypt{},
	}
}

// TextRec.
func (sr *UserServices) CreateUser(ctx context.Context) error {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Print("generate key error: ", err)
		return err
	}

	keyString := hex.EncodeToString(key)

	psw256 := sha256.Sum256([]byte(sr.cfg.UserPsw))
	psw256string := hex.EncodeToString(psw256[:])

	keyDB, err := sr.c.EncryptAES(psw256string, keyString)
	if err != nil {
		log.Print("create keyDB error: ", err)
	}

	passHex, err := bcrypt.GenerateFromPassword([]byte(sr.cfg.UserPsw), bcrypt.DefaultCost)
	if err != nil {
		log.Print("generate hex error: ", err)
		return err
	}

	sr.cfg.UserPsw = string(passHex)
	sr.cfg.Key = keyString

	err = sr.sl.CreateUser(ctx, sr.cfg.UserConfig, keyDB)
	if err != nil {
		log.Print("create user error: ", err)
		return err
	}
	return err
}

func (sr *UserServices) ReadUser(ctx context.Context) (config.UserConfig, error) {
	ucfg, err := sr.sl.ReadUser(ctx)
	if err != nil {
		log.Print("read user cfg error: ", err)
	}

	psw256 := sha256.Sum256([]byte(sr.cfg.UserPsw))
	psw256string := hex.EncodeToString(psw256[:])

	ucfg.Key, err = sr.c.DecryptAES(psw256string, ucfg.Key)

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
