package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MashinIvan/rabbitmq"
	"github.com/MashinIvan/rabbitmq/pkg/backoff"
	pb "github.com/dimsonson/pswmanager/internal/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/dimsonson/pswmanager/internal/settings"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var exchangeParams = rabbitmq.ExchangeParams{
	Name:       "records",
	Kind:       "topic",
	AutoDelete: false,
	Durable:    true,
}

var queueParams = rabbitmq.QueueParams{
	Name:       "master",
	AutoDelete: false,
	Durable:    true,
}

var qos = rabbitmq.QualityOfService{
	PrefetchCount: 5,
}

var consumer = rabbitmq.ConsumerParams{
	ConsumerName: "master",
	AutoAck:      true,
	ConsumerArgs: nil,
}

func startServerAsync(server *rabbitmq.Server) error {
	errChan := make(chan error)
	go func() {
		err := server.ListenAndServe(context.Background())
		if err != nil {
			errChan <- err
		}
	}()

	// wait for server startup
	time.Sleep(1 * time.Second)

	// check for ListenAndServe error
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

const (
	user string = "rmuser"
	psw  string = "rmpassword"
	host string = "localhost"
)

func connFactory() (*amqp.Connection, error) {
	connUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		user,
		psw,
		host,
		"",
	)

	return amqp.Dial(connUrl)
}

func main() {

	// устанавливаем соединение с сервером
	connGRPC, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer connGRPC.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewUserServicesClient(connGRPC)

	ctx := context.Background()

	newuser := &pb.CreateUserRequest{
		Login: uuid.NewString(),
		Psw:   "passw123test",
	}

	newUserCfg, err := c.CreateUser(ctx, newuser)
	if err != nil {
		log.Print("create user error: ", err)
		return
	}

	newapp := &pb.CreateAppRequest{
		Uid: newUserCfg.UserID,
		Psw: "passw123test",
	}

	newAppCfg, err := c.CreateApp(ctx, newapp)
	if err != nil {
		log.Print("create app error: ", err)
		return
	}

	connRMQ, err := rabbitmq.NewConnection(connFactory, backoff.NewDefaultSigmoidBackoff())
	if err != nil {
		log.Fatal(err)
	}
	defer connRMQ.Close()

	msgText := models.TextRec{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       newUserCfg.UserID,
		AppID:     newAppCfg.Appid,
		Text:      "secured text sending",
		Metadata:  "meta data description sample",
		Operation: models.Create,
	}

	// msgBinary := models.BinaryRec{
	// 	RecordID:  "12345",
	// 	ChngTime:  time.Now(),
	// 	UID:       "user12345",
	// 	AppID:     "app 12345",
	// 	Binary:    "secured binary",
	// 	Metadata:  "meta data description",
	// 	Operation: models.Create,
	// }

	// msgLogin := models.LoginRec{
	// 	RecordID:  "12345",
	// 	ChngTime:  time.Now(),
	// 	UID:       "user12345",
	// 	AppID:     "app 12345",
	// 	Login:     "userlog@ya.ru",
	// 	Psw:       "pasword123",
	// 	Metadata:  "meta data description",
	// 	Operation: models.Create,
	// }

	// msgCard := models.CardRec{
	// 	RecordID:  "12345",
	// 	ChngTime:  time.Now(),
	// 	UID:       "user12345",
	// 	AppID:     "app 12345",
	// 	Brand:     models.MasterCard,
	// 	ValidDate: "01/03",
	// 	Code:      234,
	// 	Number:    "2202443565854455",
	// 	Metadata:  "meta data description",
	// 	Operation: models.Create,
	// }

	publisherCh, err := connRMQ.Channel()
	if err != nil {
		log.Print(err)
	}

	// //msgText
	// //Create

	msgTextJSON, err := json.Marshal(msgText)
	if err != nil {
		log.Print("marshall error", err)
	}

	var routingKey string
	for _, v := range newAppCfg.Apps {
		if v.AppID == newAppCfg.Appid {
			routingKey = v.RoutingKey
			break
		}
	}

	err = publisherCh.Publish(exchangeParams.Name, routingKey+".text", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgTextJSON,
		DeliveryMode: 1,
	})

	if err != nil {
		log.Print(err)
	}

	// // wait for messages to arrive
	time.Sleep(1 * time.Second)

	newread := &pb.ReadUserRequest{
		Uid: newAppCfg.UserID,
	}

	newRead, err := c.ReadUser(ctx, newread)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}

	if newRead.SetTextRec[0].UID != msgText.UID ||
		newRead.SetTextRec[0].AppID != msgText.AppID ||
		newRead.SetTextRec[0].Metadata != msgText.Metadata ||
		newRead.SetTextRec[0].RecordID != msgText.RecordID ||
		newRead.SetTextRec[0].Text != msgText.Text ||
		!newRead.SetTextRec[0].ChngTime.AsTime().Round(time.Second).
			Equal(msgText.ChngTime.UTC().Round(time.Second)) {
		//fmt.Println(newRead.SetTextRec[0].ChngTime.AsTime().Round(time.Second))
		//fmt.Println(msgText.ChngTime.UTC().Round(time.Second))
		fmt.Println(">>>>>>>> integration TextTest", settings.ColorRed, "ERROR", settings.ColorReset)
		return
	}
	fmt.Println(">>>>>>>> integration TextTest", settings.ColorGreen, "OK", settings.ColorReset)
}
