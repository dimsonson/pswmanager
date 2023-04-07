package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type ExampleUser struct {
	Name string `redis:"name"`
	Age  int    `redis:"age"`
}

func main() {

	//redis.SetLogger(log)

	key := "u1:passw33"
	hkey := "u1:app1:consume_q:binds_h"
	v := "password123"

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		//TLSConfig: &tls.Config{
		//	MinVersion: tls.VersionTLS12,
		//Certificates: []tls.Certificate{cert}
		//},
	})

	defer rdb.Close()

	fmt.Println(rdb.Ping(ctx))

	err := rdb.Set(ctx, key, v, 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(key, val)

	items := ExampleUser{"jane", 22}

	err = rdb.HSet(ctx, hkey, items).Err()
	if err != nil {
		panic(err)
	}

	val22, err := rdb.HKeys(ctx, hkey).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val22)

	val2, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println(key, val2)
	}
	// Output: key value
	// key2 does not exist

	type user struct {
		user_id   string
		upassw    string
		ulogin    string
		apps_rk   map[string]string
		consume_q []string
	}

	var usr user
	usr.apps_rk = make(map[string]string)
	usr.user_id = "u888"
	usr.ulogin = "d@ya.ru"
	usr.upassw = "password12345678"
	usr.consume_q = []string{"u1app2", "u1app3", "u1app5"}
	usr.apps_rk["app1"] = "all.u8.app1"
	usr.apps_rk["app2"] = "all.u8.app2"
	usr.apps_rk["app3"] = "all.u8.app3"

	ulogin := fmt.Sprintf("%s:ulogin", usr.user_id)
	err = rdb.Set(ctx, ulogin, usr.ulogin, 0).Err()
	if err != nil {
		panic(err)
	}

	consume_q := fmt.Sprintf("%s:consume_q", usr.user_id)
	err = rdb.LPush(ctx, consume_q, usr.consume_q).Err()
	if err != nil {
		panic(err)
	}

	upassw := fmt.Sprintf("%s:upassw", usr.user_id)
	err = rdb.Set(ctx, upassw, usr.upassw, 0).Err()
	if err != nil {
		panic(err)
	}

	type MyHash struct {
		Key1 string `redis:"key1"`
		Key2 int    `redis:"key2"`
	}

	err = rdb.HSet(ctx, "myhash", MyHash{"value1", 89}).Err()
	if err != nil {
		panic(err)
	}

	match := fmt.Sprintf("%s:ulogin", usr.user_id)
	iter := rdb.Scan(ctx, 0, match, 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("keys", iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

}

type users struct {
	login_uid     map[string]string // map[login]uid
	uid_psw       map[string]string // map[uid]psw
	app_uid       map[string]string // map[app]uid
	app_exch      map[string]string // map[app]exch
	app_route     map[string]string // map[app]pub_route_key
	app_consume_q map[string]string // map[app]uid
	consume_q     []string          // []route_keys
}

// регистрация пользователя
// запись [login]uid
// запись [uid]psw

// создание exchange
// создание очереди
// создание bindings - добавление всех routing key пользователя в bindings за исключнием pub routing key этого пользователя
// выборка всех routing key по пользователю user
// создание pub routing key

// запись pub routing key во все очереди пользователя кроме очерерди текущего приложения

// publishing по exchange и routing key
// consume по consume_queue
