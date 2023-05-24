package storage

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

// CreateText.
func (sl *SQLite) CreateBinary(ctx context.Context, record models.BinaryRecord) error {
	// создаем текст запроса
	q := `INSERT INTO binary_records 
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
	_, err := sl.db.ExecContext(
		ctx,
		q,
		record.Metadata,
		record.Binary,
		record.UID,
		record.AppID,
		record.RecordID,
		record.ChngTime,
		false,
		false)
	return err
}

// UpdateText.
func (sl *SQLite) UpdateBinary(ctx context.Context, record models.BinaryRecord) error {
	// создаем текст запроса
	q := `UPDATE binary_records	
	SET metadata = $3, "binary" = $4
	WHERE recordid = $1 
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID, record.Metadata, record.Binary)
	return err
}

// DeleteText.
func (sl *SQLite) DeleteBinary(ctx context.Context, record models.BinaryRecord) error {
	// создаем текст запроса
	q := `UPDATE binary_records 
	SET  deleted = 1 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}

func (sl *SQLite) SearchBinary(ctx context.Context, searchInput string) ([]models.BinaryRecord, error) {
	binaryRecords := new([]models.BinaryRecord)
	searchInput = "%" + searchInput + "%"
	// создаем текст запроса
	q := `SELECT metadata, binary, uid, appid, recordid, chng_time FROM binary_records WHERE metadata LIKE $1 AND deleted <> 1`
	// делаем запрос в SQL, получаем строку
	rows, err := sl.db.QueryContext(ctx, q, searchInput)
	if err != nil {
		log.Print("select login_records SQL reqest error :", err)
		return nil, err
	}
	defer rows.Close()
	// пишем результат запроса в слайс
	for rows.Next() {
		binaryRecord := new(models.BinaryRecord)
		err = rows.Scan(
			&binaryRecord.Metadata,
			&binaryRecord.Binary,
			&binaryRecord.UID,
			&binaryRecord.AppID,
			&binaryRecord.RecordID,
			&binaryRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan login_records error :", err)
		}
		*binaryRecords = append(*binaryRecords, *binaryRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request text_records iteration scan error:", err)
	}
	return *binaryRecords, err
}

// MarkBinarySent.
func (sl *SQLite) MarkBinarySent(ctx context.Context, record models.BinaryRecord) error {
	// создаем текст запроса
	q := `UPDATE binary_records
	SET  sent = 1
	WHERE recordid = $1
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}
