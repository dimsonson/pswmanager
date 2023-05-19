package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type ServerRMQProviver interface {
	PublishRecord(exchName string, routingKey string, body []byte) error
}

// Services структура конструктора бизнес логики.
type ServerRMQservices struct {
	clientRMQ ClientRMQProviver
}

// New.
func NewConsume(cRMQ ClientRMQProviver) *ClientRMQservices {
	return &ClientRMQservices{
		clientRMQ: cRMQ,
	}
}

// TextRec.
func (sr *ClientRMQservices) ConsumeRecord(ctx context.Context, exchName string, routingKey string, record interface{}) error {
	switch record.(type) {
	case models.LoginRecord, models.TextRecord, models.BinaryRecord, models.CardRecord:
		msgJSON, err := json.Marshal(record)
		if err != nil {
			log.Print("marshall error", err)
			return err
		}
		err = sr.clientRMQ.PublishRecord(exchName, routingKey, msgJSON)
		if err != nil {
			log.Print("create text record error: ", err)
		}
		return err
	default:
		return errors.New("unknown type is not supprted")
	}
}
