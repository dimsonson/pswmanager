package initstart

import (
	"context"
	"sync"

	"github.com/derailed/tview"
	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/userclient/client/clientgrpc"
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
	// создание конфигураци клиента
	init.cfg = config.New()
	// парсинг конфигурациии клиента
	init.cfg.Parse()

	sl, err := storage.New(init.cfg.SQLight.Dsn)
	if err != nil {
		log.Print("storage new error:", err)
	}

	clientGRPC, err := clientgrpc.NewClientGRPC(&init.cfg.GRPC)
	if err != nil {
		log.Printf("new gRPC client error: %s", err)
		return
	}

	// init.cfg.UserConfig, err = srvusers.ReadUser(ctx)
	// if err != nil {
	// 	log.Print("no user data exist:", err)
	// }

	log.Print(string(init.cfg.UserConfig.Key))

	srvtext := services.NewText(sl, init.cfg)
	srvlogin := services.NewLogin(sl, init.cfg)
	srvbinary := services.NewBinary(sl, init.cfg)
	srvcard := services.NewCard(sl, init.cfg)

	srvusers := services.NewUsers(sl, clientGRPC, init.cfg)
	// testRSearchResults, err := srvtext.SearchText(ctx, "test")

	// log.Print(testRSearchResults)

	// init.cfg.UserID = "userID123456789"
	// init.cfg.AppID = "appID123456789"

	// init.cfg.UserLogin = "userlogin12345"
	// init.cfg.UserPsw = "userpassw12345"

	// err = srvusers.CreateUser(ctx, &init.cfg.UserConfig)
	// if err != nil {
	// 	log.Print("create user error:", err)
	// }

	log.Print(init.cfg.UserLogin)
	log.Print(init.cfg.UserPsw)

	ui := ui.NewUI(ctx, init.cfg, srvusers, srvtext, srvlogin, srvbinary, srvcard)
	ui.Init(uiLog)
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

	// sl, err := storage.New(init.cfg.SQLight.Dsn)
	// if err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// srvtext := services.NewText(sl, init.cfg)

	// testMsg := models.TextRecord{
	// 	RecordID:  uuid.NewString(),
	// 	ChngTime:  time.Now(),
	// 	UID:       "uidtest",
	// 	AppID:     "AppId Test",
	// 	Text:      "secured text sending",
	// 	Metadata:  "meta data description sample",
	// 	Operation: models.Create,
	// }
	// if err = srvtext.ProcessingText(ctx, testMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// updateTxtMsg := testMsg
	// updateTxtMsg.Text = "update secured text sending"
	// updateTxtMsg.Metadata = "123"
	// updateTxtMsg.Operation = models.Update
	// if err = srvtext.ProcessingText(ctx, updateTxtMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// resultTextSearch, err := srvtext.SearchText(ctx, "123")
	// if err != nil {
	// 	log.Print("search login_records error :", err)
	// }

	// delTxtMsg := testMsg
	// delTxtMsg.Operation = models.Delete
	// if err = srvtext.ProcessingText(ctx, delTxtMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// log.Print(resultTextSearch)

	// srvlogin := services.NewLogin(sl, init.cfg)

	// msgLogin := models.LoginRecord{
	// 	RecordID:  uuid.NewString(),
	// 	ChngTime:  time.Now(),
	// 	UID:       "uid1",
	// 	AppID:     "app1",
	// 	Login:     "login0001",
	// 	Psw:       "password001",
	// 	Metadata:  "meta data description sample",
	// 	Operation: models.Create,
	// }
	// if err = srvlogin.ProcessingLogin(ctx, msgLogin); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// updateLoginMsg := msgLogin
	// updateLoginMsg.Login = "update login"
	// updateLoginMsg.Metadata = "123update"
	// updateLoginMsg.Operation = models.Update
	// if err = srvlogin.ProcessingLogin(ctx, updateLoginMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// resultSearch, err := srvlogin.SearchLogin(ctx, "12")
	// if err != nil {
	// 	log.Print("search login_records error :", err)
	// }
	// delLoginMsg := updateLoginMsg
	// delLoginMsg.Operation = models.Delete
	// if err = srvlogin.ProcessingLogin(ctx, delLoginMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// log.Print(resultSearch)

	// srvBinary := services.NewBinary(sl, init.cfg)

	// msgBinary := models.BinaryRecord{
	// 	RecordID:  uuid.NewString(),
	// 	ChngTime:  time.Now(),
	// 	UID:       "newUserCfg.UserID",
	// 	AppID:     "newAppCfg.Appid",
	// 	Binary:    "secured binary sending",
	// 	Metadata:  "meta data description sample",
	// 	Operation: models.Create,
	// }

	// if err = srvBinary.ProcessingBinary(ctx, msgBinary); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// updateBinaryMsg := msgBinary
	// updateBinaryMsg.Binary = "update binary"
	// updateBinaryMsg.Metadata = "123updateBinary"
	// updateBinaryMsg.Operation = models.Update
	// if err = srvBinary.ProcessingBinary(ctx, updateBinaryMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// delBinaryMsg := updateBinaryMsg
	// delBinaryMsg.Operation = models.Delete
	// if err = srvBinary.ProcessingBinary(ctx, delBinaryMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// binaryResultSearch, err := srvBinary.SearchBinary(ctx, "sample")
	// if err != nil {
	// 	log.Print("search login_records error :", err)
	// }

	// log.Print(binaryResultSearch)

	// srcvCard := services.NewCard(sl, init.cfg)

	// msgCard := models.CardRecord{
	// 	RecordID:  uuid.NewString(),
	// 	ChngTime:  time.Now(),
	// 	UID:       "CardUID",
	// 	AppID:     "AppIDCard",
	// 	Brand:     "1",
	// 	ValidDate: "01/28",
	// 	Number:    "2202245445789856",
	// 	Code:      "123",
	// 	Holder:    "DMTIRY BO",
	// 	Metadata:  "meta data card description sample",
	// 	Operation: models.Create,
	// }

	// if err = srcvCard.ProcessingCard(ctx, msgCard); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// updateCardMsg := msgCard
	// updateCardMsg.ValidDate = "01/30"
	// updateCardMsg.Metadata = "card123"
	// updateCardMsg.Number = "2202245400000000"
	// updateCardMsg.Operation = models.Update
	// if err = srcvCard.ProcessingCard(ctx, updateCardMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// cardSearchResult, err := srcvCard.SearchCard(ctx, "card")
	// if err != nil {
	// 	log.Print("search login_records error :", err)
	// }

	// delCardMsg := updateCardMsg
	// delCardMsg.Operation = models.Delete
	// if err = srcvCard.ProcessingCard(ctx, delCardMsg); err != nil {
	// 	log.Print("storage new error:", err)
	// }

	// log.Print(cardSearchResult)
}
