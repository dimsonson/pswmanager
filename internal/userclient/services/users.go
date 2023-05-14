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
	//CreateUser(ctx context.Context, ucfg config.UserConfig, keyDb string) error
	CreateUser(ctx context.Context, ucfg config.UserConfig) error
	ReadUser(ctx context.Context) (config.UserConfig, error)
	CheckUser(ctx context.Context, login string) (string, error)
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	cfg *config.ServiceConfig
	sl  UsersStorageProviver
	c   CryptProvider
}

// NewUsers.
func NewUsers(s UsersStorageProviver, cfg *config.ServiceConfig) *UserServices {
	return &UserServices{
		sl:  s,
		cfg: cfg,
		c:   &Crypt{},
	}
}

// CreateUser метод создания профиля пользователя из данных конфигурации.
func (sr *UserServices) CreateUser(ctx context.Context) error {
	// генереация ключа шифрования данных
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Print("generate key error: ", err)
		return err
	}
	keyString := hex.EncodeToString(key)
	// генерация из пароля ключа, которым будет зашифрован ключ шифрования данных
	psw256 := sha256.Sum256([]byte(sr.cfg.UserPsw))
	psw256string := hex.EncodeToString(psw256[:])
	// шифрование ключа шифрования данных для хранения в базе
	keyDB, err := sr.c.EncryptAES(psw256string, keyString)
	if err != nil {
		log.Print("create keyDB error: ", err)
	}
	// создание хеша пароля пользователя для хранения в базе данных
	passHex, err := bcrypt.GenerateFromPassword([]byte(sr.cfg.UserPsw), bcrypt.DefaultCost)
	if err != nil {
		log.Print("generate hex error: ", err)
		return err
	}
	// сохранение ключа шифрования и пароля в БД
	sr.cfg.UserPsw = string(passHex)
	sr.cfg.Key = keyDB
	err = sr.sl.CreateUser(ctx, sr.cfg.UserConfig)
	if err != nil {
		log.Print("create user error: ", err)
		return err
	}
	return err
}

// ReadUser чтение профиля пользователя из БД.
func (sr *UserServices) ReadUser(ctx context.Context) (config.UserConfig, error) {
	// получаем профиль из хранилища
	ucfg, err := sr.sl.ReadUser(ctx)
	if err != nil {
		log.Print("read user cfg error: ", err)
	}
	// генерация из пароля ключа, которым будет расшифрован ключ шифрования данных
	psw256 := sha256.Sum256([]byte(sr.cfg.UserPsw))
	psw256string := hex.EncodeToString(psw256[:])
	// сохранение в память расшифрованного ключа пользователя
	ucfg.Key, err = sr.c.DecryptAES(psw256string, ucfg.Key)
	return ucfg, err
}

// CheckUser проверка пароля пользователя с записю в хранилище.
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
