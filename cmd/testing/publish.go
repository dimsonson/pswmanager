package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MashinIvan/rabbitmq"
	"github.com/MashinIvan/rabbitmq/pkg/backoff"
	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/streadway/amqp"
)

// var exchangeParams = rabbitmq.ExchangeParams{
// 	Name:       "records",
// 	Kind:       "topic",
// 	AutoDelete: false,
// 	Durable:    true,
// }

// var queueParams = rabbitmq.QueueParams{
// 	Name:       "master",
// 	AutoDelete: false,
// 	Durable:    true,
// }

// var qos = rabbitmq.QualityOfService{
// 	PrefetchCount: 5,
// }

// var consumer = rabbitmq.ConsumerParams{
// 	ConsumerName: "master",
// 	AutoAck:      true,
// 	ConsumerArgs: nil,
// }

// func startServerAsync(server *rabbitmq.Server) error {
// 	errChan := make(chan error)
// 	go func() {
// 		err := server.ListenAndServe(context.Background())
// 		if err != nil {
// 			errChan <- err
// 		}
// 	}()

// 	// wait for server startup
// 	time.Sleep(1 * time.Second)

// 	// check for ListenAndServe error
// 	select {
// 	case err := <-errChan:
// 		return err
// 	default:
// 		return nil
// 	}
// }

// const (
// 	user string = "rmuser"
// 	psw  string = "rmpassword"
// 	host string = "localhost"
// )

func connFactory0() (*amqp.Connection, error) {
	connUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		user,
		psw,
		host,
		"",
	)

	return amqp.Dial(connUrl)
}

func main0() {
	conn, err := rabbitmq.NewConnection(connFactory, backoff.NewDefaultSigmoidBackoff())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	msgText := models.TextRec{
		RecordID:  "12345",
		ChngTime:  time.Now(),
		UID:       "user12345",
		AppID:     "app 12345",
		Text:      "secured text",
		Metadata:  "meta data description",
		Operation: models.Create,
	}

	msgBinary := models.BinaryRec{
		RecordID:  "12345",
		ChngTime:  time.Now(),
		UID:       "user12345",
		AppID:     "app 12345",
		Binary:    "secured binary",
		Metadata:  "meta data description",
		Operation: models.Create,
	}

	msgLogin := models.LoginRec{
		RecordID:  "12345",
		ChngTime:  time.Now(),
		UID:       "user12345",
		AppID:     "app 12345",
		Login:     "userlog@ya.ru",
		Psw:       "pasword123",
		Metadata:  "meta data description",
		Operation: models.Create,
	}

	msgCard := models.CardRec{
		RecordID:  "12345",
		ChngTime:  time.Now(),
		UID:       "user12345",
		AppID:     "app 12345",
		Brand:     models.MasterCard,
		ValidDate: "01/03",
		Code:      234,
		Number:    "2202443565854455",
		Metadata:  "meta data description",
		Operation: models.Create,
	}

	publisherCh, err := conn.Channel()
	if err != nil {
		log.Print(err)
	}

	msgText.RecordID = "3984c07c-4d41-493b-882a-e59fb1831606"   //uuid.New().String()
	msgBinary.RecordID = "3984c07c-4d41-493b-882a-e59fb1831603" //uuid.New().String()
	msgLogin.RecordID = "3984c07c-4d41-493b-882a-e59fb1831604"  //uuid.New().String()
	msgCard.RecordID = "3984c07c-4d41-493b-882a-e59fb1831605"   //uuid.New().String()

	//msgText
	//Create

	msgTextJSON, err := json.Marshal(msgText)
	if err != nil {
		log.Print("marshall error", err)
	}
	fmt.Println(string(msgTextJSON))
	err = publisherCh.Publish(exchangeParams.Name, "all.text", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgTextJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	//Update
	msgText.Operation = models.Update
	msgText.Metadata = fmt.Sprintf("test meta data : %v", 8)
	msgTextJSON, err = json.Marshal(msgText)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.text", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgTextJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	//Delete
	msgText.Operation = models.Delete
	msgTextJSON, err = json.Marshal(msgText)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.text", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgTextJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}

	//msgBinary
	//Create
	msgBinaryJSON, err := json.Marshal(msgBinary)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.binary", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgBinaryJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	//Update
	msgBinary.Operation = models.Update
	msgBinary.Metadata = "test meta data : update"
	msgBinary.Binary = "binary data"
	msgBinaryJSON, err = json.Marshal(msgBinary)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.binary", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgBinaryJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	//Delete
	msgBinary.Operation = models.Delete
	msgBinaryJSON, err = json.Marshal(msgBinary)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.binary", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgBinaryJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}

	//msgLogin
	//Create
	msgLoginJSON, err := json.Marshal(msgLogin)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.login", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgLoginJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	//Update
	msgLogin.Operation = models.Update
	msgLogin.Metadata = fmt.Sprintf("test meta data : %v", 8)
	msgLoginJSON, err = json.Marshal(msgLogin)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.login", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgLoginJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	//Delete
	msgLogin.Operation = models.Delete
	msgLoginJSON, err = json.Marshal(msgLogin)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.login", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgLoginJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}

	//msgCard
	//Create
	msgCardJSON, err := json.Marshal(msgCard)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.card", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgCardJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	//Update
	msgCard.Operation = models.Update
	msgCard.Metadata = fmt.Sprintf("test meta data : %v", 8)
	msgCardJSON, err = json.Marshal(msgCard)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.login", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgCardJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	//Delete
	msgCard.Operation = models.Delete
	msgCardJSON, err = json.Marshal(msgCard)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, "all.card", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgCardJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}

	// for i := 0; i < 10; i++ {

	//msg.Metadata = fmt.Sprintf("test meta data : %v", i)
	// 	updateJSON, err := json.Marshal(msg)
	// 	if err != nil {
	// 		log.Print("marshall error", err)
	// 	}
	// 	err = publisherCh.Publish(exchangeParams.Name, "all.text", false, false, amqp.Publishing{
	// 		ContentType:  "application/json",
	// 		Body:         updateJSON,
	// 		DeliveryMode: 1,
	// 	})
	// 	if err != nil {
	// 		log.Print(err)
	// 	}
	// }

	// for i := 0; i < 10; i++ {
	// 	err = publisherCh.Publish(exchangeParams.Name, "all.update", false, false, amqp.Publishing{
	// 		ContentType:  "text/plain",
	// 		Body:         updateJSON,
	// 		DeliveryMode: 1,
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// }

	// wait for messages to arrive
	time.Sleep(3 * time.Second)
}
