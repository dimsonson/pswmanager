package initstart

import (
	"context"
	"sync"
	"time"

	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/google/uuid"

	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
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

	<-ctx.Done()

}
