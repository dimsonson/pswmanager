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
func NewClientRMQ(cfg models.RabbitmqSrv) *ClientRMQ {
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
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Print("rabbitmq chanel creation error: ", err)
		return nil
	}
	return &ClientRMQ{
		Cfg:  cfg,
		Conn: conn,
		Ch:   ch,
	}
}


func (r *ClientRMQ) Close() {
	r.Conn.Close()
}
