package storage

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

// CreateText.
func (sl *SQLite) CreateText(ctx context.Context, record models.TextRecord) error {
	// создаем текст запроса
	q := `INSERT INTO text_records
			VALUES (
			$1,
			$2,
			$3,
			$4,
			$5, 
			$6,
			$7
			)`
	_, err := sl.db.ExecContext(ctx, q, record.Metadata, record.Text, record.UID, record.AppID, record.RecordID, record.ChngTime, false)
	return err
}

// UpdateText.
func (sl *SQLite) UpdateText(ctx context.Context, record models.TextRecord) error {
	// создаем текст запроса
	q := `UPDATE text_records 
	SET metadata = @Metadata, 
	textdata = @Text 
	WHERE recordid = @recordID 
	AND uid = @UID`
	// q := `UPDATE text_records
	// SET metadata = $1,
	// textdata = $2
	// WHERE recordid = $3
	// AND uid = $4`
	_, err := sl.db.ExecContext(ctx, q,
		sl.Param("Metadata", record.Metadata),
		sl.Param("Text", record.Text),
		sl.Param("recordID", record.RecordID),
		sl.Param("UID", record.UID))
	//	_, err := sl.db.ExecContext(ctx, q, record.Metadata, record.Text, record.RecordID, record.UID)
	return err
}

// DeleteText.
func (sl *SQLite) DeleteText(ctx context.Context, record models.TextRecord) error {
	// создаем текст запроса
	q := `UPDATE text_records
	SET  deleted = 1
	WHERE recordid = $1
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}

func (sl *SQLite) SearchText(ctx context.Context, searchInput string) ([]models.TextRecord, error) {
	textRecords := new([]models.TextRecord)
	searchInput = "%" + searchInput + "%"
	// создаем текст запроса
	q := `SELECT metadata, textdata, uid, appid, recordid, chng_time FROM text_records WHERE metadata LIKE $1 AND deleted <> 1`
	// делаем запрос в SQL, получаем строку
	rows, err := sl.db.QueryContext(ctx, q, searchInput)
	if err != nil {
		log.Print("select login_records SQL reqest error :", err)
		return nil, err
	}
	defer rows.Close()
	// пишем результат запроса в слайс
	for rows.Next() {
		textRecord := new(models.TextRecord)
		err = rows.Scan(
			&textRecord.Metadata,
			&textRecord.Text,
			&textRecord.UID,
			&textRecord.AppID,
			&textRecord.RecordID,
			&textRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan login_records error :", err)
		}
		*textRecords = append(*textRecords, *textRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request text_records iteration scan error:", err)
	}
	return *textRecords, err
}
