package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/dimsonson/pswmanager/internal/settings"
)

type UserStorageProviver interface {
	Close()
	CreateUser(ctx context.Context, login string, psw string, uid string, usercfg []byte) error
	ReadUser(ctx context.Context, login string) (string, string, []byte, error)
	CheckPsw(ctx context.Context, uid string, psw string) (bool, error)
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	storage UserStorageProviver
	Cfg     models.RabbitmqSrv
}

// New.
func NewUserData(s UserStorageProviver, cfg models.RabbitmqSrv) *UserServices {
	return &UserServices{
		s,
		cfg,
	}
}

// CreateUser.
func (sr *UserServices) CreateUser(ctx context.Context, login string, psw string) (models.UserConfig, error) {
	var err error

	// проверка существования пользователя

	// создание uidcfg
	usercfg := models.UserConfig{}
	usercfg.UserID = uuid.New().String()
	// создание конфигурации rmq
	userapp := new(models.App)
	userapp.AppID = uuid.New().String()
	userapp.RoutingKey = fmt.Sprintf("%s.%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumeQueue = fmt.Sprintf("%s%s%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumerName = fmt.Sprintf("%s%s%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ExchangeBindings = []string{} // одно приложение - нет rk

	usercfg.Apps = append(usercfg.Apps, *userapp)
	// bindings отсутвует, т.к. одно приложение

	// сериализация для хранения в Redis
	bytesUserCfg, err := json.Marshal(usercfg)
	if err != nil {
		log.Print("usercfg encoding error: ", err)
		return models.UserConfig{}, err
	}
	// сохраняем в хранилище
	sr.storage.CreateUser(ctx, login, psw, usercfg.UserID, bytesUserCfg)
	// возращаем приложению клиента
	return usercfg, err
}

// CreateUser получаем psw хешированный base64 и .
func (sr *UserServices) CreateApp(ctx context.Context, uid string, psw string) (string, models.UserConfig, error) {
	// проекрка логина и пароля пользователя
	ok, err := sr.storage.CheckPsw(ctx, uid, psw)
	if err != nil {
		log.Print("check psw error: ", err)
		return "", models.UserConfig{}, err
	}
	if !ok {
		log.Print("uid or psw incorret")
		return "", models.UserConfig{}, errors.New("uid or psw incorret")
	}
	// получаем конфигурацию из хранилища

	// генерируем AppID

	// генерируем очередь и routingkey

	// добавлем в структуру конфигурации

	// добавляем routing key в биндинги всех приложений клиента

	// сохраняем обновленную конфигурацию в хранилище 

	// возвращаем обновленную конфигурацию приложению клиента

	// создание uidcfg
	usercfg := models.UserConfig{}
	usercfg.UserID = uuid.New().String()
	// создание конфигурации rmq
	userapp := new(models.App)
	userapp.AppID = uuid.New().String()
	userapp.RoutingKey = fmt.Sprintf("%s.%s.%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumeQueue = fmt.Sprintf("%s%s%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ConsumerName = fmt.Sprintf("%s%s%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)
	userapp.ExchangeBindings = []string{} // одно приложение - нет rk

	usercfg.Apps = append(usercfg.Apps, *userapp)
	// bindings

	// обновление exchange bindings

	// сериализация для хранения в Redis
	bytesUserCfg, err := json.Marshal(usercfg)
	if err != nil {
		log.Print("usercfg encoding error: ", err)
		return "", models.UserConfig{}, err
	}
	// сохраняем в хранилище
	sr.storage.CreateUser(ctx, uid, psw, usercfg.UserID, bytesUserCfg)
	// возращаем приложению клиента
	return "", usercfg, err
}

func (sr *UserServices) RegUser(ctx context.Context, login string, psw string) {
}

func (sr *UserServices) AuthUser(ctx context.Context, login string, psw string) {
}
