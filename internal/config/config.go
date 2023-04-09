package config

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/dimsonson/pswmanager/internal/clientrmq"
	"github.com/dimsonson/pswmanager/internal/handlers/rmqhandlers"
	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/dimsonson/pswmanager/internal/router"
	"github.com/dimsonson/pswmanager/internal/servers/grpc"
	"github.com/dimsonson/pswmanager/internal/servers/rabbitmq"
	"github.com/dimsonson/pswmanager/internal/services"
	"github.com/dimsonson/pswmanager/internal/storage/nosql"
	"github.com/dimsonson/pswmanager/internal/storage/sql"
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
	ServerAddress   string     `json:"server_address"`
	FileStoragePath string     `json:"file_storage_path"`
	EnableTLS       bool       `json:"enable_https"`
	TrustedSubnet   string     `json:"trusted_subnet"`
	TrustedCIDR     *net.IPNet `json:"-"`
	ConfigJSONpath  string     `json:"-"`
	Rabbitmq        models.RabbitmqSrv
	Redis           models.Redis
	Postgree        models.PostgreSQL
	GRPC            models.GRPC
	Wg              sync.WaitGroup
}

// NewConfig конструктор создания конфигурации сервера из переменных оружения, флагов, конфиг файла, а так же значений по умолчанию.
func New() *ServiceConfig {
	return &ServiceConfig{}
}

// Parse метод парсинга и получения значений из переменных оружения, флагов, конфиг файла, а так же значений по умолчанию.
func (cfg *ServiceConfig) Parse() {

	cfg.Postgree.Dsn = "postgres://postgres:1818@localhost:5432/pswm"

	cfg.Rabbitmq.RoutingWorkers = 8
	cfg.Rabbitmq.Controllers = make([]models.ControllerParams, 10)
	cfg.Rabbitmq.Controllers[0].RoutingKey = "all.text.*"
	cfg.Rabbitmq.Controllers[1].RoutingKey = "all.login.*"
	cfg.Rabbitmq.Controllers[2].RoutingKey = "all.binary.*"
	cfg.Rabbitmq.Controllers[3].RoutingKey = "all.card.*"

	cfg.Rabbitmq.User = "rmuser"
	cfg.Rabbitmq.Psw = "rmpassword"

	cfg.Rabbitmq.Exchange.Kind = "topic"
	cfg.Rabbitmq.Exchange.Name = "records"
	cfg.Rabbitmq.Exchange.AutoDelete = false
	cfg.Rabbitmq.Exchange.Durable = true

	cfg.Rabbitmq.Queue.Name = "master"
	cfg.Rabbitmq.Queue.Durable = true
	cfg.Rabbitmq.Queue.AutoDelete = true

	cfg.Rabbitmq.Consumer.ConsumerName = "master"
	// cfg.Rabbitmq.Consumer.AutoAck = true
	cfg.Rabbitmq.Consumer.ConsumerArgs = nil

	//cfg.Rabbitmq.

	cfg.GRPC.Network = "tcp"
	cfg.GRPC.Port = ":8080"

	//cfg.RedisDsn = "redis"
	//cfg.PostgreDsn = "redis"
	// описываем флаги
	//addrFlag := flag.String("a", "", "master server address")
	// baseFlag := flag.String("b", "", "dase URL")
	// //pathFlag := flag.String("f", "", "File storage path")
	// //dlinkFlag := flag.String("d", "", "database DSN link")
	// tlsFlag := flag.Bool("s", false, "run as HTTPS server")
	// cfgFlag := flag.String("c", "", "config json path")
	// trustFlag := flag.String("t", "", "trusted subnet CIDR for /api/internal/stats")
	// grpcFlag := flag.Bool("g", false, "run as GRPC server")
	// // парсим флаги в переменные
	// flag.Parse()
	// var ok bool
	// // используем структуру cfg models.Config для хранения параментров необходимых для запуска сервера
	// // проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	// cfg.ConfigJSONpath, ok = os.LookupEnv("CONFIG")
	// if !ok && *cfgFlag != "" {
	// 	log.Print("eviroment variable CONFIG is empty or has wrong value ", cfg.ConfigJSONpath)
	// 	cfg.ConfigJSONpath = *cfgFlag
	// }
	// // читаем конфигурвационный файл и парксим в стркутуру
	// if cfg.ConfigJSONpath != "" {
	// 	configFile, err := os.ReadFile(*cfgFlag)
	// 	if err != nil {
	// 		log.Print("reading config file error:", err)
	// 	}
	// 	if err == nil {
	// 		err = json.Unmarshal(configFile, &cfg)
	// 		if err != nil {
	// 			log.Printf("unmarshal config file error: %s", err)
	// 		}
	// 	}
	// }
	// // проверяем наличие флага или пременной окружения для CIDR доверенной сети эндпойнта /api/internal/stats
	// TrustedSubnet, ok := os.LookupEnv("TRUSTED_SUBNET")
	// if ok {
	// 	cfg.TrustedSubnet = TrustedSubnet
	// }
	// if *trustFlag != "" {
	// 	cfg.TrustedSubnet = *trustFlag
	// }
	// if cfg.TrustedSubnet != "" {
	// 	var err error
	// 	_, cfg.TrustedCIDR, err = net.ParseCIDR(cfg.TrustedSubnet)
	// 	if err != nil {
	// 		log.Print("parse CIDR error: ", err)
	// 	}
	// }
	// // проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	// ServerAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	// if ok {
	// 	cfg.ServerAddress = ServerAddress
	// }
	// //if (!ok || !govalidator.IsURL(cfg.ServerAddress) || cfg.ServerAddress == "") && *addrFlag != "" {
	// //	log.Print("eviroment variable SERVER_ADDRESS is empty or has wrong value ")
	// //cfg.ServerAddress = *addrFlag
	// //}
	// // если нет флага или переменной окружения используем переменную по умолчанию
	// //if !ok && *addrFlag == "" {
	// cfg.ServerAddress = settings.DefServAddr
	// //}
	// // проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	// //BaseURL, ok := os.LookupEnv("BASE_URL")
	// if ok {
	// 	//cfg.BaseURL = BaseURL
	// }
	// //if (!ok || !govalidator.IsURL(cfg.BaseURL) || cfg.BaseURL == "") && *baseFlag != "" {
	// log.Print("eviroment variable BASE_URL is empty or has wrong value ")
	// //	cfg.BaseURL = *baseFlag
	// //}
	// // если нет флага или переменной окружения используем переменную по умолчанию
	// if !ok && *baseFlag == "" {
	// 	//	cfg.BaseURL = settings.DefBaseURL
	// }
	// // проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	// // DatabaseDsn, ok := os.LookupEnv("DATABASE_DSN")
	// // if ok {
	// // 	cfg.DatabaseDsn = DatabaseDsn
	// // }
	// // if !ok && *dlinkFlag != "" {
	// // 	log.Print("eviroment variable DATABASE_DSN is not exist")
	// // 	cfg.DatabaseDsn = *dlinkFlag
	// // }
	// // проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	// FileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	// if ok {
	// 	cfg.FileStoragePath = FileStoragePath
	// }
	// if !ok || (cfg.FileStoragePath == "" || !govalidator.IsUnixFilePath(cfg.FileStoragePath) || govalidator.IsWinFilePath(cfg.FileStoragePath)) {
	// 	log.Print("eviroment variable FILE_STORAGE_PATH is empty or has wrong value ")
	// 	//cfg.FileStoragePath = *pathFlag
	// }
	// // проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	// EnableGRPC, ok := os.LookupEnv("ENABLE_GRPC")
	// if ok && EnableGRPC == "true" || *grpcFlag {
	// 	//cfg.EnableGRPC = true
	// }
	// if !ok {
	// 	log.Print("eviroment variable ENABLE_GRPC is empty or has wrong value ")
	// }
	// // проверяем наличие флага или пременной окружения для старта в https (tls)
	// EnableHTTPS, ok := os.LookupEnv("ENABLE_HTTPS")
	// if ok && EnableHTTPS == "true" || *tlsFlag {
	// 	cfg.EnableTLS = true
	// 	return
	// }
	// // если нет флага или переменной окружения используем переменную по умолчанию
	// cfg.EnableTLS = settings.DefHTTPS
	// log.Print("eviroment variable ENABLE_HTTPS is empty or has wrong value ")
}

