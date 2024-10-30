package utils

import (
	"encoding/json"
	"fmt"
)

// EventData represents the structure of the event.
type EventData struct {
	Index  string      `json:"index"`
	ID     string      `json:"id"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func ParseEvent(event string) (*EventData, error) {
	var eventData EventData
	if err := json.Unmarshal([]byte(event), &eventData); err != nil {
		return nil, fmt.Errorf("error parsing event data: %w", err)
	}
	return &eventData, nil
}
