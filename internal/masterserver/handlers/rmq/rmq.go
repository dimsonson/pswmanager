package rmq

import (
	"context"

	rmq "github.com/MashinIvan/rabbitmq"
	"github.com/dimsonson/pswmanager/internal/masterserver/config"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
	"github.com/rs/zerolog/log"
)

// ServiceProvider интерфейс методов бизнес логики.
type ServiceProviderText interface {
	ProcessingText(ctx context.Context, record models.TextRecord) error
}

type ServiceProviderLogin interface {
	ProcessingLogin(ctx context.Context, record models.LoginRecord) error
}

type ServiceProviderBinary interface {
	ProcessingBinary(ctx context.Context, record models.BinaryRecord) error
}

type ServiceProviderCard interface {
	ProcessingCard(ctx context.Context, record models.CardRecord) error
}

// Handlers структура для конструктура обработчика.
type Handlers struct {
	servText   ServiceProviderText
	servLogin  ServiceProviderLogin
	servBinary ServiceProviderBinary
	servCard   ServiceProviderCard
}

// New конструктор обработчика.
func New(txt ServiceProviderText, lg ServiceProviderLogin, bin ServiceProviderBinary, crd ServiceProviderCard) *Handlers {
	return &Handlers{
		txt,
		lg,
		bin,
		crd,
	}
}

// TextRec.
func (hnd *Handlers) TextRec(ctx context.Context, cfg config.RabbitmqSrv) func(ctx *rmq.DeliveryContext) {
	return func(ctx *rmq.DeliveryContext) {
		// process delivery
		create := models.TextRecord{}
		err := ctx.BindJSON(&create)
		if err != nil {
			log.Print(err)
		}
		err = hnd.servText.ProcessingText(ctx, create)
		if err != nil {
			return
		}
		err = ctx.Delivery.Ack(true)
		if err != nil {
			log.Print("received Ask error:", settings.ColorRed, err, settings.ColorReset)
		}
	}
}

// LoginRec.
func (hnd *Handlers) LoginRec(ctx context.Context, cfg config.RabbitmqSrv) func(ctx *rmq.DeliveryContext) {
	return func(ctx *rmq.DeliveryContext) {
		// process delivery
		loginRec := models.LoginRecord{}
		err := ctx.BindJSON(&loginRec)
		if err != nil {
			log.Print(err)
		}
		err = hnd.servLogin.ProcessingLogin(ctx, loginRec)
		if err != nil {
			return
		}
		err = ctx.Delivery.Ack(true)
		if err != nil {
			log.Print("received Ask error:", settings.ColorRed, err, settings.ColorReset)
		}
	}
}

// BinaryRec.
func (hnd *Handlers) BinaryRec(ctx context.Context, cfg config.RabbitmqSrv) func(ctx *rmq.DeliveryContext) {
	return func(ctx *rmq.DeliveryContext) {
		binaryRec := models.BinaryRecord{}
		err := ctx.BindJSON(&binaryRec)
		if err != nil {
			log.Print(err)
		}
		err = hnd.servBinary.ProcessingBinary(ctx, binaryRec)
		if err != nil {
			return
		}
		err = ctx.Delivery.Ack(true)
		if err != nil {
			log.Print("received Ask error:", settings.ColorRed, err, settings.ColorReset)
		}
	}
}

// CardRec.
func (hnd *Handlers) CardRec(ctx context.Context, cfg config.RabbitmqSrv) func(ctx *rmq.DeliveryContext) {
	return func(ctx *rmq.DeliveryContext) {
		// process delivery
		cardRec := models.CardRecord{}
		err := ctx.BindJSON(&cardRec)
		if err != nil {
			log.Print(err)
		}
		err = hnd.servCard.ProcessingCard(ctx, cardRec)
		if err != nil {
			return
		}
		err = ctx.Delivery.Ack(true)
		if err != nil {
			log.Print("received Ask error:", settings.ColorRed, err, settings.ColorReset)
		}
	}
}
