package clientrmq

import (
	"fmt"
	"log"

	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/streadway/amqp"
)

// Services структура конструктора бизнес логики.
type ClientRMQ struct {
	Cfg  models.RabbitmqSrv
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// New.
func NewClientRMQ(cfg models.RabbitmqSrv) (*ClientRMQ, error) {
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

func (r *ClientRMQ) QueueDeclare(queueName string) (amqp.Queue, error) {
	q, err := r.Ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
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
