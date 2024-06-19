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
	var txtype string
	switch tx.Type {
	case "A":
		txtype = "Purchase"
	case "D":
		txtype = "Sale"
	}
	if tx.PricePerShare == 0.000000 {
		txtype = "Option Exercise"
	}
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
						Name:  "Type",
						Value: txtype,
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
