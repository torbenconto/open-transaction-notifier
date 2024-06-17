package main

import (
	"fmt"
	"log"
)

func Notify(method, info string, tx Transaction) {
	switch method {
	case "discord_webhook":
		NotifyDiscordWebhook(info, tx)
	}
}

func NotifyDiscordWebhook(link string, tx Transaction) {
	err := sendEmbedToWebhook(link, WebhookMessage{
		Content: "Open Transaction Notifier - New Transaction",
		Embeds: []Embed{
			{
				Title: "New Transaction",
				Fields: []EmbedField{
					{
						Name:  "Ticker",
						Value: tx.Symbol,
					},
					{
						Name:  "Price Per Share",
						Value: fmt.Sprintf("%f", tx.PricePerShare),
					},
					{
						Name:  "Shares",
						Value: fmt.Sprintf("%v", tx.Shares),
					},
					{
						Name:  "Owner",
						Value: fmt.Sprintf("%s", tx.Owner),
					},
					{
						Name:  "Date",
						Value: fmt.Sprintf("%v", tx.Date),
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
