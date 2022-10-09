package slack

import (
	"context"
	"fmt"

	"github.com/slack-go/slack/slackevents"

	"github.com/slack-go/slack/socketmode"

	slackLib "github.com/slack-go/slack"
	"github.com/vzxw/gobot/internal/pkg/logger"
	"github.com/vzxw/gobot/internal/pkg/message"
)

type slack struct {
	api    *slackLib.Client
	client *socketmode.Client
}

func New(appToken string, botToken string) *slack {
	api := slackLib.New(
		botToken, slackLib.OptionDebug(true),
		slackLib.OptionLog(logger.NewInfo("API")),
		slackLib.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(logger.NewInfo("SOCKET")),
	)

	return &slack{
		api:    api,
		client: client,
	}
}

func (s *slack) Listen(ctx context.Context) (<-chan message.Message, error) {
	_, _, err := s.api.ConnectRTMContext(ctx)
	if err != nil {
		return nil, err
	}

	result := make(chan message.Message, 10)
	go func() {
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				fmt.Println("Shutting down socketmode listener")
				return
			case event := <-s.client.Events:
				fmt.Println("Slack event", event.Type)
				// We have a new Events, let's type switch the event
				// Add more use cases here if you want to listen to other events.
				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						fmt.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					s.client.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					fmt.Println(eventsAPIEvent)
				}
			}
		}
	}()

	go func() {
		err := s.client.Run()
		if err != nil {
			result <- message.Message{
				Text: "",
				Err:  err,
			}
		}
	}()

	return result, nil
}
