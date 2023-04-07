package nosqlstorage

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/dimsonson/pswmanager/internal/settings"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// StorageSQL структура хранилища PostgreSQL.
type StorageNoSQL struct {
	RedisNoSQL *redis.Client
}

// NewNoSQLStorage конструктор нового хранилища PostgreSQL.
func New(p string) *StorageNoSQL {
	//redis.SetLogger(internal.Logger)

	// создаем контекст и оснащаем его таймаутом
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		//TLSConfig: &tls.Config{
		//	MinVersion: tls.VersionTLS12,
		//Certificates: []tls.Certificate{cert}
		//},
	})

	log.Print(rdb.Ping(ctx))

	return &StorageNoSQL{
		RedisNoSQL: rdb,
	}
}

func (rdb *StorageNoSQL) Close() {
	rdb.RedisNoSQL.Close()
}

func (rdb *StorageNoSQL) CreateUser(ctx context.Context, login string, psw string, uid string, usercfg models.UserConfig) error {

	err := rdb.RedisNoSQL.HSet(ctx, "login", login, uid).Err()
	if err != nil {
		log.Print("login set to redis error: ", err)
		return err
	}

	err = rdb.RedisNoSQL.HSet(ctx, "psw", uid, psw).Err()
	if err != nil {
		log.Print("psw set to redis error: ", err)
		return err
	}

	err = rdb.RedisNoSQL.HSet(ctx, "usercfg", uid, usercfg).Err()
	if err != nil {
		log.Print("login set to redis error: ", err)
		return err
	}

	return err
}
