package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"sync"

	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/streadway/amqp"
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
	ServerAddress  string         `json:"server_address"`
	EnableTLS      bool           `json:"enable_tls"`
	ConfigJSONpath string         `json:"-"`
	SQLight        SQLight        `json:"postgre"`
	GRPC           GRPC           `json:"grpc"`
	Wg             sync.WaitGroup `json:"-"`
}

// GRPC.
type GRPC struct {
	Network string
	Port    string
}

type SQLight struct {
	Dsn  string  `json:"sqlite_dsn"`
	Conn *sql.DB `json:"-"`
}

type ClientRMQ struct {
	Conn *amqp.Connection `json:"-"`
	Ch   *amqp.Channel    `json:"-"`
}

// UserConfig.
type UserConfig struct {
	UserID       string `redis:"userid"`
	RmqHost      string `redis:"rmqhost"`
	RmqPort      string `redis:"rmqpost"`
	RmqUID       string `redis:"rmquid"`
	RmqPsw       string `redis:"rmqpsw"`
	ExchangeName string `redis:"exchangename"`
	ExchangeKind string `redis:"exchangekind"`
	Apps         []App  `redis:"apps"`
}

// App .
type App struct {
	AppID            string   `redis:"appid"`
	RoutingKey       string   `redis:"routingkey"`
	ConsumeQueue     string   `redis:"consumequeue"`
	ConsumerName     string   `redis:"consumername"`
	ExchangeBindings []string `redis:"bindings"`
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
	//сохранение congig.json
	cfg.SQLight.Dsn = "file:test.s3db" //?_auth&_auth_user=admin&_auth_pass=admin&_auth_crypt=sha1"
	configFile, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Printf("marshal config file error: %s", err)
	}
	err = os.WriteFile("config.json", configFile, 0666)
	if err != nil {
		log.Printf("write config file error: %s", err)
	}
}

// func (cfg *ServiceConfig) ServerStart(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup) {

// 	sl, err := storage.New(cfg.SQLight.Dsn)
// 	if err != nil {
// 		log.Print("storage new error:", err)
// 	}

// 	srvtext := services.NewText(sl)

// 	testMsg := models.TextRecord{
// 		RecordID:  uuid.NewString(),
// 		ChngTime:  time.Now(),
// 		UID:       "uidtest",
// 		AppID:     "AppId Test",
// 		Text:      "secured text sending",
// 		Metadata:  "meta data description sample",
// 		Operation: models.Create,
// 	}

// 	if err = srvtext.ProcessingText(ctx, testMsg); err != nil {
// 		log.Print("storage new error:", err)
// 	}

// 	// noSQLstorage := nosql.New(cfg.Redis)

// 	// clientRMQ, err := clientrmq.NewClientRMQ(cfg.Rabbitmq)
// 	// if err != nil {
// 	// 	log.Printf("new client error: %s", err)
// 	// 	return
// 	// }
// 	// cfg.Rabbitmq.ClientRMQ.Conn = clientRMQ.Conn
// 	// cfg.Rabbitmq.ClientRMQ.Ch = clientRMQ.Ch
// 	// cfgUser := services.NewUserData(noSQLstorage, clientRMQ, cfg.Rabbitmq)

// 	// SQLstorage := sql.New(cfg.Postgree.Dsn)
// 	// cfg.Postgree.Conn = SQLstorage.PostgreConn

// 	// servLogin := services.NewLogin(SQLstorage)
// 	// servText := services.NewText(SQLstorage)
// 	// servCard := services.NewCard(SQLstorage)
// 	// servBinary := services.NewBinary(SQLstorage)

// 	// cfgReadUsers := services.NewReadUser(SQLstorage)

// 	// handlers := rmq.New(servText, servLogin, servBinary, servCard)

// 	// grpcSrv := grpc.NewServer(ctx, stop, cfg.GRPC, wg)
// 	// grpcSrv.InitGRPCservice(cfgReadUsers, cfgUser)
// 	// //wg.Add(1)
// 	// //grpcSrv.GrpcGracefullShotdown()
// 	// //wg.Add(1)
// 	// //go grpcSrv.StartGRPC()

// 	// rmqRouter := router.New(ctx, cfg.Rabbitmq, *handlers)
// 	// rmqSrv := rabbitmq.NewServer(ctx, stop, cfg.Rabbitmq, wg)
// 	// rmqSrv.Init()
// 	// wg.Add(1)
// 	// rmqSrv.Shutdown()
// 	// wg.Add(1)
// 	// log.Print("rmq starting...")
// 	// rmqSrv.Start(ctx, rmqRouter)
// 	<-ctx.Done()
// }

func (cfg *ServiceConfig) ConnClose(ctx context.Context) {
	if cfg.SQLight.Conn != nil {
		cfg.SQLight.Conn.Close()
	}
}
