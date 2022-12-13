package whatsapp

import (
	"context"
	"fmt"
	"os"

	"github.com/macedo/whatsapp-rememberme/internal/handler"
	"github.com/macedo/whatsapp-rememberme/internal/logadapter"
	"github.com/macedo/whatsapp-rememberme/internal/store/sqlstore"
	"github.com/mdp/qrterminal"
	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
)

type Service interface {
	Start() <-chan error
	Stop()
}

func New(container *sqlstore.Container, evtHandler *handler.EventHandler) Service {
	name := "whatsapp"

	return &WhatsApp{
		container:  container,
		evtHandler: evtHandler,
		name:       name,
	}
}

type WhatsApp struct {
	container  *sqlstore.Container
	name       string
	evtHandler *handler.EventHandler
	ctx        context.Context
	cancel     context.CancelFunc
}

func (s *WhatsApp) Start() <-chan error {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	errc := make(chan error)
	go func() {
		defer close(errc)
		if err := s.run(); err != nil {
			errc <- err
		}
	}()

	return errc
}

func (s *WhatsApp) Stop() {
	s.cancel()
}

func (s *WhatsApp) run() error {
	log.Info().Str("service", s.name).Msg("service started")

	deviceStore, err := s.container.GetFirstDevice()
	if err != nil {
		return err
	}

	cliLog := log.With().Logger()
	cli := whatsmeow.NewClient(deviceStore, logadapter.WALogAdapter(cliLog))
	defer cli.Disconnect()

	s.evtHandler.SetWAClient(cli)

	cli.AddEventHandler(s.evtHandler.Func)

	if cli.Store.ID == nil {
		qrCh, _ := cli.GetQRChannel(s.ctx)
		err := cli.Connect()
		if err != nil {
			return err
		}
		for qrItem := range qrCh {
			if qrItem.Event == "code" {
				qrterminal.GenerateHalfBlock(qrItem.Code, qrterminal.L, os.Stdout)
				fmt.Println("qrcode: ", qrItem.Code)
			} else {
				fmt.Println("loggin event", qrItem.Event)
			}
		}
	} else {
		err := cli.Connect()
		if err != nil {
			return err
		}
	}

	<-s.ctx.Done()
	log.Info().Str("service", s.name).Msg("service stopped")
	return nil
}
