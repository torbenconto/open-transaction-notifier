package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func Notify(method, info string, tx Transaction) {
	switch method {
	case "discord_webhook":
		NotifyDiscordWebhook(info, tx)
	}
}

func NotifyDiscordWebhook(link string, tx Transaction) {
	var txtype string
	var color int
	var relationship string

	switch tx.Type {
	case "P":
		color = 0x00FF00 // Green
	case "S":
		color = 0xFF0000 // Red
	}
	if tx.PricePerShare == 0.000000 {
		color = 0x0000FF
	}

	txtype = transactionCodes[tx.Type]

	// Get color from config
	configColor := config.Discord.Embed.Color
	if configColor != "default" {
		// Parse hex color code to int
		hexColor, err := strconv.ParseInt(strings.TrimPrefix(configColor, "#"), 16, 64)
		if err == nil {
			color = int(hexColor)
		}
	}

	switch {
	case tx.Relationship.IsDirector:
		relationship = "Director"
	case tx.Relationship.IsOfficer:
		relationship = tx.Relationship.OfficerTitle
	case tx.Relationship.IsTenPercentOwner:
		relationship = "10% Owner"
	case tx.Relationship.IsOther:
		relationship = tx.Relationship.OtherText
	}

	// Create embed fields based on config
	var fields []EmbedField
	if config.Discord.Embed.Fields.Ticker {
		fields = append(fields, EmbedField{Name: "Ticker", Value: tx.Symbol})
	}
	if config.Discord.Embed.Fields.Type {
		fields = append(fields, EmbedField{Name: "Type", Value: txtype})
	}
	if config.Discord.Embed.Fields.PricePerShare {
		fields = append(fields, EmbedField{Name: "Price Per Share", Value: fmt.Sprintf("%f", tx.PricePerShare)})
	}
	if config.Discord.Embed.Fields.Shares {
		fields = append(fields, EmbedField{Name: "Shares", Value: fmt.Sprintf("%v", tx.Shares)})
	}
	if config.Discord.Embed.Fields.Owner {
		fields = append(fields, EmbedField{Name: "Owner", Value: fmt.Sprintf("%s", tx.Owner)})
	}
	if config.Discord.Embed.Fields.Relationship {
		fields = append(fields, EmbedField{Name: "Relationship", Value: relationship})
	}
	if config.Discord.Embed.Fields.Date {
		fields = append(fields, EmbedField{Name: "Date", Value: fmt.Sprintf("%v", tx.Date)})
	}

	err := sendEmbedToWebhook(link, WebhookMessage{
		Content: "Open Transaction Notifier - New Transaction",
		Embeds: []Embed{
			{
				Title:  "New Transaction",
				Fields: fields,
				Color:  color,
			},
		},
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// Sleep 1/4 second to avoid rate limiting
	time.Sleep(250 * time.Millisecond)
}
