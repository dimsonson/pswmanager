package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
)

type UserStorageProviver interface {
	Close()
	CreateUser(ctx context.Context, login string, psw string, uid string, usercfg []byte) error
	ReadUserCfg(ctx context.Context, uid string) ([]byte, error)
	UpdateUser(ctx context.Context, uid string, bytesUserCfg []byte) error
	CheckPsw(ctx context.Context, uid string, psw string) (bool, error)
	IsUserLoginExist(ctx context.Context, login string) (bool, error)
}

type ClientRMQProvider interface {
	Close()
	ExchangeDeclare(exchName string) error
	QueueDeclare(queueName string) (amqp.Queue, error)
	QueueBind(queueName string, routingKey string) error
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	storage   UserStorageProviver
	clientRMQ ClientRMQProvider
	Cfg       models.RabbitmqSrv
}

// New.
func NewUserData(s UserStorageProviver, clientrmq ClientRMQProvider, cfg models.RabbitmqSrv) *UserServices {
	return &UserServices{
		storage:   s,
		clientRMQ: clientrmq,
		Cfg:       cfg,
	}
}

// CreateUser.
func (s *UserServices) CreateUser(ctx context.Context, login string, psw string) (models.UserConfig, error) {
	// проверка существования пользователя
	ok, err := s.storage.IsUserLoginExist(ctx, login)
	if ok {
		log.Printf("login \"%s\" already exist:", login)
		return models.UserConfig{}, err
	}
	if err != nil {
		log.Print("check login error: ", err)
		return models.UserConfig{}, err
	}
	// создание uidcfg
	usercfg := models.UserConfig{}
	usercfg.UserID = uuid.New().String()
	// создание конфигурации rmq
	userapp := new(models.App)
	userapp.AppID = uuid.New().String()
	userapp.RoutingKey = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumeQueue = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumerName = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ExchangeBindings = []string{} // одно приложение - нет rk

	usercfg.Apps = append(usercfg.Apps, *userapp)
	// bindings отсутвует, т.к. одно приложение
	_, err = s.clientRMQ.QueueDeclare(
		userapp.ConsumeQueue, // name
	)
	if err != nil {
		log.Print("rabbitmq queue creation error: ", err)
		return models.UserConfig{}, err
	}
	// сериализация для хранения в Redis
	bytesUserCfg, err := json.Marshal(usercfg)
	if err != nil {
		log.Print("usercfg encoding error: ", err)
		return models.UserConfig{}, err
	}
	// сохраняем в хранилище
	err = s.storage.CreateUser(ctx, login, psw, usercfg.UserID, bytesUserCfg)
	if err != nil {
		log.Print("usercfg creating in storage error: ", err)
		return models.UserConfig{}, err
	}
	// возращаем приложению клиента
	return usercfg, err
}

// CreateUser получаем psw хешированный base64 и .
func (s *UserServices) CreateApp(ctx context.Context, uid string, psw string) (string, models.UserConfig, error) {
	// проекрка логина и пароля пользователя
	ok, err := s.storage.CheckPsw(ctx, uid, psw)
	if !ok {
		log.Print("uid or psw incorret")
		return "", models.UserConfig{}, errors.New("uid or psw incorret")
	}
	if err != nil {
		log.Print("check psw error: ", err)
		return "", models.UserConfig{}, err
	}
	// получаем конфигурацию из хранилища
	usercfg := models.UserConfig{}
	bytesUserCfg, err := s.storage.ReadUserCfg(ctx, uid)
	if err != nil {
		log.Print("read usercfg from storage error: ", err)
		return "", models.UserConfig{}, err
	}
	err = json.Unmarshal(bytesUserCfg, &usercfg)
	if err != nil {
		log.Print("usercfg decoding error: ", err)
		return "", models.UserConfig{}, err
	}
	// генерируем AppID
	userapp := models.App{}
	userapp.AppID = uuid.New().String()
	// генерируем очередь и routingkey
	userapp.ConsumeQueue = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumerName = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.RoutingKey = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	// добавляем очередь для нового приложения пользователя
	q, err := s.clientRMQ.QueueDeclare(userapp.ConsumeQueue)
	if err != nil {
		log.Print("rabbitmq queue creation error: ", err)
		return "", models.UserConfig{}, err
	}
	err = s.clientRMQ.QueueBind(q.Name, userapp.RoutingKey)
	if err != nil {
		log.Print("rabbitmq queue bindings error: ", err)
		return "", models.UserConfig{}, err
	}
	// добавляем routingkey в биндинги всех приложений клиента
	// добавление routingkey в очереди rabbit (bindings)
	for _, v := range usercfg.Apps {
		v.ExchangeBindings = append(v.ExchangeBindings, userapp.RoutingKey)
		err = s.clientRMQ.QueueBind(
			v.ConsumeQueue,     // queue name
			userapp.RoutingKey, // routing key
		)
		if err != nil {
			log.Print("rabbitmq queue bindings error: ", err)
			return "", models.UserConfig{}, err
		}
	}
	// добавлем новое приложение в структуру конфигурации
	usercfg.Apps = append(usercfg.Apps, userapp)
	// сохраняем обновленную конфигурацию в хранилище
	// сериализация для хранения в Redis
	bytesUserCfg, err = json.Marshal(usercfg)
	if err != nil {
		log.Print("usercfg encoding error: ", err)
		return "", models.UserConfig{}, err
	}
	// сохраняем в хранилище
	s.storage.UpdateUser(ctx, uid, bytesUserCfg)
	// возвращаем AppID и конфигурацию приложению клиента
	return userapp.AppID, usercfg, err
}
