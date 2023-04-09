package sql

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/models"
)

// CreateText.
func (ms *StorageSQL) CreateLogin(ctx context.Context, record models.LoginRec) error {
	// создаем текст запроса
	q := `INSERT INTO login_records 
			VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5,
			$6			
			)`
	_, err := ms.PostgreConn.ExecContext(ctx, q, record.Metadata, record.Login, record.Psw, record.UID, record.AppID, record.RecordID)
	return err
}

// UpdateText.
func (ms *StorageSQL) UpdateLogin(ctx context.Context, record models.LoginRec) error {
	// создаем текст запроса
	q := `UPDATE login_records 
	SET  metadata = $3, login = $4, psw = $5 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := ms.PostgreConn.ExecContext(ctx, q, record.RecordID, record.UID, record.Metadata, record.Login, record.Psw)
	return err
}

// DeleteText.
func (ms *StorageSQL) DeleteLogin(ctx context.Context, record models.LoginRec) error {
	// создаем текст запроса
	q := `UPDATE login_records 
	SET  deleted = true 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := ms.PostgreConn.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}