package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MashinIvan/rabbitmq"
	"github.com/MashinIvan/rabbitmq/pkg/backoff"
	pb "github.com/dimsonson/pswmanager/internal/masterserver/handlers/protobuf"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
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

const (
	user     string = "rmuser"
	psw      string = "rmpassword"
	host     string = "localhost"
	grpcAddr string = "localhost:8080"
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

func connect() (*grpc.ClientConn, *rabbitmq.Connection, *amqp.Channel, error) {
	connGRPC, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
	}
	connRMQ, err := rabbitmq.NewConnection(connFactory, backoff.NewDefaultSigmoidBackoff())
	if err != nil {
		log.Print(err)
	}
	publisherCh, err := connRMQ.Channel()
	if err != nil {
		log.Print(err)
	}
	return connGRPC, connRMQ, publisherCh, err
}

func NewUserApp(ctx context.Context, c pb.UserServicesClient) (*pb.CreateUserResponse, *pb.CreateAppResponse, error) {
	// получаем переменную интерфейсного типа UsersClient, через которую будем отправлять сообщения
	newuser := &pb.CreateUserRequest{
		Login: uuid.NewString(),
		Psw:   "passw123test",
	}
	newUserCfg, err := c.CreateUser(ctx, newuser)
	if err != nil {
		log.Print("create user error: ", err)
	}
	newapp := &pb.CreateAppRequest{
		Uid: newUserCfg.UserID,
		Psw: "passw123test",
	}
	newAppCfg, err := c.CreateApp(ctx, newapp)
	if err != nil {
		log.Print("create app error: ", err)
	}
	return newUserCfg, newAppCfg, err
}

func NewReadUser(ctx context.Context, c pb.UserServicesClient, newAppCfg *pb.CreateAppResponse) (*pb.ReadUserResponse, error) {
	// переменная запроса всех записей пользователя
	newread := &pb.ReadUserRequest{
		Uid: newAppCfg.UserID,
	}
	// запрос всех записей пользователя
	newRead, err := c.ReadUser(ctx, newread)
	if err != nil {
		log.Print("read records error: ", err)
	}
	return newRead, err
}

func main() {
	ctx := context.Background()
	connGRPC, connRMQ, publisherCh, err := connect()
	if err != nil {
		log.Print(err)
	}
	defer connGRPC.Close()
	defer connRMQ.Close()
	defer publisherCh.Close()
	c := pb.NewUserServicesClient(connGRPC)
	newUserCfg, newAppCfg, err := NewUserApp(ctx, c)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	msgTextTest(ctx, newUserCfg, newAppCfg, publisherCh, c)
	msgBinaryTest(ctx, newUserCfg, newAppCfg, publisherCh, c)
	msgLoginTest(ctx, newUserCfg, newAppCfg, publisherCh, c)
	msgCardTest(ctx, newUserCfg, newAppCfg, publisherCh, c)
}

