package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/pkg/log"
)

type UsersStorageProviver interface {
	//CreateUser(ctx context.Context, ucfg config.UserConfig, keyDb string) error
	CreateUser(ctx context.Context, ucfg config.UserConfig) error
	ReadUser(ctx context.Context) (config.UserConfig, error)
	CheckUser(ctx context.Context, login string) (string, error)
}

type ClientGRPCProvider interface {
	NewUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	NewApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error)
	ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error)
	IsOnline() bool
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	cfg        *config.ServiceConfig
	clientGRPC ClientGRPCProvider
	sl         UsersStorageProviver
	c          CryptProvider
}

// NewUsers.
func NewUsers(s UsersStorageProviver, clientGRPC ClientGRPCProvider, cfg *config.ServiceConfig) *UserServices {
	return &UserServices{
		sl:         s,
		clientGRPC: clientGRPC,
		cfg:        cfg,
		c:          &Crypt{},
	}
}

// CreateUser метод создания профиля пользователя из данных конфигурации.
func (sr *UserServices) CreateUser(ctx context.Context) error {
	if !sr.clientGRPC.IsOnline() {
		log.Print("clientGRPC status: ", sr.cfg.GRPC.ClientConn.GetState().String())
		err := errors.New("clien gRPC status is not online")
		return err
	}

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

	psw512 := sha256.Sum256([]byte(sr.cfg.UserPsw))
	passHex := hex.EncodeToString(psw512[:])

	// passHex, err := bcrypt.GenerateFromPassword([]byte(sr.cfg.UserPsw), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Print("generate hex error: ", err)
	// 	return err
	// }

	// сохранение ключа шифрования и пароля в БД
	sr.cfg.UserPsw = passHex
	sr.cfg.Key = keyDB

	uConfig, err := sr.clientGRPC.NewUser(ctx, &pb.CreateUserRequest{
		Login: sr.cfg.UserLogin,
		Psw:   sr.cfg.UserPsw,
		CKey:  keyDB,
	})
	if err != nil {
		log.Print("gRPC call NewUser error: ", err)
		return err
	}
	l := len(uConfig.Apps) - 1
	if l < 0 {
		sr.cfg.UserID = uConfig.UserID

		log.Printf("login %s already exist", sr.cfg.UserLogin)
		uConfigApp, err := sr.clientGRPC.NewApp(ctx, &pb.CreateAppRequest{
			Uid: sr.cfg.UserID,
			Psw: sr.cfg.UserPsw,
		})
		if err != nil {
			log.Print("gRPC call NewUser error: ", err)
			return err
		}
		lApp := len(uConfigApp.Apps) - 1
		//sr.cfg.UserID = uConfig.UserID
		sr.cfg.Key = uConfigApp.CKey
		sr.cfg.AppID = uConfigApp.Apps[lApp].AppID
		sr.cfg.ExchName = uConfigApp.ExchangeName
		sr.cfg.RoutingKey = uConfigApp.Apps[lApp].RoutingKey
		sr.cfg.ConsumeQueue = uConfigApp.Apps[lApp].ConsumeQueue
		sr.cfg.ConsumeRkey = uConfigApp.UserID + ".*.*"

	}
	if l >= 0 {
		sr.cfg.UserID = uConfig.UserID
		sr.cfg.AppID = uConfig.Apps[l].AppID
		sr.cfg.ExchName = uConfig.ExchangeName
		sr.cfg.RoutingKey = uConfig.Apps[l].RoutingKey
		sr.cfg.ConsumeQueue = uConfig.Apps[l].ConsumeQueue
		sr.cfg.ConsumeRkey = uConfig.UserID + ".*.*"
	}
	log.Print(sr.cfg.UserConfig)

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

	log.Print("psw256string", psw256string)

	// сохранение в память расшифрованного ключа пользователя
	ucfg.Key, err = sr.c.DecryptAES(psw256string, ucfg.Key)
	if err != nil {
		log.Print("dicrypt user key error: ", err)
	}

	log.Print("Read User ucfg.Key", ucfg.Key)

	return ucfg, err
}

// CheckUser проверка пароля пользователя с записю в хранилище.
func (sr *UserServices) CheckUser(ctx context.Context, ulogin, upsw string) error {
	passwDB, err := sr.sl.CheckUser(ctx, ulogin)
	if err != nil {
		log.Print("check user cfg error: ", err)
	}

	psw512 := sha256.Sum256([]byte(upsw))
	passHex := hex.EncodeToString(psw512[:])

	if passHex != passwDB {
		return errors.New("wrong password or login")
	}

	// err = bcrypt.CompareHashAndPassword([]byte(passwDB), []byte(upsw))
	// if err != nil {
	// 	log.Print("check psw error: ", err)
	// }
	return nil
}
