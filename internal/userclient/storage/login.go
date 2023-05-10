package storage

import (
	"context"
	"github.com/dimsonson/pswmanager/pkg/log"

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
	_, err := sl.db.ExecContext(ctx, q, record.Metadata, record.Login, record.Psw, record.UID, record.AppID, record.RecordID, record.ChngTime, false)
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

func (sl *SQLite) SearchLogin(ctx context.Context, searchInput string) ([]models.LoginRecord, error) {
	loginRecords := new([]models.LoginRecord)
	searchInput = "%"+searchInput+"%"
	// создаем текст запроса
	q := `SELECT metadata, login, psw, uid, appid, recordid, chng_time FROM login_records WHERE metadata LIKE $1 AND deleted <> 1`
	// делаем запрос в SQL, получаем строку
	rows, err := sl.db.QueryContext(ctx, q, searchInput)
	if err != nil {
		log.Print("select login_records SQL reqest error :", err)
		return nil, err
	}
	defer rows.Close()
	// пишем результат запроса в слайс
	for rows.Next() {
		loginRecord := new(models.LoginRecord)
		err = rows.Scan(
			&loginRecord.Metadata,
			&loginRecord.Login,
			&loginRecord.Psw,
			&loginRecord.UID,
			&loginRecord.AppID,
			&loginRecord.RecordID,
			&loginRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan login_records error :", err)
		}
		*loginRecords = append(*loginRecords, *loginRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request text_records iteration scan error:", err)
	}
	return *loginRecords, err
}
