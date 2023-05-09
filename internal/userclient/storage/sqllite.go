package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
	_ "github.com/mutecomm/go-sqlcipher"
)

type SQLite struct {
	db *sql.DB
}

func New(dsn string) (*SQLite, error) {
	// создаем контекст и оснащаем его таймаутом
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()

	// key := "2DD29CA851E7B56E4697B0E1F08507293D761A05CE4D1B628663F411A8086D99"
	// dbname := fmt.Sprintf("db?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", key)

	db, err := sql.Open("sqlite3", "db?_pragma_key=123") // dbname)
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
	"chng_time" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"deleted" BOOLEAN DEFAULT 'false'
	);

	CREATE TABLE IF NOT EXISTS text_records (
	"metadata" TEXT NOT NULL,
	"textdata" TEXT,
	"uid" TEXT NOT NULL,
	"appid" TEXT NOT NULL,
	"recordid" TEXT NOT NULL UNIQUE,
	"deleted" BOOLEAN DEFAULT 'false', 
	"chng_time" timestamp with time zone DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS binary_records (
	"metadata" TEXT NOT NULL,
	"binary" TEXT,
	"uid" TEXT NOT NULL,
	"appid" TEXT NOT NULL,
	"recordid" TEXT NOT NULL UNIQUE,
	"deleted" BOOLEAN DEFAULT 'false',
	"chng_time" timestamp with time zone DEFAULT CURRENT_TIMESTAMP
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
	"chng_time" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	"deleted" BOOLEAN DEFAULT 'false'
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
