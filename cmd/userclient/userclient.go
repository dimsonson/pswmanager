package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dimsonson/pswmanager/internal/userclient/config"
	"github.com/dimsonson/pswmanager/internal/userclient/ui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	ui := ui.New()
	ui.Init()
	ui.TextConfig()
	ui.ListConfig()
	ui.FlexConfig()
	ui.PagesConfig()
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:          ui.TextView.LogWindow,
		TimeFormat:   "2006/01/02 15:04:05",
		NoColor:      true,
		FormatCaller: func(i interface{}) string { return fmt.Sprintf("%s:", i) },
		PartsOrder:   []string{zerolog.TimestampFieldName, zerolog.LevelFieldName, zerolog.MessageFieldName, zerolog.CallerFieldName}}).
		With().Caller().Logger()
	stdlog.SetFlags(stdlog.Lshortfile)
	stdlog.SetOutput(log.Logger)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	go ui.UIRun()
}

// Глобальные переменные для использования при сборке - go run -ldflags "-X main.buildVersion=v0.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d')' -X main.buildCommit=final"  main.go.
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	// Вывод данных о версии, дате, коммите сборки.
	log.Printf("version=%s, date=%s, commit=%s", buildVersion, buildDate, buildCommit)

	var wg sync.WaitGroup
	// опередяляем контекст уведомления о сигнале прерывания
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// создание конфигурацию сервера
	cfg := config.New()
	// парсинг конфигурации сервера
	cfg.Parse()
	// старт серверов
	cfg.ServerStart(ctx, stop, &wg)
	// остановка всех сущностей, куда передан контекст по прерыванию
	stop()
	// закрываем соединения
	cfg.ConnClose(ctx)
	// ожидаем выполнение горутин
	wg.Wait()
	// логирование закрытия сервера без ошибок
	log.Print("rmq & grpc servers gracefully shutdown")
}
