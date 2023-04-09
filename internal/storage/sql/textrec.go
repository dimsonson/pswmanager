package sql

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/models"
)

// CreateText.
func (ms *StorageSQL) CreateText(ctx context.Context, record models.TextRec) error {
	// создаем текст запроса
	q := `INSERT INTO text_records 
			VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5			
			)`
	_, err := ms.PostgreConn.ExecContext(ctx, q, record.Metadata, record.Text, record.UID, record.AppID, record.RecordID)
	return err
}

// UpdateText.
func (ms *StorageSQL) UpdateText(ctx context.Context, record models.TextRec) error {
	// создаем текст запроса
	q := `UPDATE text_records 
	SET  metadata = $3, textdata = $4
	WHERE recordid = $1 
	AND uid = $2`
	_, err := ms.PostgreConn.ExecContext(ctx, q, record.RecordID, record.UID, record.Metadata, record.Text)
	return err
}

// DeleteText.
func (ms *StorageSQL) DeleteText(ctx context.Context, record models.TextRec) error {
	// создаем текст запроса
	q := `UPDATE text_records 
	SET  deleted = true 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := ms.PostgreConn.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}
