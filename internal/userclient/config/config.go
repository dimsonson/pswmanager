package config

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/clientrmq"

	"github.com/dimsonson/pswmanager/internal/masterserver/handlers/rmq"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/masterserver/router"
	"github.com/dimsonson/pswmanager/internal/masterserver/servers/grpc"
	rabbitmq "github.com/dimsonson/pswmanager/internal/masterserver/servers/rmq"
	"github.com/dimsonson/pswmanager/internal/masterserver/services"
	"github.com/dimsonson/pswmanager/internal/masterserver/storage/nosql"
	"github.com/dimsonson/pswmanager/internal/masterserver/storage/sql"
)

// Константы по умолчанию.
const (
	defServAddr = "localhost:8080"
	defDBlink   = "postgres://postgres:1818@localhost:5432/dbo"
	defTLS      = false
)

// ServiceConfig структура конфигурации сервиса, при запуске сервиса с флагом -c/config
// и отсутствии иных флагов и переменных окружения заполняется из файла указанного в этом флаге или переменной окружения CONFIG.
type ServiceConfig struct {
	ServerAddress  string             `json:"server_address"`
	EnableTLS      bool               `json:"enable_tls"`
	ConfigJSONpath string             `json:"-"`
	Rabbitmq       models.RabbitmqSrv `json:"rabbitmq"`
	Redis          models.Redis       `json:"redis"`
	Postgree       models.PostgreSQL  `json:"postgre"`
	GRPC           models.GRPC        `json:"grpc"`
	Wg             sync.WaitGroup     `json:"-"`
}

// NewConfig конструктор создания конфигурации сервера из переменных оружения, флагов, конфиг файла, а так же значений по умолчанию.
func New() *ServiceConfig {
	return &ServiceConfig{}
}

// Parse метод парсинга и получения значений из переменных оружения, флагов, конфиг файла, а так же значений по умолчанию.
func (cfg *ServiceConfig) Parse() {
	// описываем флаги
	cfgFlag := flag.String("c", "", "config json path")
	// парсим флаги в переменные
	flag.Parse()
	cfg.ConfigJSONpath = *cfgFlag
	// используем структуру cfg models.Config для хранения параментров необходимых для запуска сервера
	// читаем конфигурвационный файл и парксим в стркутуру
	if cfg.ConfigJSONpath != "" {
		configFile, err := os.ReadFile(*cfgFlag)
		if err != nil {
			log.Print("reading config file error:", err)
		}
		if err == nil {
			err = json.Unmarshal(configFile, &cfg)
			if err != nil {
				log.Printf("unmarshal config file error: %s", err)
			}
		}
	}
	// сохранение congig.json
	// configFile, err := json.MarshalIndent(cfg, "", "  ")
	// if err != nil {
	// 	log.Printf("marshal config file error: %s", err)
	// }
	// err = os.WriteFile("config.json", configFile, 0666)
	// if err != nil {
	// 	log.Printf("write config file error: %s", err)
	// }
}

func (cfg *ServiceConfig) ServerStart(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup) {

	noSQLstorage := nosql.New(cfg.Redis)

	clientRMQ, err := clientrmq.NewClientRMQ(cfg.Rabbitmq)
	if err != nil {
		log.Printf("new client error: %s", err)
		return
	}
	cfg.Rabbitmq.ClientRMQ.Conn = clientRMQ.Conn
	cfg.Rabbitmq.ClientRMQ.Ch = clientRMQ.Ch
	cfgUser := services.NewUserData(noSQLstorage, clientRMQ, cfg.Rabbitmq)

	SQLstorage := sql.New(cfg.Postgree.Dsn)
	cfg.Postgree.Conn = SQLstorage.PostgreConn

	servLoginRec := services.NewLoginRec(SQLstorage)
	servTextRec := services.NewTextRec(SQLstorage)
	servCardRec := services.NewCardRec(SQLstorage)
	servBinaryRec := services.NewBinaryServices(SQLstorage)

	cfgReadUsers := services.NewReadUser(SQLstorage)

	handlers := rmq.New(servTextRec, servLoginRec, servBinaryRec, servCardRec)

	grpcSrv := grpc.NewServer(ctx, stop, cfg.GRPC, wg)
	grpcSrv.InitGRPCservice(cfgReadUsers, cfgUser)
	//wg.Add(1)
	//grpcSrv.GrpcGracefullShotdown()
	//wg.Add(1)
	//go grpcSrv.StartGRPC()

	rmqRouter := router.New(ctx, cfg.Rabbitmq, *handlers)
	rmqSrv := rabbitmq.NewServer(ctx, stop, cfg.Rabbitmq, wg)
	rmqSrv.Init()
	wg.Add(1)
	rmqSrv.Shutdown()
	wg.Add(1)
	log.Print("rmq starting...")
	rmqSrv.Start(ctx, rmqRouter)
}

func (cfg *ServiceConfig) ConnClose(ctx context.Context) {
	if cfg.Postgree.Conn != nil {
		cfg.Postgree.Conn.Close()
	}
	if cfg.Rabbitmq.ClientRMQ.Conn != nil {
		cfg.Rabbitmq.ClientRMQ.Conn.Close()
	}
	if cfg.Rabbitmq.ClientRMQ.Ch != nil {
		cfg.Rabbitmq.ClientRMQ.Ch.Close()
	}
}
