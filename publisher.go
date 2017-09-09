package evently

import (
	"encoding/json"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// EventPublisher is an application event publisher.
type EventPublisher struct {
	Config ServerConfig
	Server *machinery.Server
	Errors []error
}

// NewEventPublisher creates a new application event publisher.
func NewEventPublisher(settings map[string]string) *EventPublisher {
	publisher := &EventPublisher{}
	var err error

	// Attempting to load event publishers configuration
	publisher.Config = GetConfiguration(settings)

	// Attempting to start the event server
	publisher.Server, err = machinery.NewServer(publisher.Config)
	publisher.processError(err)

	return publisher
}

// PublishEvent is a helper method that triggers a user event
func (publisher *EventPublisher) PublishEvent(name string, data interface{}) *EventPublisher {
	// JSON encoding the event data
	jsonData, err := json.Marshal(data)
	publisher.processError(err)

	// Preparing event data
	event := &tasks.Signature{
		RoutingKey: publisher.Config.DefaultQueue,
		Name:       name,
		Args:       []tasks.Arg{{Type: "string", Value: string(jsonData)}},
	}

	// Attempting to trigger the event
	_, err = publisher.Server.SendTask(event)
	publisher.processError(err)

	return publisher
}

// processError processes an event publisher error
func (publisher *EventPublisher) processError(err error) {
	if err != nil {
		publisher.Errors = append(publisher.Errors, err)
	}
}

// IsOK determins the ok status of the event publisher
func (publisher *EventPublisher) IsOK() bool {
	return (len(publisher.Errors) == 0)
}
