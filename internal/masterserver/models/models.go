package models

import (
	"time"
)

// MsgType тип исплльзуемый для проставления признака типа сообщения / операции.
type MsgType int

// Константы типа MsgType исплльзуемые для проставления признака типа сообщения / операции.
const (
	Create MsgType = iota + 1
	Read
	Update
	Delete
)

// CardType тип исплльзуемый для проставления признака типа банковской карты.
type CardType int

// Константы типа CardType исплльзуемые для проставления признака типа банковской карты.
const (
	Mir CardType = iota 
	MasterCard
	Visa
	AmEx
)

// SetRecords .
type SetOfRecords struct {
	SetLoginRec  []LoginRecord
	SetTextRec   []TextRecord
	SetBinaryRec []BinaryRecord
	SetCardRec   []CardRecord
}

// LoginRec структура сообщния для опараций с парами логин/пароль.
type LoginRecord struct {
	RecordID  string
	ChngTime  time.Time
	UID       string
	AppID     string
	Login     string
	Psw       string
	Metadata  string
	Operation MsgType
}

// LoginRec структура сообщния для опараций с текстовыми данными пользователя.
type TextRecord struct {
	RecordID  string
	ChngTime  time.Time
	UID       string
	AppID     string
	Text      string
	Metadata  string
	Operation MsgType
}

// BinaryRec структура сообщния для опараций с бинарными данными пользователя.
type BinaryRecord struct {
	RecordID  string
	ChngTime  time.Time
	UID       string
	AppID     string
	Binary    string
	Metadata  string
	Operation MsgType
}

// CardRec структура сообщния для опараций с данными карт пользователя.
type CardRecord struct {
	RecordID  string
	ChngTime  time.Time
	UID       string
	AppID     string
	Brand     string
	Number    string
	ValidDate string
	Code      string
	Holder    string
	Metadata  string
	Operation MsgType
}

// Queue.
type Queue struct {
	Name      string // server confirmed or generated name
	Messages  int    // количество сообщений не требующих ask
	Consumers int    // количество потребителей, получающих сообщения
}
