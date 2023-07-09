package services

import (
	"context"

	"github.com/dimsonson/pswmanager/internal/masterserver/models"
)

type ServicesProvider interface {
	ProcessingText(ctx context.Context, record models.TextRecord, key string) error
	ProcessingLogin(ctx context.Context, record models.LoginRecord) error
	ProcessingBinary(ctx context.Context, record models.BinaryRecord) error
	ProcessingCard(ctx context.Context, record models.CardRecord) error
	MarkTextSent(ctx context.Context, record models.TextRecord) error
	MarkLoginSent(ctx context.Context, record models.LoginRecord) error
	MarkCardSent(ctx context.Context, record models.CardRecord) error
	MarkBinarySent(ctx context.Context, record models.BinaryRecord) error
}
