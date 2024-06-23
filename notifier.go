package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Notifier struct {
	config *Config
}

func NewNotifier(config *Config) *Notifier {
	return &Notifier{
		config: config,
	}
}

func (n *Notifier) Notify(transaction Transaction) {
	switch n.config.OpenTransactionNotifier.NotificationMethod {
	case "discord_webhook":
		n.DiscordWebhook(transaction)
	}
}

func (n *Notifier) DiscordWebhook(transaction Transaction) {
	var txtype string
	var color int
	var relationship string

	switch transaction.Type {
	case "P":
		color = 0x00FF00 // Green
	case "S":
		color = 0xFF0000 // Red
	}

	fmt.Println(transaction.Symbol, transaction.Type)

	if n.config.OpenTransactionNotifier.OnlyReportSalesAndBuys {
		if transaction.Type != "P" && transaction.Type != "S" {
			return
		}
	}

	txtype = transactionCodes[transaction.Type]

	// Get color from config
	configColor := n.config.Discord.Embed.Color
	if configColor != "default" {
		// Parse hex color code to int
		hexColor, err := strconv.ParseInt(strings.TrimPrefix(configColor, "#"), 16, 64)
		if err == nil {
			color = int(hexColor)
		}
	}

	switch {
	case transaction.Relationship.IsDirector:
		relationship = "Director"
	case transaction.Relationship.IsOfficer:
		relationship = transaction.Relationship.OfficerTitle
	case transaction.Relationship.IsTenPercentOwner:
		relationship = "10% Owner"
	case transaction.Relationship.IsOther:
		relationship = transaction.Relationship.OtherText
	}

	decimalFormatString := fmt.Sprintf("%%.%df", conf.Discord.Embed.Decimals)

	// Create embed fields based on config
	var fields []EmbedField
	if n.config.Discord.Embed.Fields.Ticker {
		fields = append(fields, EmbedField{Name: "Ticker", Value: transaction.Symbol})
	}
	if n.config.Discord.Embed.Fields.Type {
		fields = append(fields, EmbedField{Name: "Type", Value: txtype})
	}
	if n.config.Discord.Embed.Fields.PricePerShare {
		fields = append(fields, EmbedField{Name: "Price Per Share", Value: fmt.Sprintf(decimalFormatString, transaction.PricePerShare)})
	}
	if n.config.Discord.Embed.Fields.Shares {
		fields = append(fields, EmbedField{Name: "Shares", Value: fmt.Sprintf(decimalFormatString, transaction.Shares)})
	}
	if n.config.Discord.Embed.Fields.Owner {
		fields = append(fields, EmbedField{Name: "Owner", Value: fmt.Sprintf("%s", transaction.Owner)})
	}
	if n.config.Discord.Embed.Fields.Relationship {
		fields = append(fields, EmbedField{Name: "Relationship", Value: relationship})
	}
	if n.config.Discord.Embed.Fields.Date {
		fields = append(fields, EmbedField{Name: "Date", Value: fmt.Sprintf("%v", transaction.Date)})
	}

	err := sendEmbedToWebhook(conf.OpenTransactionNotifier.NotificationInfo, WebhookMessage{
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
