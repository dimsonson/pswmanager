package clientrmq

import (
	"fmt"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Services структура конструктора бизнес логики.
type ClientRMQ struct {
	Cfg  config.RabbitmqSrv
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// New.
func NewClientRMQ(cfg config.RabbitmqSrv) (*ClientRMQ, error) {
	var err error
	rabbitConnURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Psw,
		cfg.Host,
		cfg.Port,
	)
	conn, err := amqp.Dial(rabbitConnURL)
	if err != nil {
		log.Print("rabbitmq connection error: ", err)
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Print("rabbitmq chanel creation error: ", err)
		return nil, err
	}
	return &ClientRMQ{
		Cfg:  cfg,
		Conn: conn,
		Ch:   ch,
	}, err
}

func (r *ClientRMQ) Close() {
	r.Ch.Close()
	r.Conn.Close()
}

func (r *ClientRMQ) ExchangeDeclare(exchName string) error {
	err := r.Ch.ExchangeDeclare(
		exchName,            // name
		r.Cfg.Exchange.Kind, // kind
		true,                // durable
		false,               // delete when unused
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	return err
}

func (r *ClientRMQ) QueueDeclare(queueName string) (models.Queue, error) {
	amqpq, err := r.Ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	q := models.Queue{
		Name:      amqpq.Name,
		Messages:  amqpq.Messages,
		Consumers: amqpq.Consumers,
	}
	return q, err
}

func (r *ClientRMQ) QueueBind(queueName string, routingKey string) error {
	err := r.Ch.QueueBind(
		queueName,           // queue name
		routingKey,          // routing key
		r.Cfg.Exchange.Name, // exchange
		false,               // no-wait
		nil,                 // arguments
	)
	return err
}

func (r *ClientRMQ) UserInit() (config.UserConfig, *config.App) {
	usercfg := config.UserConfig{}
	usercfg.UserID = uuid.New().String()
	userapp := new(config.App)
	userapp.AppID = uuid.New().String()
	userapp.RoutingKey = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumeQueue = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumerName = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ExchangeBindings = []string{}
	usercfg.Apps = append(usercfg.Apps, *userapp)
	return usercfg, userapp
}

func (r *ClientRMQ) AppInit(usercfg config.UserConfig) config.App {
	userapp := config.App{}
	userapp.AppID = uuid.New().String()
	userapp.ConsumeQueue = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumerName = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.RoutingKey = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	return userapp
}
