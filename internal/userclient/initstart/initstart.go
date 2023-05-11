package initstart

import (
	"context"
	"sync"
	"time"

	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/google/uuid"
	"github.com/derailed/tview"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/internal/userclient/services"
	"github.com/dimsonson/pswmanager/internal/userclient/storage"
	"github.com/dimsonson/pswmanager/internal/userclient/ui"
)

type Init struct {
	cfg *config.ServiceConfig
}

func New() *Init {
	return &Init{}
}

func (init *Init) InitAndStart(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup, uiLog *tview.TextView) {
	// создание конфигурацию сервера
	init.cfg = config.New()
	// парсинг конфигурации сервера
	init.cfg.Parse()

	sl, err := storage.New(init.cfg.SQLight.Dsn)
	if err != nil {
		log.Print("storage new error:", err)
	}

	srvusers := services.NewUsers(sl)

	init.cfg.UserConfig, err = srvusers.ReadUser(ctx)
	if err != nil {
		log.Print("storage new error:", err)
	}

	ui := ui.NewUI(ctx, init.cfg, srvusers)
	ui.Init(uiLog)
	//ui.LogWindow = uiLog
	//log.LogInit()
	//log.Output(ui.LogWindow)
	go ui.UIRun()

	//init.test(ctx)

	// srvtext := services.NewText(sl)
	// srvlogin := services.NewLogin(sl)
	// srvBinary := services.NewBinary(sl)
	// srcvCard := services.NewCard(sl)

	<-ctx.Done()
}

func (init *Init) ConnClose(ctx context.Context) {
	if init.cfg.SQLight.Conn != nil {
		init.cfg.SQLight.Conn.Close()
	}
}

func (init *Init) test(ctx context.Context) {

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

	resultTextSearch, err := srvtext.SearchText(ctx, "123")
	if err != nil {
		log.Print("search login_records error :", err)
	}

	delTxtMsg := testMsg
	delTxtMsg.Operation = models.Delete
	if err = srvtext.ProcessingText(ctx, delTxtMsg); err != nil {
		log.Print("storage new error:", err)
	}

	log.Print(resultTextSearch)

	srvlogin := services.NewLogin(sl)

	msgLogin := models.LoginRecord{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       "uid1",
		AppID:     "app1",
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

	resultSearch, err := srvlogin.SearchLogin(ctx, "12")
	if err != nil {
		log.Print("search login_records error :", err)
	}
	delLoginMsg := updateLoginMsg
	delLoginMsg.Operation = models.Delete
	if err = srvlogin.ProcessingLogin(ctx, delLoginMsg); err != nil {
		log.Print("storage new error:", err)
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

	binaryResultSearch, err := srvBinary.SearchBinary(ctx, "sample")
	if err != nil {
		log.Print("search login_records error :", err)
	}

	log.Print(binaryResultSearch)

	srcvCard := services.NewCard(sl)

	msgCard := models.CardRecord{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       "CardUID",
		AppID:     "AppIDCard",
		Brand:     1,
		ValidDate: "01/28",
		Number:    "2202245445789856",
		Code:      123,
		Holder:    "DMTIRY BO",
		Metadata:  "meta data card description sample",
		Operation: models.Create,
	}

	if err = srcvCard.ProcessingCard(ctx, msgCard); err != nil {
		log.Print("storage new error:", err)
	}

	updateCardMsg := msgCard
	updateCardMsg.ValidDate = "01/30"
	updateCardMsg.Metadata = "card123"
	updateCardMsg.Number = "2202245400000000"
	updateCardMsg.Operation = models.Update
	if err = srcvCard.ProcessingCard(ctx, updateCardMsg); err != nil {
		log.Print("storage new error:", err)
	}

	cardSearchResult, err := srcvCard.SearchCard(ctx, "card")
	if err != nil {
		log.Print("search login_records error :", err)
	}

	delCardMsg := updateCardMsg
	delCardMsg.Operation = models.Delete
	if err = srcvCard.ProcessingCard(ctx, delCardMsg); err != nil {
		log.Print("storage new error:", err)
	}

	log.Print(cardSearchResult)
}
