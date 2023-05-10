package storage

import (
	"context"
	"log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

// CreateText.
func (sl *SQLite) CreateLogin(ctx context.Context, record models.LoginRecord) error {
	// создаем текст запроса
	q := `INSERT INTO login_records 
			VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5,
			$6,
			$7,
			$8			
			)`
	_, err := sl.db.ExecContext(ctx, q, record.Metadata, record.Login, record.Psw, record.UID, record.AppID, record.RecordID, false, record.ChngTime)
	return err
}

// UpdateText.
func (sl *SQLite) UpdateLogin(ctx context.Context, record models.LoginRecord) error {
	// создаем текст запроса
	q := `UPDATE login_records 
	SET  metadata = $1, login = $2, psw = $3 
	WHERE recordid = $4 
	AND uid = $5`
	_, err := sl.db.ExecContext(ctx, q, record.Metadata, record.Login, record.Psw, record.RecordID, record.UID)
	return err
}

// DeleteText.
func (sl *SQLite) DeleteLogin(ctx context.Context, record models.LoginRecord) error {
	// создаем текст запроса
	q := `UPDATE login_records 
	SET  deleted = 1 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}

func (sl *SQLite) SearchLogin(ctx context.Context, record models.LoginRecord) error {
	userRecords := new(models.SetOfRecords)
	// создаем текст запроса
	q := `SELECT metadata, textdata, uid, appid, recordid, chng_time FROM text_records WHERE uid = $1 AND deleted <> true`
	// делаем запрос в SQL, получаем строку
	rows, err := sl.db.QueryContext(ctx, q, record.RecordID)
	if err != nil {
		log.Print("select text_records SQL reuest error :", err)
	}
	defer rows.Close()
	// пишем результат запроса в структуру
	for rows.Next() {
		userTextRecord := new(models.TextRecord)
		err = rows.Scan(&userTextRecord.Metadata, &userTextRecord.Text, &userTextRecord.UID, &userTextRecord.AppID, &userTextRecord.RecordID, &userTextRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan text_records error :", err)
		}
		userRecords.SetTextRec = append(userRecords.SetTextRec, *userTextRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request text_records iteration scan error:", err)
	}
	return err
}
