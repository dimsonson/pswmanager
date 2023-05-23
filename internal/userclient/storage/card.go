package storage

import (
	"context"

	"github.com/dimsonson/pswmanager/pkg/log"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

// CreateText.
func (sl *SQLite) CreateCard(ctx context.Context, record models.CardRecord) error {
	// создаем текст запроса
	q := `INSERT INTO card_records 
			VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12			
			)`
	_, err := sl.db.ExecContext(
		ctx,
		q,
		record.Metadata,
		record.Brand,
		record.Number,
		record.ValidDate,
		record.Code,
		record.Holder,
		record.UID,
		record.AppID,
		record.RecordID,
		record.ChngTime,
		false,
		false)
	return err
}

// UpdateText.
func (sl *SQLite) UpdateCard(ctx context.Context, record models.CardRecord) error {
	// создаем текст запроса
	q := `UPDATE card_records 
	SET  metadata = $1, brand = $2, num = $3, date = $4, code = $5, holder = $6
	WHERE recordid = $7 
	AND uid = $8`
	_, err := sl.db.ExecContext(
		ctx,
		q,
		record.Metadata,
		record.Brand,
		record.Number,
		record.ValidDate,
		record.Code,
		record.Holder,
		record.RecordID,
		record.UID,
	)
	return err
}

// DeleteText.
func (sl *SQLite) DeleteCard(ctx context.Context, record models.CardRecord) error {
	// создаем текст запроса
	q := `UPDATE card_records 
	SET  deleted = 1 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}

func (sl *SQLite) SearchCard(ctx context.Context, searchInput string) ([]models.CardRecord, error) {
	cardRecords := new([]models.CardRecord)
	searchInput = "%" + searchInput + "%"
	// создаем текст запроса
	q := `SELECT 
			metadata, 
			brand,
			num,
			date,
			code,
			holder,
			uid, 
			appid, 
			recordid, 
			chng_time 
			FROM card_records 
			WHERE metadata 
			LIKE $1 
			AND deleted <> 1`
	// делаем запрос в SQL, получаем строку
	rows, err := sl.db.QueryContext(ctx, q, searchInput)
	if err != nil {
		log.Print("select login_records SQL reqest error :", err)
		return nil, err
	}
	defer rows.Close()
	// пишем результат запроса в слайс
	for rows.Next() {
		cardRecord := new(models.CardRecord)
		err = rows.Scan(
			&cardRecord.Metadata,
			&cardRecord.Brand,
			&cardRecord.Number,
			&cardRecord.ValidDate,
			&cardRecord.Code,
			&cardRecord.Holder,
			&cardRecord.UID,
			&cardRecord.AppID,
			&cardRecord.RecordID,
			&cardRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan login_records error :", err)
		}
		*cardRecords = append(*cardRecords, *cardRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request text_records iteration scan error:", err)
	}
	return *cardRecords, err
}

// MarkCardSent.
func (sl *SQLite) MarkCardSent(ctx context.Context, record models.CardRecord) error {
	// создаем текст запроса
	q := `UPDATE card_records
	SET  sent = 1
	WHERE recordid = $1
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}
