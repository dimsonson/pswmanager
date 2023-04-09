package sql

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/rs/zerolog/log"
)

// ReadUserRecords.
func (ms *StorageSQL) ReadUserRecords(ctx context.Context, userID string) (models.SetRecords, error) {
	userRecords := models.SetRecords{}
	// создаем текст запроса
	q := `SELECT metadata, textdata, uid, recordid, deleted,  chng_time FROM text_records WHERE userid = $1 AND deleted <> true`
	// делаем запрос в SQL, получаем строку
	rows, err := ms.PostgreConn.QueryContext(ctx, q, userID)
	if err != nil {
		log.Print("select text_records SQL reuest error :", err)
	}
	defer rows.Close()
	// пишем результат запроса в структуру
	for rows.Next() {
		userTextRecord := models.TextRec{}
		err = rows.Scan(&userTextRecord.Metadata, &userTextRecord.Text, &userTextRecord.UID, &userTextRecord.AppID, &userTextRecord.RecordID, &userTextRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan text_records error :", err)
		}
		userRecords.SetTextRec = append(userRecords.SetTextRec, userTextRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request text_records iteration scan error:", err)
	}

	// создаем текст запроса
	q = `SELECT metadata, "binary", uid, recordid, deleted,  chng_time FROM binary_records WHERE userid = $1 AND deleted <> true`
	// делаем запрос в SQL, получаем строку
	rows, err = ms.PostgreConn.QueryContext(ctx, q, userID)
	if err != nil {
		log.Print("select binary_records SQL reuest error :", err)
	}
	// пишем результат запроса в структуру
	for rows.Next() {
		userBinaryRecord := models.BinaryRec{}
		err = rows.Scan(&userBinaryRecord.Metadata, &userBinaryRecord.Binary, &userBinaryRecord.UID, &userBinaryRecord.AppID, &userBinaryRecord.RecordID, &userBinaryRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan binary_records error :", err)
		}
		userRecords.SetBinaryRec = append(userRecords.SetBinaryRec, userBinaryRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request binary_records iteration scan error:", err)
	}

	return userRecords, err
}
