package models

import (
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/MashinIvan/rabbitmq"
	"github.com/streadway/amqp"
)

type GRPC struct {
	Network string
	Port    string
}

// RabbitmqSrv обобщающая структура конфигурации RabbitMq server.
type RabbitmqSrv struct {
	//Dsn            string `json:"rabbitmq_dsn"`
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

// ExchangeParams generalizes amqp exchange settings
type ExchangeParams struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// QueueParams generalizes amqp queue settings
type QueueParams struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// QualityOfService generalizes amqp qos settings
type QualityOfService struct {
	PrefetchCount int
	PrefetchSize  int
}

// ConsumerParams generalizes amqp consumer settings
type ConsumerParams struct {
	ConsumerName string
	AutoAck      bool
	ConsumerArgs amqp.Table
}

// ControllerParams generalizes rmq server controllers settings
type ControllerParams struct {
	RoutingKey string
	Controller rabbitmq.ControllerFunc
}

type Redis struct {
	Dsn       string `json:"redis_dsn"`
	Username  string
	Password  string
	Network   string 
	Addr      string // host:port address.
	DB        int
	TLSConfig *tls.Config
}

type PostgreSQL struct {
	Dsn  string `json:"postgre_dsn"`
	Conn *sql.DB
}

type ClientRMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// MsgType тип исплльзуемый для проставления признака типа сообщения / операции.
type MsgType int

// Константы типа MsgType исплльзуемые для проставления признака типа сообщения / операции.
const (
	Create MsgType = iota + 1
	Read
	Update
	Delete
)

// CardType тип исплльзуемый для проставления признака типа банковской карты.
type CardType int

// Константы типа CardType исплльзуемые для проставления признака типа банковской карты.
const (
	Mir CardType = iota + 1
	MasterCard
	Visa
	AmEx
)

// SetRecords .
type SetRecords struct {
	SetLoginRec  []LoginRec
	SetTextRec   []TextRec
	SetBinaryRec []BinaryRec
	SetCardRec   []CardRec
}

// LoginRec структура сообщния для опараций с парами логин/пароль.
type LoginRec struct {
	RecordID  string
	ChngTime  time.Time
	UID       string
	AppID     string
	Login     string
	Psw       string
	Metadata  string
	Operation MsgType
}

// LoginRec структура сообщния для опараций с текстовыми данными пользователя.
type TextRec struct {
	RecordID  string
	ChngTime  time.Time
	UID       string
	AppID     string
	Text      string
	Metadata  string
	Operation MsgType
}

// BinaryRec структура сообщния для опараций с бинарными данными пользователя.
type BinaryRec struct {
	RecordID  string
	ChngTime  time.Time
	UID       string
	AppID     string
	Binary    string
	Metadata  string
	Operation MsgType
}

// CardRec структура сообщния для опараций с данными карт пользователя.
type CardRec struct {
	RecordID  string
	ChngTime  time.Time
	UID       string
	AppID     string
	Brand     CardType
	Number    string
	ValidDate string
	Code      int
	Holder    string
	Metadata  string
	Operation MsgType
}

// UserConfig
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
