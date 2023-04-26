package rmq

import (
	"context"

	"github.com/rs/zerolog/log"

	rmq "github.com/MashinIvan/rabbitmq"
	"github.com/dimsonson/pswmanager/internal/masterserver/models"
	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
)

// ServiceProvider интерфейс методов бизнес логики.
type ServiceProviderText interface {
	TextRec(ctx context.Context, record models.TextRecord) error
}

type ServiceProviderLogin interface {
	LoginRec(ctx context.Context, record models.LoginRecord) error
}

type ServiceProviderBinary interface {
	BinaryRec(ctx context.Context, record models.BinaryRecord) error
}

type ServiceProviderCard interface {
	CardRec(ctx context.Context, record models.CardRecord) error
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
func (hnd *Handlers) TextRec(ctx context.Context, cfg models.RabbitmqSrv) func(ctx *rmq.DeliveryContext) {
	return func(ctx *rmq.DeliveryContext) {
		// process delivery
		create := models.TextRecord{}
		err := ctx.BindJSON(&create)
		if err != nil {
			log.Print(err)
		}
		err = hnd.servText.TextRec(ctx, create)
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
func (hnd *Handlers) LoginRec(ctx context.Context, cfg models.RabbitmqSrv) func(ctx *rmq.DeliveryContext) {
	return func(ctx *rmq.DeliveryContext) {
		// process delivery
		loginRec := models.LoginRecord{}
		err := ctx.BindJSON(&loginRec)
		if err != nil {
			log.Print(err)
		}
		err = hnd.servLogin.LoginRec(ctx, loginRec)
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
func (hnd *Handlers) BinaryRec(ctx context.Context, cfg models.RabbitmqSrv) func(ctx *rmq.DeliveryContext) {
	return func(ctx *rmq.DeliveryContext) {
		// process delivery
		binaryRec := models.BinaryRecord{}
		err := ctx.BindJSON(&binaryRec)
		if err != nil {
			log.Print(err)
		}
		err = hnd.servBinary.BinaryRec(ctx, binaryRec)
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
func (hnd *Handlers) CardRec(ctx context.Context, cfg models.RabbitmqSrv) func(ctx *rmq.DeliveryContext) {
	return func(ctx *rmq.DeliveryContext) {
		// process delivery
		cardRec := models.CardRecord{}
		err := ctx.BindJSON(&cardRec)
		if err != nil {
			log.Print(err)
		}
		err = hnd.servCard.CardRec(ctx, cardRec)
		if err != nil {
			return
		}
		err = ctx.Delivery.Ack(true)
		if err != nil {
			log.Print("received Ask error:", settings.ColorRed, err, settings.ColorReset)
		}
	}
}
