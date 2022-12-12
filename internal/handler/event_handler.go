package handler

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/goodsign/monday"
	"github.com/olebedev/when"
	"github.com/procyon-projects/chrono"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type EventHandler struct {
	nlDTParser *when.Parser
	scheduler  chrono.TaskScheduler
	waClient   *whatsmeow.Client
}

func NewEventHandler(w *when.Parser, scheduler chrono.TaskScheduler) *EventHandler {
	return &EventHandler{
		nlDTParser: w,
		scheduler:  scheduler,
	}
}

func (h *EventHandler) SetWAClient(c *whatsmeow.Client) {
	h.waClient = c
}

func (h *EventHandler) Func(evt interface{}) {
	switch evtType := evt.(type) {
	case *events.Message:
		message := parseConversation(evtType.Message.GetConversation())
		switch message.Action {
		case "remember me", "me lembra de":
			r, err := h.nlDTParser.Parse(message.Body, time.Now())
			if err != nil {
				log.Fatal(err)
			}

			if r == nil {
				h.waClient.SendMessage(context.TODO(), evtType.Info.Sender, "", &waProto.Message{
					Conversation: proto.String("Desculpe, não entendi direito. Eu sou facilmente confundido. Talvez tente as palavras em uma ordem diferente. Isso geralmente funciona: me lembra de [o que] [quando]"),
				})
			} else {
				todo := strings.ReplaceAll(message.Body, message.Body[r.Index:r.Index+len(r.Text)], "")
				todo = strings.TrimSpace(todo)

				fmt.Println(r.Time)

				_, err = h.scheduler.Schedule(func(ctx context.Context) {
					h.waClient.SendMessage(context.TODO(), evtType.Info.Sender, "", &waProto.Message{
						Conversation: proto.String(todo),
					})
				}, chrono.WithTime(r.Time))

				if err == nil {
					h.waClient.SendMessage(context.TODO(), evtType.Info.Sender, "", &waProto.Message{
						Conversation: proto.String(fmt.Sprintf("Vou te lembrar de %s as %s", todo, monday.Format(r.Time, "15:04 de Monday dia 02 de January", monday.LocalePtBR))),
					})
					log.Print("Task has been scheduled successfully.")
				}
			}
		default:
			log.Printf("unhandled action %s", message.Action)
		}
	}
}

type eventMessage struct {
	Action string
	Body   string
}

func parseConversation(conversation string) *eventMessage {
	message := &eventMessage{}

	re := regexp.MustCompile(`^(?P<action>remember me|me lembra de) ?(?P<body>.*)?$`)
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

	message.Action = result["action"]
	message.Body = result["body"]

	return message
}
