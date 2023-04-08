package sql

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/pswmanager/internal/settings"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// StorageSQL структура хранилища PostgreSQL.
type StorageSQL struct {
	PostgreSQL *sql.DB
}

// NewSQLStorage конструктор нового хранилища PostgreSQL.
func New(p string) *StorageSQL {
	// создаем контекст и оснащаем его таймаутом
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()
	// открываем базу данных
	db, err := sql.Open("pgx", p)
	if err != nil {
		log.Print("database opening error:", settings.ColorRed, err, settings.ColorReset)
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
			"uid" TEXT NOT NULL,
			"appid" TEXT NOT NULL,
			"recordid" TEXT NOT NULL UNIQUE,
			"chng_time" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
			"deleted" BOOLEAN DEFAULT 'false'
			);
		
		`
	// создаем таблицу в SQL базе, если не существует
	_, err = db.ExecContext(ctx, q)
	if err != nil {
		log.Print("request NewSQLStorage to sql db returned error:", settings.ColorRed, err, settings.ColorReset)
	}
	return &StorageSQL{
		PostgreSQL: db,
	}
}

func (ms *StorageSQL) Close() {
	ms.PostgreSQL.Close()
}
