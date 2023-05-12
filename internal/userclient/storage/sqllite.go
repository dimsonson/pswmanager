package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
	_ "github.com/mutecomm/go-sqlcipher"
)

type SQLite struct {
	db *sql.DB
	NamedParam
}

func New(dsn string) (*SQLite, error) {
	// создаем контекст и оснащаем его таймаутом
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()

	// key := "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99"
	// dbname := fmt.Sprintf("db?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", key)

	dbname := fmt.Sprintf("%s?_pragma_key=123", dsn)

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("can't connect database: %w", err)
	}

	// создаем текст запроса
	q := `CREATE TABLE IF NOT EXISTS login_records (
	"metadata" TEXT NOT NULL,
	"login" TEXT,
	"psw" TEXT,
	"uid" TEXT NOT NULL,
	"appid" TEXT NOT NULL,
	"recordid" TEXT NOT NULL UNIQUE,
	"chng_time" DATETIME,
	"deleted" INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS text_records (
	"metadata" TEXT NOT NULL,
	"textdata" TEXT,
	"uid" TEXT NOT NULL,
	"appid" TEXT NOT NULL,
	"recordid" TEXT NOT NULL UNIQUE,
	"chng_time" DATETIME,
	"deleted" BOOLEAN DEFAULT 'false' 
	);

	CREATE TABLE IF NOT EXISTS binary_records (
	"metadata" TEXT NOT NULL,
	"binary" TEXT,
	"uid" TEXT NOT NULL,
	"appid" TEXT NOT NULL,
	"recordid" TEXT NOT NULL UNIQUE,
	"chng_time" DATETIME,
	"deleted" BOOLEAN DEFAULT 'false'
	);

	CREATE TABLE IF NOT EXISTS card_records (
	"metadata" TEXT NOT NULL,
	"brand" DECIMAL,
	"num" TEXT,
	"date" TEXT,
	"code" DECIMAL,
	"holder" TEXT,
	"uid" TEXT NOT NULL,
	"appid" TEXT NOT NULL,
	"recordid" TEXT NOT NULL UNIQUE,
	"chng_time" DATETIME,
	"deleted" BOOLEAN DEFAULT 'false'
	);

	CREATE TABLE IF NOT EXISTS ucfg (
		"ulogin" TEXT NOT NULL,
		"upsw" TEXT,
		"uid" TEXT,
		"appid" TEXT,
		"key" TEXT		
	);

	`
	if _, err = db.ExecContext(ctx, q); err != nil {
		return nil, fmt.Errorf("can't create tables: %w", err)
	}

	return &SQLite{db: db}, nil
}

func (sl *SQLite) Close() {
	sl.db.Close()
}

type NamedParam string

func (sl *NamedParam) Param(name string, value string) NamedParam {
	return NamedParam(value)
}

// CryptEncoderSHA256 encodes a password with SHA256
func CryptEncoderSHA256(pass []byte, hash interface{}) []byte {
	h := sha256.Sum256(pass)
	return h[:]
}
