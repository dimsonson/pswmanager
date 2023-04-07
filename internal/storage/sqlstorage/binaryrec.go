package sqlstorage

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/models"
)

// CreateText.
func (ms *StorageSQL) CreateBinary(ctx context.Context, record models.BinaryRec) error {
	// создаем текст запроса
	q := `INSERT INTO binary_records 
			VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5			
			)`
	_, err := ms.PostgreSQL.ExecContext(ctx, q, record.Metadata, record.Binary, record.UID, record.AppID, record.RecordID)
	return err
}

// UpdateText.
func (ms *StorageSQL) UpdateBinary(ctx context.Context, record models.BinaryRec) error {
	// создаем текст запроса
	q := `UPDATE binary_records	
	SET metadata = $3, "binary" = $4
	WHERE recordid = $1 
	AND uid = $2`
	_, err := ms.PostgreSQL.ExecContext(ctx, q, record.RecordID, record.UID, record.Metadata, record.Binary)
	return err
}

// DeleteText.
func (ms *StorageSQL) DeleteBinary(ctx context.Context, record models.BinaryRec) error {
	// создаем текст запроса
	q := `UPDATE binary_records 
	SET  deleted = true 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := ms.PostgreSQL.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}
