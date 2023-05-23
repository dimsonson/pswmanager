package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/pkg/log"
)

type UsersStorageProviver interface {
	//CreateUser(ctx context.Context, ucfg config.UserConfig, keyDb string) error
	CreateUser(ctx context.Context, ucfg config.UserConfig) error
	ReadUser(ctx context.Context) (config.UserConfig, error)
	CheckUser(ctx context.Context, login string) (string, error)
	ReadAppLogin(ctx context.Context) (string, error)

	MarkTextSent(ctx context.Context, record models.TextRecord) error
	CreateText(ctx context.Context, record models.TextRecord) error

	MarkLoginSent(ctx context.Context, record models.LoginRecord) error
	CreateLogin(ctx context.Context, record models.LoginRecord) error

	MarkBinarySent(ctx context.Context, record models.BinaryRecord) error
	CreateBinary(ctx context.Context, record models.BinaryRecord) error

	MarkCardSent(ctx context.Context, record models.CardRecord) error
	CreateCard(ctx context.Context, record models.CardRecord) error
}

type ClientGRPCProvider interface {
	NewUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	NewApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error)
	ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error)
	IsOnline() bool
	CreateText(ctx context.Context, in *pbpub.PublishTextRequest) error
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	cfg        *config.ServiceConfig
	clientGRPC ClientGRPCProvider
	sl         UsersStorageProviver
	Crypt
}

