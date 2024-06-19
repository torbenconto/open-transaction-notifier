package main

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var processedTransactions = make(map[string]bool)
var isFirstRun = true

func main() {
	fmt.Println(os.Environ())
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

	for _, entry := range feed.Entries {
		link := entry.Link.Href
		shortURL := shortenURL(link)

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

		// Get index of transaction
		req, err := http.NewRequest(http.MethodGet, shortURL, http.NoBody)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
		req.Header.Set("User-Agent", UA)

		resp, err := SecHttpClient.Do(req)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
		defer resp.Body.Close()

		// Create a goquery document from the HTTP response
		document, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal("Error loading HTTP response body. ", err)
		}

		// Find and request all links with hrefs ending with ".xml"
		document.Find("tbody a").Each(func(index int, element *goquery.Selection) {
			href, exists := element.Attr("href")
			if exists && strings.HasSuffix(href, ".xml") {
				// Send a request to the .xml URL
				xmlReq, err := http.NewRequest(http.MethodGet, "https://sec.gov"+href, nil)
				if err != nil {
					log.Fatalf("error: %v\n", err)
				}
				xmlReq.Header.Set("User-Agent", UA)

				xmlResp, err := SecHttpClient.Do(xmlReq)
				if err != nil {
					log.Fatalf("error: %v\n", err)
				}
				defer xmlResp.Body.Close()

				// Parse the XML response into an OwnershipDocument
				var doc OwnershipDocument
				decoder := xml.NewDecoder(xmlResp.Body)
				decoder.CharsetReader = charset.NewReaderLabel
				err = decoder.Decode(&doc)
				if err != nil {
					log.Fatalf("error: %v\n", err)
				}

				// Call handleNewDocument with the parsed OwnershipDocument
				handleNewDocument(doc)
			}
		})
	}

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
