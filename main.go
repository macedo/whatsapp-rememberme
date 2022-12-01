package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lib/pq"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

func init() {
	sqlstore.PostgresArrayWrapper = pq.Array
}

func main() {
	dbLog := waLog.Stdout("DATABASE", "INFO", true)
	container, err := sqlstore.New("postgres", "postgres://postgres:0519@localhost:5432/whatsapp-rememberme", dbLog)
	check_error(err)

	deviceStore, err := container.GetFirstDevice()
	check_error(err)

	clientLog := waLog.Stdout("CLIENT", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		check_error(err)

		for event := range qrChan {
			if event.Event == "code" {
				qrterminal.GenerateHalfBlock(event.Code, qrterminal.L, os.Stdout)
				fmt.Println("QR code: ", event.Code)
			} else {
				fmt.Println("Login event: ", event.Event)
			}
		}
	} else {
		err = client.Connect()
		check_error(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Disconnect()
}

func check_error(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
