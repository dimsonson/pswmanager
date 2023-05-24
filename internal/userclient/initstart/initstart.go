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

	srvtext := services.NewText(sl, clientGRPC, init.cfg)
	srvlogin := services.NewLogin(sl, clientGRPC, init.cfg)
	srvbinary := services.NewBinary(sl, clientGRPC, init.cfg)
	srvcard := services.NewCard(sl, clientGRPC, init.cfg)

	srvusers := services.NewUsers(sl, clientGRPC, init.cfg)
	
	log.Print(init.cfg.UserLogin)
	log.Print(init.cfg.UserPsw)

	ui := ui.NewUI(ctx, init.cfg, srvusers, srvtext, srvlogin, srvbinary, srvcard)
	ui.Init(uiLog)
	go ui.UIRun()

	<-ctx.Done()
}

func (init *Init) ConnClose(ctx context.Context) {
	if init.cfg.SQLight.Conn != nil {
		init.cfg.SQLight.Conn.Close()
	}
}
