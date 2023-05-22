package nosql

import (
	"context"
	"encoding/json"

	"github.com/dimsonson/pswmanager/internal/masterserver/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
	"github.com/dimsonson/pswmanager/pkg/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

// StorageNoSQL структура хранилища Redis.
type StorageNoSQL struct {
	RedisNoSQL *redis.Client
}

// NewNoSQLStorage конструктор нового хранилища PostgreSQL.
func New(cfg config.Redis) *StorageNoSQL {
	//redis.SetLogger(internal.Logger)

	// создаем контекст и оснащаем его таймаутом
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		Username: cfg.Username,
		//TLSConfig: &tls.Config{
		//MinVersion: tls.VersionTLS12,
		//Certificates: []tls.Certificate{cert}
		//},
	})

	log.Print("redisDB ", rdb.Ping(ctx))

	return &StorageNoSQL{
		RedisNoSQL: rdb,
	}
}

func (rdb *StorageNoSQL) Close() {
	rdb.RedisNoSQL.Close()
}

func (rdb *StorageNoSQL) CreateUser(ctx context.Context, login string, psw string, uid string, usercfg config.UserConfig) error {
	// сериализация для хранения в Redis
	bytesUserCfg, err := json.Marshal(usercfg)
	if err != nil {
		log.Print("usercfg encoding error: ", err)
		return err
	}

	pipe := rdb.RedisNoSQL.Pipeline()

	err = pipe.HSet(ctx, "login", login, uid).Err()
	if err != nil {
		log.Print("login set to redis error: ", err)
		return err
	}

	err = pipe.HSet(ctx, "psw", uid, psw).Err()
	if err != nil {
		log.Print("psw set to redis error: ", err)
		return err
	}

	err = pipe.HSet(ctx, "key", uid, usercfg.CryptoKey).Err()
	if err != nil {
		log.Print("psw set to redis error: ", err)
		return err
	}

	err = pipe.HSet(ctx, "usercfg", uid, bytesUserCfg).Err()
	if err != nil {
		log.Print("usercfg set to redis error: ", err)
		return err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Print("pipe redis error: ", err)
		return err
	}

	return err
}

func (rdb *StorageNoSQL) ReadUserCfg(ctx context.Context, uid string) (config.UserConfig, error) {
	cmd := rdb.RedisNoSQL.HGet(ctx, "usercfg", uid)
	bytesUserCfg, err := cmd.Bytes()
	if err != nil {
		log.Print("usercfg get from redis error: ", err)
		return config.UserConfig{}, err
	}
	var usercfg config.UserConfig
	err = json.Unmarshal(bytesUserCfg, &usercfg)
	if err != nil {
		log.Print("usercfg decoding error: ", err)
		return config.UserConfig{}, err
	}
	return usercfg, err
}

func (rdb *StorageNoSQL) CheckPsw(ctx context.Context, uid string, psw string) (bool, error) {
	pswStorage, err := rdb.RedisNoSQL.HGet(ctx, "psw", uid).Result()
	if err != nil {
		log.Print("psw check redis error: ", err)
		return false, err
	}
	return psw == pswStorage, err
}

func (rdb *StorageNoSQL) UpdateUser(ctx context.Context, uid string, usercfg config.UserConfig) error {
	// сохраняем обновленную конфигурацию в хранилище
	// сериализация для хранения в Redis
	bytesUserCfg, err := json.Marshal(usercfg)
	if err != nil {
		log.Print("usercfg encoding error: ", err)
		return err
	}
	err = rdb.RedisNoSQL.HSet(ctx, "usercfg", uid, bytesUserCfg).Err()
	if err != nil {
		log.Print("usercfg set to redis error: ", err)
		return err
	}
	return err
}

func (rdb *StorageNoSQL) IsUserLoginExist(ctx context.Context, login string) (string, bool, error) {
	uid, err := rdb.RedisNoSQL.HGet(ctx, "login", login).Result()
	if err != nil {
		if err == redis.Nil {
			log.Print("login key doenst exist: ", err)
			return uid, false, nil
		}
		log.Print("login check redis error: ", err)
		return uid, true, err
	}
	return uid, true, err
}
