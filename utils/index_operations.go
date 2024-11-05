package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mike-mo-coposit/kafka-event-engine/opensearch"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func CreateDocument(index string, id string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create the IndexRequest with opensearchapi
	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(jsonData),
	}

	res, err := req.Do(context.Background(), opensearch.Client)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from OpenSearch: %s", res.String())
	}

	log.Printf("Document created: ID=%s, Index=%s", id, index)
	return nil
}

func UpdateDocument(index string, id string, data interface{}) error {
	jsonData, err := json.Marshal(map[string]interface{}{
		"doc": data,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Use opensearchapi.UpdateRequest instead for partial updates
	req := opensearchapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(jsonData),
	}

	res, err := req.Do(context.Background(), opensearch.Client)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from OpenSearch: %s", res.String())
	}

	log.Printf("Document updated: ID=%s, Index=%s", id, index)
	return nil
}