// NewUsers.
func NewUsers(s UsersStorageProviver, clientGRPC ClientGRPCProvider, cfg *config.ServiceConfig) *UserServices {
	return &UserServices{
		sl:         s,
		clientGRPC: clientGRPC,
		cfg:        cfg,
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
	keyDB, err := sr.EncryptAES(psw256string, keyString)
	if err != nil {
		log.Print("create keyDB error: ", err)
	}
	// создание хеша пароля пользователя для хранения в базе данных

	psw := sha256.Sum256([]byte(sr.cfg.UserPsw))
	passHex := hex.EncodeToString(psw[:])

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

	log.Print("uConfig: ", uConfig)

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

		log.Print("sr.cfg.UserID: ", sr.cfg.UserID)

		recordsApp, err := sr.clientGRPC.ReadUser(ctx, &pb.ReadUserRequest{
			Uid: sr.cfg.UserID,
		})
		if err != nil {
			log.Print("gRPC call ReadUser error: ", err)
			return err
		}

		log.Print("recordsApp: ", recordsApp)

		log.Print("sr.cfg.Key: ", sr.cfg.Key)

		textRecord := models.TextRecord{}
		for i := range recordsApp.SetTextRec {
			textRecord.RecordID = recordsApp.SetTextRec[i].RecordID
			textRecord.ChngTime = recordsApp.SetTextRec[i].ChngTime.AsTime()
			textRecord.UID = recordsApp.SetTextRec[i].UID
			textRecord.AppID = recordsApp.SetTextRec[i].AppID
			textRecord.Text = recordsApp.SetTextRec[i].Text
			textRecord.Metadata = recordsApp.SetTextRec[i].Metadata
			textRecord.Operation = models.Create

			log.Print("textRecord  ", textRecord)

			textRecord.Text, err = sr.EncryptAES(keyString, textRecord.Text)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			err = sr.sl.CreateText(ctx, textRecord)
			if err != nil {
				log.Print("create text record error: ", err)
				return err
			}
			err = sr.sl.MarkTextSent(ctx, textRecord)
			if err != nil {
				log.Print("mark Text sent for new App error: ", err)
				return err
			}
		}

		loginRecord := models.LoginRecord{}
		for i := range recordsApp.SetLoginRec {
			loginRecord.RecordID = recordsApp.SetLoginRec[i].RecordID
			loginRecord.ChngTime = recordsApp.SetLoginRec[i].ChngTime.AsTime()
			loginRecord.UID = recordsApp.SetLoginRec[i].UID
			loginRecord.AppID = recordsApp.SetLoginRec[i].AppID
			loginRecord.Login = recordsApp.SetLoginRec[i].Login
			loginRecord.Psw = recordsApp.SetLoginRec[i].Psw
			loginRecord.Metadata = recordsApp.SetLoginRec[i].Metadata
			loginRecord.Operation = models.Create
			loginRecord.Login, err = sr.EncryptAES(keyString, loginRecord.Login)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			loginRecord.Psw, err = sr.EncryptAES(keyString, loginRecord.Psw)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			err = sr.sl.CreateLogin(ctx, loginRecord)
			if err != nil {
				log.Print("create text record error: ", err)
				return err
			}
			err = sr.sl.MarkLoginSent(ctx, loginRecord)
			if err != nil {
				log.Print("mark Text sent for new App error: ", err)
				return err
			}
		}

		binaryRecord := models.BinaryRecord{}
		for i := range recordsApp.SetBinaryRec {
			binaryRecord.RecordID = recordsApp.SetBinaryRec[i].RecordID
			binaryRecord.ChngTime = recordsApp.SetBinaryRec[i].ChngTime.AsTime()
			binaryRecord.UID = recordsApp.SetBinaryRec[i].UID
			binaryRecord.AppID = recordsApp.SetBinaryRec[i].AppID
			binaryRecord.Binary = recordsApp.SetBinaryRec[i].Binary
			binaryRecord.Metadata = recordsApp.SetBinaryRec[i].Metadata
			binaryRecord.Operation = models.Create
			binaryRecord.Binary, err = sr.EncryptAES(keyString, binaryRecord.Binary)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			err = sr.sl.CreateBinary(ctx, binaryRecord)
			if err != nil {
				log.Print("create text record error: ", err)
				return err
			}
			err = sr.sl.MarkBinarySent(ctx, binaryRecord)
			if err != nil {
				log.Print("mark Text sent for new App error: ", err)
				return err
			}
		}

		cardRecord := models.CardRecord{}
		for i := range recordsApp.SetCardRec {
			cardRecord.RecordID = recordsApp.SetCardRec[i].RecordID
			cardRecord.ChngTime = recordsApp.SetCardRec[i].ChngTime.AsTime()
			cardRecord.UID = recordsApp.SetCardRec[i].UID
			cardRecord.AppID = recordsApp.SetCardRec[i].AppID
			cardRecord.Brand = recordsApp.SetCardRec[i].Brand
			cardRecord.Number = recordsApp.SetCardRec[i].Number
			cardRecord.ValidDate = recordsApp.SetCardRec[i].ValidDate
			cardRecord.Code = recordsApp.SetCardRec[i].Code
			cardRecord.Holder = recordsApp.SetCardRec[i].Holder
			cardRecord.Metadata = recordsApp.SetCardRec[i].Metadata
			cardRecord.Operation = models.Create
			cardRecord.Brand, err = sr.EncryptAES(keyString, cardRecord.Brand)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			cardRecord.Number, err = sr.EncryptAES(keyString, cardRecord.Number)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			cardRecord.ValidDate, err = sr.EncryptAES(keyString, cardRecord.ValidDate)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			cardRecord.Code, err = sr.EncryptAES(keyString, cardRecord.Code)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			cardRecord.Holder, err = sr.EncryptAES(keyString, cardRecord.Holder)
			if err != nil {
				log.Print("encrypt error: ", err)
				return err
			}
			err = sr.sl.CreateCard(ctx, cardRecord)
			if err != nil {
				log.Print("create text record error: ", err)
				return err
			}
			err = sr.sl.MarkCardSent(ctx, cardRecord)
			if err != nil {
				log.Print("mark Text sent for new App error: ", err)
				return err
			}
		}
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

	log.Print("ucfg.Key bsfore decrypt", ucfg.Key)

	// сохранение в память расшифрованного ключа пользователя
	ucfg.Key, err = sr.DecryptAES(psw256string, ucfg.Key)
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

// IsAppRegistered.
func (sr *UserServices) IsAppRegistered(ctx context.Context) (bool, error) {
	ulogin, err := sr.sl.ReadAppLogin(ctx)
	if err != nil {
		log.Print("check app registration error: ", err)
	}
	if ulogin == "" {
		return false, err
	}
	return true, err
}
