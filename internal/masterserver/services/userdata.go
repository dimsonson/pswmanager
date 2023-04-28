package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
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
	usercfg := config.UserConfig{}
	usercfg.UserID = uuid.New().String()
	// создание конфигурации rmq
	userapp := new(config.App)
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
	userapp := config.App{}
	userapp.AppID = uuid.New().String()
	// генерируем очередь и routingkey
	userapp.ConsumeQueue = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumerName = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.RoutingKey = fmt.Sprintf("%s%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	// добавляем очередь для нового приложения пользователя
	q, err := s.clientRMQ.QueueDeclare(userapp.ConsumeQueue)
	if err != nil {
		log.Print("rabbitmq queue creation error: ", err)
		return "", config.UserConfig{}, err
	}
	err = s.clientRMQ.QueueBind(q.Name, userapp.RoutingKey)
	if err != nil {
		log.Print("rabbitmq queue bindings error: ", err)
		return "", config.UserConfig{}, err
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
			return "", config.UserConfig{}, err
		}
	}
	// добавлем новое приложение в структуру конфигурации
	usercfg.Apps = append(usercfg.Apps, userapp)
	// сохраняем в хранилище
	s.storage.UpdateUser(ctx, uid, usercfg)
	// возвращаем AppID и конфигурацию приложению клиента
	return userapp.AppID, usercfg, err
}
