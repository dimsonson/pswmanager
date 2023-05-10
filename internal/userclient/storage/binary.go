package storage

import (
	"context"

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
			$7			
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
