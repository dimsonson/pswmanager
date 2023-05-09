package storage

import (
	"context"
	"time"

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
	_, err := sl.db.ExecContext(ctx, q, record.Metadata, record.Text, record.UID, record.AppID, record.RecordID, false, time.Now())
	return err
}

// UpdateText.
func (sl *SQLite) UpdateText(ctx context.Context, record models.TextRecord) error {
	// создаем текст запроса
	q := `UPDATE text_records
	SET  metadata = $3, textdata = $4
	WHERE recordid = $1
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID, record.Metadata, record.Text)
	return err
}

// DeleteText.
func (sl *SQLite) DeleteText(ctx context.Context, record models.TextRecord) error {
	// создаем текст запроса
	q := `UPDATE text_records
	SET  deleted = true
	WHERE recordid = $1
	AND uid = $2`
	_, err := sl.db.ExecContext(ctx, q, record.RecordID, record.UID)
	return err
}
