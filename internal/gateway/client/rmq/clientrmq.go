package clientrmq

import (
	"fmt"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/gateway/config"
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
	 if !conn.IsClosed() {log.Print("rmq connection open")}
	
	return &ClientRMQ{
		Cfg:  cfg,
		Conn: conn,
		Ch:   ch,
	}, err
}

// func (r *ClientRMQ) Close() {
// 	r.Ch.Close()
// 	r.Conn.Close()
// }

func (r *ClientRMQ) PublishRecord(exchName string, routingKey string, body []byte) error {
	log.Print("exchName :", exchName)
	err := r.Ch.Publish(exchName, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: 1,
	})
	return err
}
