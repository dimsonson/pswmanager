package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/dimsonson/pswmanager/internal/settings"
)

type UserStorageProviver interface {
	Close()
	CreateUser(ctx context.Context, login string, psw string, uid string, usercfg models.UserConfig) error
}

// Services структура конструктора бизнес логики.
type UserServices struct {
	storage UserStorageProviver
}

// New.
func NewUserData(s UserStorageProviver, cfg models.RabbitmqSrv) *UserServices {
	return &UserServices{
		s,
	}
}

// CreateUser.
func (sr *UserServices) Create(ctx context.Context, login string, psw string) (*models.UserConfig, error) {
	var err error
	// проверка что логин не существует

	// создание uidcfg
	usercfg := new(models.UserConfig)
	usercfg.UserID = uuid.New().String()

	// создание конфигурации rmq

	userapp := new(models.App)
	userapp.AppID = uuid.New().String()
	userapp.RoutingKey = fmt.Sprintf("%s%s%s", settings.MasterQueue, usercfg.UserID, userapp.AppID)

	usercfg.Apps = append(usercfg.Apps, *userapp)

	// bindings

	// сохраняем в хранилище
	sr.storage.CreateUser(ctx, login, psw, usercfg.UserID, *usercfg)

	// возращаем приложению клиента

	//

	return usercfg, err
}

// CreateUser.
// func (sr *UserServices) CreateApp(ctx context.Context, record models.BinaryRec) error {
// 	var err error
// 	switch record.Operation {
// 	case models.Create:
// 		err := sr.storage.CreateBinary(ctx, record)
// 		if err != nil {
// 			log.Print("create binary record error: ", err)
// 		}
// 		return err
// 	case models.Update:
// 		err := sr.storage.UpdateBinary(ctx, record)
// 		if err != nil {
// 			log.Print("update binary record error: ", err)
// 		}
// 		return err
// 	case models.Delete:
// 		err := sr.storage.DeleteBinary(ctx, record)
// 		if err != nil {
// 			log.Print("delete binary record error: ", err)
// 		}
// 		return err
// 	}
// 	return err
// }
