package sql

import (
	"context"
)

// CreateText.
func (ms *StorageSQL) ReadUserRecords(ctx context.Context, userID string) error {
	// создаем текст запроса
	q := `INSERT INTO text_records 
			VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5			
			)`
	// записываем в хранилице userid, id, URL PostgreSQL.
	_, err := ms.PostgreConn.ExecContext(ctx, q, userID)
	return err
}
