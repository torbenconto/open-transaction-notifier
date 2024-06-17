package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var processedTransactions = make(map[string]bool)
var isFirstRun = true

func main() {
	ticker := time.NewTicker(10 * time.Second) // Adjust the interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkForNewTransactions()
		}
	}
}

func checkForNewTransactions() {
	log.Println("Checking for new transactions...")
	feed, err := GetTransactions()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	var wg sync.WaitGroup

	for _, entry := range feed.Entries {
		link := entry.Link.Href
		shortURL := shortenURL(link) + "/ownership.xml"

		// Check if the transaction has already been processed
		if _, exists := processedTransactions[shortURL]; exists && !isFirstRun {
			continue
		}

		// Mark the transaction as processed
		processedTransactions[shortURL] = true

		// If it's the first run, don't handle the transaction
		if isFirstRun {
			continue
		}

		wg.Add(1)
		go func(entry Entry) {
			defer wg.Done()

			fmt.Print("shortURL: ", shortURL, "\n")

			req, err := http.NewRequest(http.MethodGet, shortURL, http.NoBody)
			if err != nil {
				log.Fatalf("error: %v\n", err)
			}
			req.Header.Set("User-Agent", UA)

			res, err := (&SecHttpClient).Do(req)
			defer res.Body.Close()

			if err != nil {
				log.Fatalf("error: %v\n", err)
			}

			reader := bufio.NewReader(res.Body)

			var doc OwnershipDocument
			decoder := xml.NewDecoder(reader)
			decoder.CharsetReader = charset.NewReaderLabel
			err = decoder.Decode(&doc)
			if err != nil {
				log.Fatalf("error: %v\n", err)
			}

			// Call the function for each new doc
			handleNewDocument(doc)
		}(entry)
	}

	wg.Wait()

	isFirstRun = false
}

func handleNewDocument(doc OwnershipDocument) {
	var transactions []Transaction
	// Handle non-derivative transactions
	for _, transaction := range doc.NonDerivativeTable.NonDerivativeTransaction {
		date, err := time.Parse("2006-01-02", transaction.TransactionDate.Value)
		if err != nil {
			log.Fatalf("error parsing date: %v", err)
		}

		shares, err := strconv.Atoi(transaction.TransactionAmounts.TransactionShares.Value)
		if err != nil {
			log.Fatalf("error parsing shares: %v", err)
		}

		pricePerShare, err := strconv.ParseFloat(transaction.TransactionAmounts.TransactionPricePerShare.Value, 64)

		tx := Transaction{
			Symbol:        doc.Issuer.IssuerTradingSymbol,
			Owner:         doc.Issuer.IssuerName,
			Date:          date,
			Shares:        shares,
			PricePerShare: pricePerShare,
			Type:          transaction.TransactionAmounts.TransactionAcquiredDisposedCode.Value,
		}

		transactions = append(transactions, tx)
	}

	for _, transaction := range transactions {
		method := os.Getenv("NOTIFICATION_METHOD") // "discord_webhook" for now
		info := os.Getenv("NOTIFICATION_INFO")     // Discord webhook URL or email address (soon)
		if method == "" || info == "" {
			log.Fatalf("missing notification method or info")
		}
		Notify(method, info, transaction)
	}
	// Handle derivative transactions (WIP)
	/*
		for _, transaction := range doc.DerivativeTable.DerivativeTransaction {
			// Here, transaction is an individual derivative transaction
			// You can access its fields like transaction.SecurityTitle.Value, transaction.TransactionDate.Value, etc.
			fmt.Printf("Derivative Transaction: %+v\n", transaction)
		}
	*/
}

func shortenURL(url string) string {
	lastSlashIndex := strings.LastIndex(url, "/")
	return url[:lastSlashIndex]
}
