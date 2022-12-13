package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/macedo/whatsapp-rememberme/internal/handler"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Service interface {
	Start() <-chan error
	Stop()
}

func New(db *sql.DB, evtHandler *handler.EventHandler) Service {
	name := "WHATSAPP"
	logLevel := "INFO"

	return &WhatsApp{
		db:         db,
		evtHandler: evtHandler,
		log:        waLog.Stdout(name, logLevel, true),
		name:       name,
	}
}

type WhatsApp struct {
	db         *sql.DB
	name       string
	evtHandler *handler.EventHandler
	log        waLog.Logger
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
	log.Printf("service %s started", s.name)

	container := sqlstore.NewWithDB(s.db, "sqlite3", s.log.Sub("DATABASE"))

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return err
	}

	cli := whatsmeow.NewClient(deviceStore, s.log.Sub("CLIENT"))
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
	log.Printf("service %s stopped", s.name)
	return nil
}
