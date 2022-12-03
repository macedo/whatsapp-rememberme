package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/lib/pq"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

func init() {
	sqlstore.PostgresArrayWrapper = pq.Array
}

func HelloWorldHandler(evt interface{}) {
	switch evtType := evt.(type) {
	case *events.Message:
		message := parseConversation(evtType.Message.GetConversation())
		switch message.Program {
		case "remember me":
			fmt.Println("remember me program")
		default:
			client.Log.Debugf("unhandled program %s", message.Program)
		}
		fmt.Printf("%v\n", message)
	}
}

type Message struct {
	Program string
	Todo    string
	When    string
}

func parseConversation(conversation string) *Message {
	var message = &Message{}
	re := regexp.MustCompile(`^(?P<program>remember me) ?(?P<todo>[0-9a-z ]+)? ?(?P<when>@ [0-9a-z/-:]+)?$`)
	result := make(map[string]string)
	matches := re.FindStringSubmatch(conversation)

	if len(matches) == 0 {
		return message
	}

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	message.Program = result["program"]
	message.Todo = result["todo"]
	message.When = result["when"]

	return message
}

func main() {
	dbLog := waLog.Stdout("DATABASE", "INFO", true)
	container, err := sqlstore.New("postgres", "postgres://postgres:0519@localhost:5432/whatsapp-rememberme", dbLog)
	check_error(err)

	deviceStore, err := container.GetFirstDevice()
	check_error(err)

	clientLog := waLog.Stdout("CLIENT", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	client.AddEventHandler(HelloWorldHandler)

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
