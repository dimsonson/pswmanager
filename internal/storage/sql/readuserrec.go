package sql

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/models"
	"github.com/rs/zerolog/log"
)

// ReadUserRecords.
func (ms *StorageSQL) ReadUserRecords(ctx context.Context, userID string) (*models.SetRecords, error) {
	userRecords := new(models.SetRecords)
	// создаем текст запроса
	q := `SELECT metadata, textdata, uid, appid, recordid, chng_time FROM text_records WHERE uid = $1 AND deleted <> true`
	// делаем запрос в SQL, получаем строку
	rows, err := ms.PostgreConn.QueryContext(ctx, q, userID)
	if err != nil {
		log.Print("select text_records SQL reuest error :", err)
	}
	defer rows.Close()
	// пишем результат запроса в структуру
	for rows.Next() {
		userTextRecord := new(models.TextRec)
		err = rows.Scan(&userTextRecord.Metadata, &userTextRecord.Text, &userTextRecord.UID, &userTextRecord.AppID, &userTextRecord.RecordID, &userTextRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan text_records error :", err)
		}
		userRecords.SetTextRec = append(userRecords.SetTextRec, *userTextRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request text_records iteration scan error:", err)
	}

	// создаем текст запроса
	q = `SELECT metadata, "binary", uid, appid, recordid, chng_time FROM binary_records WHERE uid = $1 AND deleted <> true`
	// делаем запрос в SQL, получаем строку
	rows, err = ms.PostgreConn.QueryContext(ctx, q, userID)
	if err != nil {
		log.Print("select binary_records SQL reuest error :", err)
	}
	// пишем результат запроса в структуру
	for rows.Next() {
		userBinaryRecord := new(models.BinaryRec)
		err = rows.Scan(&userBinaryRecord.Metadata, &userBinaryRecord.Binary, &userBinaryRecord.UID, &userBinaryRecord.AppID, &userBinaryRecord.RecordID, &userBinaryRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan binary_records error :", err)
		}
		userRecords.SetBinaryRec = append(userRecords.SetBinaryRec, *userBinaryRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request binary_records iteration scan error:", err)
	}

	// создаем текст запроса
	q = `SELECT metadata, login, psw, uid, appid, recordid, chng_time FROM login_records WHERE uid = $1 AND deleted <> true`
	// делаем запрос в SQL, получаем строку
	rows, err = ms.PostgreConn.QueryContext(ctx, q, userID)
	if err != nil {
		log.Print("select login_records SQL reuest error :", err)
	}
	// пишем результат запроса в структуру
	for rows.Next() {
		userLoginRecord := new(models.LoginRec)
		err = rows.Scan(&userLoginRecord.Metadata, &userLoginRecord.Login, &userLoginRecord.Psw, &userLoginRecord.UID, &userLoginRecord.AppID, &userLoginRecord.RecordID, &userLoginRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan login_records error :", err)
		}
		userRecords.SetLoginRec = append(userRecords.SetLoginRec, *userLoginRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request login_records iteration scan error:", err)
	}

	// создаем текст запроса
	q = `SELECT metadata, brand, num, date, code, uid, appid, recordid, chng_time FROM card_records WHERE uid = $1 AND deleted <> true`
	// делаем запрос в SQL, получаем строку
	rows, err = ms.PostgreConn.QueryContext(ctx, q, userID)
	if err != nil {
		log.Print("select card_records SQL reuest error :", err)
	}
	// пишем результат запроса в структуру
	for rows.Next() {
		userCardRecord := new(models.CardRec)
		err = rows.Scan(&userCardRecord.Metadata, &userCardRecord.Brand, &userCardRecord.Number, &userCardRecord.ValidDate, &userCardRecord.Code, &userCardRecord.UID, &userCardRecord.AppID, &userCardRecord.RecordID, &userCardRecord.ChngTime)
		if err != nil {
			log.Print("row by row scan card_records error :", err)
		}
		userRecords.SetCardRec = append(userRecords.SetCardRec, *userCardRecord)
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Print("request card_records iteration scan error:", err)
	}
	return userRecords, err
}
