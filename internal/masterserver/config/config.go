package config

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"sync"

	"github.com/MashinIvan/rabbitmq"
	"github.com/dimsonson/pswmanager/pkg/log"
	"github.com/streadway/amqp"
)

// ServiceConfig структура конфигурации сервиса, при запуске сервиса с флагом -c/config
// и отсутствии иных флагов и переменных окружения заполняется из файла указанного в этом флаге или переменной окружения CONFIG.
type ServiceConfig struct {
	ServerAddress  string         `json:"server_address"`
	EnableTLS      bool           `json:"enable_tls"`
	ConfigJSONpath string         `json:"-"`
	Rabbitmq       RabbitmqSrv    `json:"rabbitmq"`
	Redis          Redis          `json:"redis"`
	Postgree       PostgreSQL     `json:"postgre"`
	GRPC           GRPC           `json:"grpc"`
	Wg             sync.WaitGroup `json:"-"`
}

// GRPC.
type GRPC struct {
	Network string
	Port    string
}

// RabbitmqSrv обобщающая структура конфигурации RabbitMq server.
type RabbitmqSrv struct {
	User           string
	Psw            string
	Host           string
	Port           string
	ClientRMQ      ClientRMQ
	Exchange       ExchangeParams
	Queue          QueueParams
	QoS            QualityOfService
	Consumer       ConsumerParams
	Controllers    []ControllerParams
	RoutingWorkers int
}

// ExchangeParams общие настройки amqp exchange.
type ExchangeParams struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// QueueParams общие настройки amqp queue.
type QueueParams struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// QualityOfService общие настройки amqp qos.
type QualityOfService struct {
	PrefetchCount int
	PrefetchSize  int
}

// ConsumerParams общие настройки amqp consumer.
type ConsumerParams struct {
	ConsumerName string
	AutoAck      bool
	ConsumerArgs amqp.Table
}

// ControllerParams общие настройки rmq server controllers.
type ControllerParams struct {
	RoutingKey string
	Controller rabbitmq.ControllerFunc `json:"-"`
}

type Redis struct {
	Username  string
	Password  string
	Network   string
	Addr      string // host:port address.
	DB        int
	TLSConfig *tls.Config `json:"-"`
}

type PostgreSQL struct {
	Dsn  string  `json:"postgre_dsn"`
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
			err = json.Unmarshal(configFile, &cfg) //Indent(cfg, configFile, "", "  ") //Unmarshal(configFile, &cfg)
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
