package config

import (
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"sync"

	"github.com/dimsonson/pswmanager/pkg/log"
)

// Константы по умолчанию.
const (
	defServAddr = "localhost:8080"
	//defDBlink   = "postgres://postgres:1818@localhost:5432/dbo"
	defTLS = false
)

// ServiceConfig структура конфигурации сервиса, при запуске сервиса с флагом -c/config
// и отсутствии иных флагов и переменных окружения заполняется из файла указанного
// в этом флаге или переменной окружения CONFIG.
type ServiceConfig struct {
	ServerAddress  string         `json:"server_address"`
	EnableTLS      bool           `json:"enable_tls"`
	ConfigJSONpath string         `json:"-"`
	SQLight        SQLight        `json:"sqlite"`
	GRPC           GRPC           `json:"grpc"`
	Wg             sync.WaitGroup `json:"-"`
	UserConfig
}

// GRPC.
type GRPC struct {
	Network string
	Port    string
}

type SQLight struct {
	Dsn  string  `json:"sqlite_dsn"`
	Conn *sql.DB `json:"-"`
}

// UserConfig .
type UserConfig struct {
	UserID    string
	AppID     string
	UserLogin string
	UserPsw   string
	Key       string
}

// NewConfig конструктор создания конфигурации сервера из переменных оружения,
// флагов, конфиг файла, а так же значений по умолчанию.
func New() *ServiceConfig {
	return &ServiceConfig{}
}

// Parse метод парсинга и получения значений из переменных оружения, флагов,
// конфиг файла, а так же значений по умолчанию.
func (cfg *ServiceConfig) Parse() {

	// описываем флаги
	cfgFlag := flag.String("c", "", "config json path")
	// парсим флаги в переменные
	flag.Parse()
	cfg.ConfigJSONpath = *cfgFlag
	// используем структуру cfg models.Config для хранения параментров
	// необходимых для запуска сервера
	// читаем конфигурвационный файл и парксим в стркутуру
	if cfg.ConfigJSONpath != "" {
		configFile, err := os.ReadFile(*cfgFlag)
		if err != nil {
			log.Print("reading config file error:", err)
		}
		if err == nil {
			err = json.Unmarshal(configFile, &cfg)
			if err != nil {
				log.Printf("unmarshal config file error: %s", err)
			}
		}
	}
	//сохранение congig.json
	// cfg.SQLight.Dsn = "db" // "file:./db?_auth&_auth_user=admin&_auth_pass=admin&_auth_crypt=sha1"
	// configFile, err := json.MarshalIndent(cfg, "", "  ")
	// if err != nil {
	// 	log.Printf("marshal config file error: %s", err)
	// }
	// err = os.WriteFile("config.json", configFile, 0666)
	// if err != nil {
	// 	log.Printf("write config file error: %s", err)
	// }

}