func msgTextTest(ctx context.Context, newUserCfg *pb.CreateUserResponse, newAppCfg *pb.CreateAppResponse, publisherCh *amqp.Channel, c pb.UserServicesClient) {
	// msgText
	// Create
	msgText := models.TextRec{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       newUserCfg.UserID,
		AppID:     newAppCfg.Appid,
		Text:      "secured text sending",
		Metadata:  "meta data description sample",
		Operation: models.Create,
	}
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

	newReadData, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}

	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if newReadData.SetTextRec[0].UID != msgText.UID ||
		newReadData.SetTextRec[0].AppID != msgText.AppID ||
		newReadData.SetTextRec[0].Metadata != msgText.Metadata ||
		newReadData.SetTextRec[0].RecordID != msgText.RecordID ||
		newReadData.SetTextRec[0].Text != msgText.Text ||
		!newReadData.SetTextRec[0].ChngTime.AsTime().Round(time.Second).
			Equal(msgText.ChngTime.UTC().Round(time.Second)) {
		fmt.Println(">>>>>>>> integration TextCreate_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration TextCreate_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
	// msgText
	// Update
	msgText.Operation = models.Update
	msgText.ChngTime = time.Now()
	msgText.Text = "updated" + msgText.Text
	msgText.Metadata = "updated" + msgText.Metadata
	msgTextUpdateJSON, err := json.Marshal(msgText)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".text", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgTextUpdateJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)

	// запрос всех записей пользователя
	newReadUpdate, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if newReadUpdate.SetTextRec[0].UID != msgText.UID ||
		newReadUpdate.SetTextRec[0].AppID != msgText.AppID ||
		newReadUpdate.SetTextRec[0].Metadata != msgText.Metadata ||
		newReadUpdate.SetTextRec[0].RecordID != msgText.RecordID ||
		newReadUpdate.SetTextRec[0].Text != msgText.Text ||
		!newReadUpdate.SetTextRec[0].ChngTime.AsTime().Round(time.Minute).
			Equal(msgText.ChngTime.UTC().Round(time.Minute)) {
		fmt.Println(">>>>>>>> integration TextUpdate_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration TextUpdate_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
	// msgText
	// Delete
	msgText.Operation = models.Delete
	msgText.ChngTime = time.Now()
	msgTextDeleteJSON, err := json.Marshal(msgText)
	if err != nil {
		log.Print("marshall error", err)

	}
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".text", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgTextDeleteJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)
	// запрос всех записей пользователя
	newReadDelete, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if len(newReadDelete.SetTextRec) > 0 {
		fmt.Println(">>>>>>>> integration TextDelete_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration TextDelete_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
}

func msgBinaryTest(ctx context.Context, newUserCfg *pb.CreateUserResponse, newAppCfg *pb.CreateAppResponse, publisherCh *amqp.Channel, c pb.UserServicesClient) {
	// msgBinary
	// Create
	msgBinary := models.BinaryRec{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       newUserCfg.UserID,
		AppID:     newAppCfg.Appid,
		Binary:    "secured text sending",
		Metadata:  "meta data description sample",
		Operation: models.Create,
	}
	msgBinaryJSON, err := json.Marshal(msgBinary)
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
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".binary", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgBinaryJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// // wait for messages to arrive
	time.Sleep(1 * time.Second)

	newReadData, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}

	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if newReadData.SetBinaryRec[0].UID != msgBinary.UID ||
		newReadData.SetBinaryRec[0].AppID != msgBinary.AppID ||
		newReadData.SetBinaryRec[0].Metadata != msgBinary.Metadata ||
		newReadData.SetBinaryRec[0].RecordID != msgBinary.RecordID ||
		newReadData.SetBinaryRec[0].Binary != msgBinary.Binary ||
		!newReadData.SetBinaryRec[0].ChngTime.AsTime().Round(time.Second).
			Equal(msgBinary.ChngTime.UTC().Round(time.Second)) {
		fmt.Println(">>>>>>>> integration BinaryCreate_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration BinaryCreate_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
	// msgText
	// Update
	msgBinary.Operation = models.Update
	msgBinary.ChngTime = time.Now()
	msgBinary.Binary = "updated" + msgBinary.Binary
	msgBinary.Metadata = "updated" + msgBinary.Metadata
	msgBinaryUpdateJSON, err := json.Marshal(msgBinary)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".binary", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgBinaryUpdateJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)

	// запрос всех записей пользователя
	newReadUpdate, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if newReadUpdate.SetBinaryRec[0].UID != msgBinary.UID ||
		newReadUpdate.SetBinaryRec[0].AppID != msgBinary.AppID ||
		newReadUpdate.SetBinaryRec[0].Metadata != msgBinary.Metadata ||
		newReadUpdate.SetBinaryRec[0].RecordID != msgBinary.RecordID ||
		newReadUpdate.SetBinaryRec[0].Binary != msgBinary.Binary ||
		!newReadUpdate.SetBinaryRec[0].ChngTime.AsTime().Round(time.Minute).
			Equal(msgBinary.ChngTime.UTC().Round(time.Minute)) {
		fmt.Println(">>>>>>>> integration BinaryUpdate_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration BinaryUpdate_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
	// msgText
	// Delete
	msgBinary.Operation = models.Delete
	msgBinary.ChngTime = time.Now()
	msgBinaryDeleteJSON, err := json.Marshal(msgBinary)
	if err != nil {
		log.Print("marshall error", err)

	}
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".binary", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgBinaryDeleteJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)
	// запрос всех записей пользователя
	newReadDelete, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if len(newReadDelete.SetBinaryRec) > 0 {
		fmt.Println(">>>>>>>> integration BinaryDelete_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration BinaryDelete_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
}

func msgLoginTest(ctx context.Context, newUserCfg *pb.CreateUserResponse, newAppCfg *pb.CreateAppResponse, publisherCh *amqp.Channel, c pb.UserServicesClient) {
	// msgLogin
	// Create
	msgLogin := models.LoginRec{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       newUserCfg.UserID,
		AppID:     newAppCfg.Appid,
		Login:     "login0001",
		Psw:       "password001",
		Metadata:  "meta data description sample",
		Operation: models.Create,
	}
	msgLoginJSON, err := json.Marshal(msgLogin)
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
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".login", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgLoginJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// // wait for messages to arrive
	time.Sleep(1 * time.Second)

	newReadData, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}

	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if newReadData.SetLoginRec[0].UID != msgLogin.UID ||
		newReadData.SetLoginRec[0].AppID != msgLogin.AppID ||
		newReadData.SetLoginRec[0].Metadata != msgLogin.Metadata ||
		newReadData.SetLoginRec[0].RecordID != msgLogin.RecordID ||
		newReadData.SetLoginRec[0].Login != msgLogin.Login ||
		newReadData.SetLoginRec[0].Psw != msgLogin.Psw ||
		!newReadData.SetLoginRec[0].ChngTime.AsTime().Round(time.Second).
			Equal(msgLogin.ChngTime.UTC().Round(time.Second)) {
		fmt.Println(">>>>>>>> integration LoginCreate_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration LoginCreate_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
	// msgText
	// Update
	msgLogin.Operation = models.Update
	msgLogin.ChngTime = time.Now()
	msgLogin.Login = "updated" + msgLogin.Login
	msgLogin.Metadata = "updated" + msgLogin.Metadata
	msgLoginUpdateJSON, err := json.Marshal(msgLogin)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".login", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgLoginUpdateJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)

	// запрос всех записей пользователя
	newReadUpdate, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if newReadUpdate.SetLoginRec[0].UID != msgLogin.UID ||
		newReadUpdate.SetLoginRec[0].AppID != msgLogin.AppID ||
		newReadUpdate.SetLoginRec[0].Metadata != msgLogin.Metadata ||
		newReadUpdate.SetLoginRec[0].RecordID != msgLogin.RecordID ||
		newReadUpdate.SetLoginRec[0].Login != msgLogin.Login ||
		newReadData.SetLoginRec[0].Psw != msgLogin.Psw ||
		!newReadUpdate.SetLoginRec[0].ChngTime.AsTime().Round(time.Minute).
			Equal(msgLogin.ChngTime.UTC().Round(time.Minute)) {
		fmt.Println(">>>>>>>> integration LoginUpdate_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration LoginUpdate_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
	// msgText
	// Delete
	msgLogin.Operation = models.Delete
	msgLogin.ChngTime = time.Now()
	msgLoginDeleteJSON, err := json.Marshal(msgLogin)
	if err != nil {
		log.Print("marshall error", err)

	}
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".login", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgLoginDeleteJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)
	// запрос всех записей пользователя
	newReadDelete, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if len(newReadDelete.SetLoginRec) > 0 {
		fmt.Println(">>>>>>>> integration LoginDelete_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration LoginDelete_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
}

func msgCardTest(ctx context.Context, newUserCfg *pb.CreateUserResponse, newAppCfg *pb.CreateAppResponse, publisherCh *amqp.Channel, c pb.UserServicesClient) {
	// Create
	msgCard := models.CardRec{
		RecordID:  uuid.NewString(),
		ChngTime:  time.Now(),
		UID:       newUserCfg.UserID,
		AppID:     newAppCfg.Appid,
		Brand:     1,
		ValidDate: "01/28",
		Number:    "2202245445789856",
		Code:      123,
		Holder:    "DMTIRY BO",
		Metadata:  "meta data card description sample",
		Operation: models.Create,
	}
	msgCardJSON, err := json.Marshal(msgCard)
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
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".card", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgCardJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)
	newReadData, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if newReadData.SetCardRec[0].UID != msgCard.UID ||
		newReadData.SetCardRec[0].AppID != msgCard.AppID ||
		newReadData.SetCardRec[0].Metadata != msgCard.Metadata ||
		newReadData.SetCardRec[0].RecordID != msgCard.RecordID ||
		newReadData.SetCardRec[0].ValidDate != msgCard.ValidDate ||
		newReadData.SetCardRec[0].Number != msgCard.Number ||
		newReadData.SetCardRec[0].Code != int64(msgCard.Code) ||
		newReadData.SetCardRec[0].Holder != msgCard.Holder ||
		!newReadData.SetCardRec[0].ChngTime.AsTime().Round(time.Second).
			Equal(msgCard.ChngTime.UTC().Round(time.Second)) {
		fmt.Println(">>>>>>>> integration CardCreate_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration CardCreate_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
	// Update
	msgCard.Operation = models.Update
	msgCard.ChngTime = time.Now()
	msgCard.Holder = "updated" + msgCard.Holder
	msgCard.Metadata = "updated" + msgCard.Metadata
	msgCardUpdateJSON, err := json.Marshal(msgCard)
	if err != nil {
		log.Print("marshall error", err)
	}
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".card", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgCardUpdateJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)
	// запрос всех записей пользователя
	newReadUpdate, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if newReadUpdate.SetCardRec[0].UID != msgCard.UID ||
	newReadUpdate.SetCardRec[0].AppID != msgCard.AppID ||
	newReadUpdate.SetCardRec[0].Metadata != msgCard.Metadata ||
	newReadUpdate.SetCardRec[0].RecordID != msgCard.RecordID ||
	newReadUpdate.SetCardRec[0].ValidDate != msgCard.ValidDate ||
	newReadUpdate.SetCardRec[0].Number != msgCard.Number ||
	newReadUpdate.SetCardRec[0].Code != int64(msgCard.Code) ||
	newReadUpdate.SetCardRec[0].Holder != msgCard.Holder ||
		!newReadUpdate.SetCardRec[0].ChngTime.AsTime().Round(time.Minute).
			Equal(msgCard.ChngTime.UTC().Round(time.Minute)) {
		fmt.Println(">>>>>>>> integration CardUpdate_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration CardUpdate_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
	// Delete
	msgCard.Operation = models.Delete
	msgCard.ChngTime = time.Now()
	msgCardDeleteJSON, err := json.Marshal(msgCard)
	if err != nil {
		log.Print("marshall error", err)

	}
	err = publisherCh.Publish(exchangeParams.Name, routingKey+".card", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         msgCardDeleteJSON,
		DeliveryMode: 1,
	})
	if err != nil {
		log.Print(err)
	}
	// wait for messages to arrive
	time.Sleep(1 * time.Second)
	// запрос всех записей пользователя
	newReadDelete, err := NewReadUser(ctx, c, newAppCfg)
	if err != nil {
		log.Print("read records error: ", err)
		return
	}
	// проверка соответствия направленных записей от приложений клиента выгрузке всех записей клиента из мастер базы
	if len(newReadDelete.SetCardRec) > 0 {
		fmt.Println(">>>>>>>> integration CardDelete_Test: ", settings.ColorRed, "ERROR", settings.ColorReset)
	} else {
		fmt.Println(">>>>>>>> integration CardDelete_Test: ", settings.ColorGreen, "OK", settings.ColorReset)
	}
}
