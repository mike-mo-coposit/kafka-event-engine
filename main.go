package main

import (
	"context"
	"fmt"

	"shareutils/opensearch"
	"shareutils/utils"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event string) error {
	// Parse the event
	eventData, err := utils.ParseEvent(event)
	if err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	// Initialize the OpenSearch client
	opensearch.InitClient()

	// Handle create or update based on event action
	switch eventData.Action {
	case "create":
		return utils.CreateDocument(eventData.Index, eventData.ID, eventData.Data)
	case "update":
		return utils.UpdateDocument(eventData.Index, eventData.ID, eventData.Data)
	default:
		return fmt.Errorf("invalid action: %s", eventData.Action)
	}
}

func main() {
	// Start the Lambda function
	lambda.Start(Handler)
}
