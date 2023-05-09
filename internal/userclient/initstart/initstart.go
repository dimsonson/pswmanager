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

	delTxtMsg := testMsg
	delTxtMsg.Operation = models.Delete
	if err = srvtext.ProcessingText(ctx, delTxtMsg); err != nil {
		log.Print("storage new error:", err)
	}

	srvlogin := services.NewLogin(sl)




	msgLogin := models.LoginRecord{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       init.cfg.UserID,
		AppID:     init.cfg.AppID,
		Login:     "login0001",
		Psw:       "password001",
		Metadata:  "meta data description sample",
		Operation: models.Create,
	}
	if err = srvlogin.ProcessingLogin(ctx, msgLogin); err != nil {
		log.Print("storage new error:", err)
	}



	// msgBinary := models.BinaryRecord{
	// 	RecordID:  uuid.NewString(),
	// 	ChngTime:  time.Now(),
	// 	UID:       newUserCfg.UserID,
	// 	AppID:     newAppCfg.Appid,
	// 	Binary:    "secured text sending",
	// 	Metadata:  "meta data description sample",
	// 	Operation: models.Create,
	// }

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
