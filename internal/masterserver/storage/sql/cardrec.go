package sql

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

// CreateText.
func (ms *StorageSQL) CreateCard(ctx context.Context, record models.CardRec) error {
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
			$9			
			)`
	_, err := ms.PostgreConn.ExecContext(
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
		record.RecordID)
	return err
}

// UpdateText.
func (ms *StorageSQL) UpdateCard(ctx context.Context, record models.CardRec) error {
	// создаем текст запроса
	q := `UPDATE card_records 
	SET  metadata = $3, brand = $4, num = $5, date = $6, code = $7, holder = $8
	WHERE recordid = $1 
	AND uid = $2`
	_, err := ms.PostgreConn.ExecContext(
		ctx, 
		q,
		record.RecordID,
		record.UID,
		record.Metadata,
		record.Brand,
		record.Number,
		record.ValidDate,
		record.Code,
		record.Holder,
	)
	return err
}

// DeleteText.
func (ms *StorageSQL) DeleteCard(ctx context.Context, record models.CardRec) error {
	// создаем текст запроса
	q := `UPDATE card_records 
	SET  deleted = true 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := ms.PostgreConn.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}
