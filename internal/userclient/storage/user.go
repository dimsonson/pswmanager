package storage

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/userclient/config"

	"github.com/dimsonson/pswmanager/pkg/log"
)

// добавление нового пользователя в хранилище, запись в две таблицы в транзакции
func (sl *SQLite) CreateUser(ctx context.Context, ucfg config.UserConfig) error {
//func (sl *SQLite) CreateUser(ctx context.Context, ucfg config.UserConfig, keyDB string) error {
	// создаем текст запроса
	q := `INSERT INTO ucfg
			VALUES (
			$1,
			$2,
			$3,
			$4,
			$5			
			)`
	_, err := sl.db.ExecContext(ctx, q, ucfg.UserLogin, ucfg.UserPsw, ucfg.UserID, ucfg.AppID, ucfg.Key)
	return err
}

// проверка наличия нового пользователя в хранилище - авторизация
func (sl *SQLite) ReadUser(ctx context.Context) (config.UserConfig, error) {
	ucfg := config.UserConfig{}
	// создаем текст запроса
	q := `SELECT ulogin, upsw, uid, appid, key FROM ucfg`
	// делаем запрос в SQL, получаем строку и пишем результат запроса в пременную
	err := sl.db.QueryRowContext(ctx, q).Scan(&ucfg.UserLogin, &ucfg.UserPsw, &ucfg.UserID, &ucfg.AppID, &ucfg.Key)
	if err != nil {
		log.Printf("select SQL request scan error: %s", err)
		return ucfg, err
	}
	return ucfg, err
}

// проверка наличия нового пользователя в хранилище - авторизация
func (sl *SQLite) CheckUser(ctx context.Context, login string) (string, error) {
	var passwDB string
	// создаем текст запроса
	q := `SELECT upsw FROM ucfg WHERE ulogin = $1`
	// делаем запрос в SQL, получаем строку и пишем результат запроса в пременную
	err := sl.db.QueryRowContext(ctx, q, login).Scan(&passwDB)
	if err != nil {
		log.Printf("select CheckUser request scan error: %s", err)
		return "", err
	}

	return passwDB, err
}
