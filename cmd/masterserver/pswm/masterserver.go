package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dimsonson/pswmanager/internal/masterserver/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/initstart"

	"github.com/dimsonson/pswmanager/pkg/log"
)

func init() {
	//log.LogInit()
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
	// остановка всех сущностей, куда передан контекст по прерыванию
	stop()
	// закрываем соединения
	cfg.ConnClose(ctx)
	// ожидаем выполнение горутин
	wg.Wait()
	// логирование закрытия сервера без ошибок
	log.Print("rmq & grpc servers gracefully shutdown")
}
