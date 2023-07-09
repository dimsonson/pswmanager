package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"

	pbconsume "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protoconsume"
	pbpub "github.com/dimsonson/pswmanager/internal/gateway/handlers/grpc_handlers/protopub"
	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
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
	UpdateText(ctx context.Context, record models.TextRecord) error
	DeleteText(ctx context.Context, record models.TextRecord) error

	MarkLoginSent(ctx context.Context, record models.LoginRecord) error
	CreateLogin(ctx context.Context, record models.LoginRecord) error
	UpdateLogin(ctx context.Context, record models.LoginRecord) error
	DeleteLogin(ctx context.Context, record models.LoginRecord) error

	MarkBinarySent(ctx context.Context, record models.BinaryRecord) error
	CreateBinary(ctx context.Context, record models.BinaryRecord) error
	UpdateBinary(ctx context.Context, record models.BinaryRecord) error
	DeleteBinary(ctx context.Context, record models.BinaryRecord) error

	MarkCardSent(ctx context.Context, record models.CardRecord) error
	CreateCard(ctx context.Context, record models.CardRecord) error
	UpdateCard(ctx context.Context, record models.CardRecord) error
	DeleteCard(ctx context.Context, record models.CardRecord) error
}

type ClientGRPCProvider interface {
	NewUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	NewApp(ctx context.Context, in *pb.CreateAppRequest) (*pb.CreateAppResponse, error)
	ReadUser(ctx context.Context, in *pb.ReadUserRequest) (*pb.ReadUserResponse, error)
	IsOnline() bool

	PublishText(ctx context.Context, in *pbpub.PublishTextRequest) error
	PublishLogins(ctx context.Context, in *pbpub.PublishLoginsRequest) error
	PublishBinary(ctx context.Context, in *pbpub.PublishBinaryRequest) error
	PublishCard(ctx context.Context, in *pbpub.PublishCardRequest) error

	ConsumeFromStream(ctx context.Context, in *pbconsume.ConsumeRequest) (pbconsume.ServerRMQhandlers_ConsumeClient, error)
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
		sr.cfg.Key = uConfigApp.CKey
		sr.cfg.AppID = uConfigApp.Apps[lApp].AppID
		sr.cfg.ExchName = uConfigApp.ExchangeName
		sr.cfg.RoutingKey = uConfigApp.Apps[lApp].RoutingKey
		sr.cfg.ConsumeQueue = uConfigApp.Apps[lApp].ConsumeQueue
		sr.cfg.ConsumeRkey = uConfigApp.UserID + ".*.*"
		recordsApp, err := sr.clientGRPC.ReadUser(ctx, &pb.ReadUserRequest{
			Uid: sr.cfg.UserID,
		})
		if err != nil {
			log.Print("gRPC call ReadUser error: ", err)
			return err
		}
		textRecord := models.TextRecord{}
		for i := range recordsApp.SetTextRec {
			textRecord.RecordID = recordsApp.SetTextRec[i].RecordID
			textRecord.ChngTime = recordsApp.SetTextRec[i].ChngTime.AsTime()
			textRecord.UID = recordsApp.SetTextRec[i].UID
			textRecord.AppID = recordsApp.SetTextRec[i].AppID
			textRecord.Text = recordsApp.SetTextRec[i].Text
			textRecord.Metadata = recordsApp.SetTextRec[i].Metadata
			textRecord.Operation = models.Create
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

// ConsumeFromStream.
func (sr *UserServices) ConsumeFromStream(ctx context.Context) {
	log.Print(sr.cfg.ExchName, "\n", sr.cfg.ConsumeQueue, "\n", sr.cfg.ConsumeRkey)
	stream, err := sr.clientGRPC.ConsumeFromStream(ctx, &pbconsume.ConsumeRequest{
		ExchName:      sr.cfg.ExchName,
		ConsumerQname: sr.cfg.ConsumeQueue,
		RoutingKey:    "all." + sr.cfg.ConsumeRkey,
	})

	log.Print(sr.cfg.ExchName, sr.cfg.ConsumeQueue, sr.cfg.RoutingKey)
	if err != nil {
		log.Print("stream error: ", err)
	}
	var rec *pbconsume.ConsumeResponse
	for {
		select {
		case <-ctx.Done():
			log.Print(ctx.Err()) // prints "context deadline exceeded"
			return
		default:
			log.Print("starting consume worker...")
			rec, err = stream.Recv()
			if err != nil {
				log.Print("stream error: ", err)
				return
			}
			//log.Print(rec)
			err := sr.StreamDataProcessing(ctx, int(rec.RecordType), rec.Record)
			if err != nil {
				log.Print("stream peocessing error: ", err)
				stream.SendMsg(err)
			}
		}
	}
}

// StreamDataProcessing.
func (sr *UserServices) StreamDataProcessing(ctx context.Context, typeOfRecord int, record []byte) error {
	var err error
	
	switch typeOfRecord {
	case int(models.TextType):
		var textrec models.TextRecord
		err = json.Unmarshal(record, &textrec)
		if err != nil {
			log.Print("unmarshal error: ", err)
			return err
		}
		switch textrec.Operation {
		case models.Create:
			err = sr.sl.CreateText(ctx, textrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		case models.Update:
			err = sr.sl.UpdateText(ctx, textrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		case models.Delete:
			err = sr.sl.DeleteText(ctx, textrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		default:
			log.Print("unknown operation type for stream data processing", err)
			return errors.New("unknown type for stream data processing")
		}
	case int(models.LoginsType):
		var loginrec models.LoginRecord
		err = json.Unmarshal(record, &loginrec)
		if err != nil {
			log.Print("unmarshal error: ", err)
			return err
		}
		switch loginrec.Operation {
		case models.Create:
			err = sr.sl.CreateLogin(ctx, loginrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		case models.Update:
			err = sr.sl.UpdateLogin(ctx, loginrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		case models.Delete:
			err = sr.sl.DeleteLogin(ctx, loginrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		default:
			log.Print("unknown operation type for stream data processing", err)
			return errors.New("unknown type for stream data processing")
		}
	case int(models.BinaryType):
		var binaryrec models.BinaryRecord
		err = json.Unmarshal(record, &binaryrec)
		if err != nil {
			log.Print("unmarshal error: ", err)
			return err
		}
		switch binaryrec.Operation {
		case models.Create:
			err = sr.sl.CreateBinary(ctx, binaryrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		case models.Update:
			err = sr.sl.UpdateBinary(ctx, binaryrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		case models.Delete:
			err = sr.sl.DeleteBinary(ctx, binaryrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		default:
			log.Print("unknown operation type for stream data processing", err)
			return errors.New("unknown type for stream data processing")
		}
	case int(models.CardType):
		var cardrec models.CardRecord
		err = json.Unmarshal(record, &cardrec)
		if err != nil {
			log.Print("unmarshal error: ", err)
			return err
		}
		switch cardrec.Operation {
		case models.Create:
			err = sr.sl.CreateCard(ctx, cardrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		case models.Update:
			err = sr.sl.UpdateCard(ctx, cardrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		case models.Delete:
			err = sr.sl.DeleteCard(ctx, cardrec)
			if err != nil {
				log.Print("processing error: ", err)
				return err
			}
		default:
			log.Print("unknown operation type for stream data processing", err)
			return errors.New("unknown type for stream data processing")
		}
	default:
		return errors.New("unknown type for stream data processing")
	}
	return err
}
