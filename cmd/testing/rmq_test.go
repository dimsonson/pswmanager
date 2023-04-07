package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/MashinIvan/rabbitmq"
	"github.com/MashinIvan/rabbitmq/pkg/backoff"
	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func TestServer_Publish(t *testing.T) {
	conn, err := rabbitmq.NewConnection(connFactory, backoff.NewDefaultSigmoidBackoff())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	create := models.TextRec{
		RecordID:  "12345",
		ChngTime:  time.Now(),
		UID:       "user12345",
		AppID:     "app 12345",
		Text:      "secured text",
		Metadata:  "meta data description",
		Operation: models.Create,
	}
	// update := models.TextRec{
	// 	RecordID:  12345,
	// 	ChngTime:  time.Now(),
	// 	UID:       "user12345",
	// 	AppID:     "app 12345",
	// 	Text:      "secured updated text",
	// 	Metadata:  "meta data description",
	// 	Operation: models.Update,
	// }
	// updateJSON, err := json.Marshal(update)
	// if err != nil {
	// 	log.Print("marshall error", err)
	// }

	publisherCh, err := conn.Channel()
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		create.RecordID = uuid.New().String()
		createJSON, err := json.Marshal(create)
		if err != nil {
			log.Print("marshall error", err)
		}
		err = publisherCh.Publish(exchangeParams.Name, "all.create", false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         createJSON,
			DeliveryMode: 1,
		})
		if err != nil {
			t.Error(err)
		}
	}

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

func TestServer_Consume(t *testing.T) {
	conn, err := rabbitmq.NewConnection(connFactory, backoff.NewDefaultSigmoidBackoff())
	if err != nil {
		log.Fatal(err)
	}
		defer conn.Close()
		router := rabbitmq.NewRouter()

	routerGroup := router.Group(exchangeParams, queueParams, qos, consumer)


	fooController := func(ctx *rabbitmq.DeliveryContext) {
		fmt.Println("test123")
		fmt.Println(string(ctx.Delivery.Body))
		fmt.Println(ctx.Delivery.Exchange)
		fmt.Println(ctx.Delivery.RoutingKey)
		err = ctx.Delivery.Ack(true)
		if err != nil {
			t.Error(err)
		}
	}

	barController := func(ctx *rabbitmq.DeliveryContext) {
		fmt.Println("test12345678")
		fmt.Println(ctx.Delivery.Exchange)
		fmt.Println(ctx.Delivery.RoutingKey)
		err = ctx.Delivery.Ack(true)
		if err != nil {
			t.Error(err)
		}
	}

	routerGroup.Route("all.create", fooController)
	routerGroup.Route("all.update", barController)

	server := rabbitmq.NewServer(conn, router)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Shutdown(context.Background())

	err = startServerAsync(server)
	if err != nil {

		t.Error(err)
	}

	// wait for messages to arrive
	time.Sleep(15 * time.Second)

}

func TestServer_Shutdown(t *testing.T) {
	conn, err := rabbitmq.NewConnection(connFactory, backoff.NewDefaultSigmoidBackoff())
	if err != nil {
		log.Fatal(err)
	}

	// router := rabbitmq.NewRouter()
	// routerGroup := router.Group(exchangeParams, queueParams, qos, consumer)

	// gracefulStopCompleted := false
	// routerGroup.Route("logs", func(ctx *rabbitmq.DeliveryContext) {
	// 	// simulate message processing time
	// 	fmt.Println("test")
	// 	time.Sleep(5 * time.Second)
	// 	gracefulStopCompleted = true
	// })

	// server := rabbitmq.NewServer(conn, router)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// err = startServerAsync(server)
	// if err != nil {
	// 	t.Error(err)
	// }

	publisherCh, err := conn.Channel()
	if err != nil {
		t.Fatal(err)
	}

	err = publisherCh.Publish(exchangeParams.Name, "logs", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("message"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// wait for publish
	time.Sleep(1 * time.Second)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// 	err = server.Shutdown(ctx)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	//	if !gracefulStopCompleted {
	//		t.Error("graceful stop did not complete")
	//	}
}

func TestServer_Reconnect(t *testing.T) {
	conn, err := rabbitmq.NewConnection(connFactory, backoff.NewDefaultSigmoidBackoff())
	if err != nil {
		log.Fatal(err)
	}

	router := rabbitmq.NewRouter()
	routerGroup := router.Group(
		rabbitmq.ExchangeParams{
			Name:       "test-durable",
			Kind:       "topic",
			Durable:    true,
			AutoDelete: false,
		}, rabbitmq.QueueParams{
			Name:       "test-durable",
			Durable:    true,
			AutoDelete: false,
		}, qos, consumer,
	)

	messageReceivedFirstTime := false
	messageReceivedSecondTime := false

	routerGroup.Route("test.foo", func(ctx *rabbitmq.DeliveryContext) {
		// simulate message processing time
		time.Sleep(1 * time.Second)

		if !messageReceivedFirstTime {
			messageReceivedFirstTime = true
			// test send after server is restarted on message processing
			go func() {
				//err = restartRabbitAndWaitReconnect(7)
				if err != nil {
					t.Error(err)
				}
			}()

			// wait for server stop
			time.Sleep(1 * time.Second)
		}

		if !ctx.Ack() {
			return
		}

		messageReceivedSecondTime = true
	})

	server := rabbitmq.NewServer(conn, router)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Shutdown(context.Background())

	// err = startServerAsync(server)
	// if err != nil {
	// 	t.Error(err)
	// }

	publisher, err := conn.Channel()
	if err != nil {
		t.Fatal(err)
	}

	err = publisher.Publish("test-durable", "test.foo", false, false, amqp.Publishing{
		DeliveryMode: 2,
		ContentType:  "text/plain",
		Body:         []byte("message"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// wait for message to be processed with server shutdown while processing, restarted and processed again
	<-time.After(15 * time.Second)

	err = server.Shutdown(context.Background())
	if err != nil {
		t.Error(err)
	}

	if !messageReceivedFirstTime {
		t.Error("message was not received first time")
	}
	if !messageReceivedSecondTime {
		t.Error("message was not received second time")
	}
}