func (cfg *ServiceConfig) ServerStart(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup) {

	noSQLstorage := nosql.New("test")

	clientRMQ := clientrmq.NewClientRMQ(cfg.Rabbitmq)
	cfg.Rabbitmq.ClientRMQ.Ch = clientRMQ.Ch
	cfg.Rabbitmq.ClientRMQ.Conn = clientRMQ.Conn 

	servUserCreate := services.NewUserData(noSQLstorage, clientRMQ, cfg.Rabbitmq)

	ucfg, _ := servUserCreate.CreateUser(ctx, "testlogin", "passwtest")

	fmt.Println(servUserCreate.CreateApp(ctx, ucfg.UserID, "passwtest"))

	SQLstorage := sql.New(cfg.Postgree.Dsn)
	cfg.Postgree.Conn = SQLstorage.PostgreConn

	servLoginRec := services.NewLoginRec(SQLstorage)
	servTextRec := services.NewTextRec(SQLstorage)
	servCardRec := services.NewCardRec(SQLstorage)
	servBinaryRec := services.NewBinaryRec(SQLstorage)

	handlers := rmqhandlers.New(servTextRec, servLoginRec, servBinaryRec, servCardRec)

	grpcSrv := grpc.NewServer(ctx, stop, cfg.GRPC, wg)
	grpcSrv.InitGRPCservice()
	wg.Add(1)
	grpcSrv.GrpcGracefullShotdown()
	wg.Add(1)
	go grpcSrv.StartGRPC()

	rmqRouter := router.New(ctx, cfg.Rabbitmq, *handlers)
	rmqSrv := rabbitmq.NewServer(ctx, stop, cfg.Rabbitmq, wg)
	rmqSrv.Init()
	wg.Add(1)
	rmqSrv.Shutdown()
	wg.Add(1)
	rmqSrv.Start(ctx, rmqRouter)

}

func (cfg *ServiceConfig) ConnClose(ctx context.Context) {
	cfg.Postgree.Conn.Close()
	cfg.Rabbitmq.ClientRMQ.Ch.Close()
	cfg.Rabbitmq.ClientRMQ.Conn.Close()
}
