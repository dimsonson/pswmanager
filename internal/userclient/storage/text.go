package storage

import (
	"context"

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
	_, err := sl.db.ExecContext(ctx, q, record.Metadata, record.Text, record.UID, record.AppID, record.RecordID, false, record.ChngTime)
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
