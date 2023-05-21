package services

import (
	"context"
	"errors"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type UserStorageProviver interface {
	Close()
	CreateUser(ctx context.Context, login string, psw string, uid string, usercfg config.UserConfig) error
	ReadUserCfg(ctx context.Context, uid string) (config.UserConfig, error)
	UpdateUser(ctx context.Context, uid string, usercfg config.UserConfig) error
	CheckPsw(ctx context.Context, uid string, psw string) (bool, error)
	IsUserLoginExist(ctx context.Context, login string) (bool, error)
}

type ClientRMQProvider interface {
	Close()
	ExchangeDeclare(exchName string) error
	QueueDeclare(queueName string) (models.Queue, error)
	QueueBind(queueName string, routingKey string) error
	UserInit() (config.UserConfig, *config.App)
	AppInit(usercfg config.UserConfig) config.App
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	storage   UserStorageProviver
	clientRMQ ClientRMQProvider
	Cfg       config.RabbitmqSrv
}

// New.
func NewUserData(s UserStorageProviver, clientrmq ClientRMQProvider, cfg config.RabbitmqSrv) *UserServices {
	return &UserServices{
		storage:   s,
		clientRMQ: clientrmq,
		Cfg:       cfg,
	}
}

// CreateUser.
func (s *UserServices) CreateUser(ctx context.Context, login string, psw string) (config.UserConfig, error) {
	log.Print("TEST USER C")

	// проверка существования пользователя
	ok, err := s.storage.IsUserLoginExist(ctx, login)
	if ok {
		log.Printf("login \"%s\" already exist:", login)
		return config.UserConfig{}, err
	}
	if err != nil {
		log.Print("check login error: ", err)
		return config.UserConfig{}, err
	}
	// создание uidcfg
	// создание конфигурации rmq
	// одно приложение - нет rk
	usercfg, userapp := s.clientRMQ.UserInit()
	// bindings отсутвует, т.к. одно приложение
	_, err = s.clientRMQ.QueueDeclare(
		userapp.ConsumeQueue, // name
	)
	if err != nil {
		log.Print("rabbitmq queue creation error: ", err)
		return config.UserConfig{}, err
	}
	// сохраняем в хранилище
	err = s.storage.CreateUser(ctx, login, psw, usercfg.UserID, usercfg)
	if err != nil {
		log.Print("usercfg creating in storage error: ", err)
		return config.UserConfig{}, err
	}
	// возращаем приложению клиента
	return usercfg, err
}

// CreateUser получаем psw хешированный base64 и .
func (s *UserServices) CreateApp(ctx context.Context, uid string, psw string) (string, config.UserConfig, error) {
	// проекрка логина и пароля пользователя
	ok, err := s.storage.CheckPsw(ctx, uid, psw)
	if !ok {
		log.Print("uid or psw incorret")
		return "", config.UserConfig{}, errors.New("uid or psw incorret")
	}
	if err != nil {
		log.Print("check psw error: ", err)
		return "", config.UserConfig{}, err
	}
	// получаем конфигурацию из хранилища
	usercfg, err := s.storage.ReadUserCfg(ctx, uid)
	if err != nil {
		log.Print("read usercfg from storage error: ", err)
		return "", config.UserConfig{}, err
	}
	// генерируем AppID
	// генерируем очередь и routingkey
	userapp := s.clientRMQ.AppInit(usercfg)
	// добавляем очереди для нового приложения пользователя
	q, err := s.clientRMQ.QueueDeclare(userapp.ConsumeQueue)
	if err != nil {
		log.Print("rabbitmq queue creation error: ", err)
		return "", config.UserConfig{}, err
	}
	for i := range usercfg.Apps[0].ExchangeBindings {
		userapp.ExchangeBindings = append(userapp.ExchangeBindings, usercfg.Apps[0].ExchangeBindings[i] )
		err = s.clientRMQ.QueueBind(
			q.Name, // queue name
			usercfg.Apps[0].ExchangeBindings[i])          // routing key
			if err != nil {
			log.Print("rabbitmq queue bindings error: ", err)
			return "", config.UserConfig{}, err
		}
	}
	userapp.ExchangeBindings = append(userapp.ExchangeBindings, usercfg.Apps[0].ConsumeQueue)
	err = s.clientRMQ.QueueBind(q.Name, usercfg.Apps[0].ConsumeQueue)
	if err != nil {
		log.Print("rabbitmq queue bindings error: ", err)
		return "", config.UserConfig{}, err
	}
	// добавляем routingkey в биндинги всех приложений клиента
	// добавление routingkey в очереди rabbit (bindings)
	for i := range usercfg.Apps {
		usercfg.Apps[i].ExchangeBindings = append(usercfg.Apps[i].ExchangeBindings, userapp.RoutingKey)
		err = s.clientRMQ.QueueBind(
			usercfg.Apps[i].ConsumeQueue, // queue name
			userapp.RoutingKey)          // routing key
			if err != nil {
			log.Print("rabbitmq queue bindings error: ", err)
			return "", config.UserConfig{}, err
		}
	}
	// добавлем новое приложение в структуру конфигурации
	usercfg.Apps = append(usercfg.Apps, userapp)
	// сохраняем в хранилище
	err = s.storage.UpdateUser(ctx, uid, usercfg)
	if err != nil {
		log.Print("update usercfg Apps error: ", err)
	}
	// возвращаем AppID и конфигурацию приложению клиента
	return userapp.AppID, usercfg, err
}
