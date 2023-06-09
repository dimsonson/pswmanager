package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	stdlog "log"

	"github.com/dimsonson/pswmanager/internal/gateway/config"
	"github.com/dimsonson/pswmanager/internal/gateway/initstart"
	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/rs/zerolog"
)

func init() {
	//log.LogInit()
	log.Logg = log.Output(
		zerolog.ConsoleWriter{
			Out:          os.Stderr,
			NoColor:      true,
			TimeFormat:   "2006/01/02 15:04:05",
			FormatCaller: func(i interface{}) string { return fmt.Sprintf("%s:", i) },
			PartsOrder: []string{
				zerolog.TimestampFieldName,
				zerolog.LevelFieldName,
				zerolog.MessageFieldName,
				zerolog.CallerFieldName}}).
		With().Caller().Logger()
		stdlog.SetFlags(stdlog.Lshortfile)
	//stdlog.SetOutput(Logg)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	//log.Logg = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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

	init := initstart.New(cfg)
	// старт серверов
	init.InitAndStart(ctx, stop, &wg)
	// остановка всех сущностей, куда передана отмена контекста по прерыванию
	stop()
	// закрываем соединения
	init.ConnClose(ctx)
	// ожидаем выполнение горутин
	wg.Wait()
	// логирование закрытия сервера без ошибок
	log.Print("rmq & grpc servers gracefully shutdown")
}
