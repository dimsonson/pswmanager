package storage

import (
	"context"

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
			$6			
			)`
	_, err := sl.db.ExecContext(ctx, q, record.Metadata, record.Login, record.Psw, record.UID, record.AppID, record.RecordID)
	return err
}

// UpdateText.
func (sl *SQLite) UpdateLogin(ctx context.Context, record models.LoginRecord) error {
	// создаем текст запроса
	q := `UPDATE login_records 
	SET  metadata = $3, login = $4, psw = $5 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID, record.Metadata, record.Login, record.Psw)
	return err
}

// DeleteText.
func (sl *SQLite) DeleteLogin(ctx context.Context, record models.LoginRecord) error {
	// создаем текст запроса
	q := `UPDATE login_records 
	SET  deleted = true 
	WHERE recordid = $1 
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}