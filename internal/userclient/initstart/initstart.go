package initstart

import (
	"context"
	"sync"
	"time"

	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/google/uuid"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/internal/userclient/services"
	"github.com/dimsonson/pswmanager/internal/userclient/storage"
)

type Init struct {
	cfg *config.ServiceConfig
}

func New(cfg *config.ServiceConfig) *Init {
	return &Init{
		cfg: cfg,
	}
}

func (init *Init) InitAndStart(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup) {

	sl, err := storage.New(init.cfg.SQLight.Dsn)
	if err != nil {
		log.Print("storage new error:", err)
	}

	srvtext := services.NewText(sl)

	testMsg := models.TextRecord{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       "uidtest",
		AppID:     "AppId Test",
		Text:      "secured text sending",
		Metadata:  "meta data description sample",
		Operation: models.Create,
	}
	if err = srvtext.ProcessingText(ctx, testMsg); err != nil {
		log.Print("storage new error:", err)
	}

	updateTxtMsg := testMsg
	updateTxtMsg.Text = "update secured text sending"
	updateTxtMsg.Metadata = "123"
	updateTxtMsg.Operation = models.Update
	if err = srvtext.ProcessingText(ctx, updateTxtMsg); err != nil {
		log.Print("storage new error:", err)
	}

	// delTxtMsg := testMsg
	// delTxtMsg.Operation = models.Delete
	// if err = srvtext.ProcessingText(ctx, delTxtMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	resultTextSearch, err := srvtext.SearchText(ctx, "123")
	if err != nil {
		log.Print("search login_records error :", err)
	}

	log.Print(resultTextSearch)

	srvlogin := services.NewLogin(sl)

	msgLogin := models.LoginRecord{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       "uid1", //init.cfg.UserID,
		AppID:     "app1", //init.cfg.AppID,
		Login:     "login0001",
		Psw:       "password001",
		Metadata:  "meta data description sample",
		Operation: models.Create,
	}
	if err = srvlogin.ProcessingLogin(ctx, msgLogin); err != nil {
		log.Print("storage new error:", err)
	}

	updateLoginMsg := msgLogin
	updateLoginMsg.Login = "update login"
	updateLoginMsg.Metadata = "123update"
	updateLoginMsg.Operation = models.Update
	if err = srvlogin.ProcessingLogin(ctx, updateLoginMsg); err != nil {
		log.Print("storage new error:", err)
	}

	// delLoginMsg := updateLoginMsg
	// delLoginMsg.Operation = models.Delete
	// if err = srvlogin.ProcessingLogin(ctx, delLoginMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	resultSearch, err := srvlogin.SearchLogin(ctx, "12")
	if err != nil {
		log.Print("search login_records error :", err)
	}

	log.Print(resultSearch)

	srvBinary := services.NewBinary(sl)

	msgBinary := models.BinaryRecord{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       "newUserCfg.UserID",
		AppID:     "newAppCfg.Appid",
		Binary:    "secured binary sending",
		Metadata:  "meta data description sample",
		Operation: models.Create,
	}

	if err = srvBinary.ProcessingBinary(ctx, msgBinary); err != nil {
		log.Print("storage new error:", err)
	}

	updateBinaryMsg := msgBinary
	updateBinaryMsg.Binary = "update binary"
	updateBinaryMsg.Metadata = "123updateBinary"
	updateBinaryMsg.Operation = models.Update
	if err = srvBinary.ProcessingBinary(ctx, updateBinaryMsg); err != nil {
		log.Print("storage new error:", err)
	}

	delBinaryMsg := updateBinaryMsg
	delBinaryMsg.Operation = models.Delete
	if err = srvBinary.ProcessingBinary(ctx, delBinaryMsg); err != nil {
		log.Print("storage new error:", err)
	}




	// msgCard := models.CardRecord{
	// 	RecordID:  uuid.NewString(),
	// 	ChngTime:  time.Now(),
	// 	UID:       newUserCfg.UserID,
	// 	AppID:     newAppCfg.Appid,
	// 	Brand:     1,
	// 	ValidDate: "01/28",
	// 	Number:    "2202245445789856",
	// 	Code:      123,
	// 	Holder:    "DMTIRY BO",
	// 	Metadata:  "meta data card description sample",
	// 	Operation: models.Create,
	// }

	<-ctx.Done()
}

func (init *Init) ConnClose(ctx context.Context) {
	if init.cfg.SQLight.Conn != nil {
		init.cfg.SQLight.Conn.Close()
	}
}
