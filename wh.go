package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields"`
}

type WebhookMessage struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

func sendEmbedToWebhook(webhookURL string, message WebhookMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
