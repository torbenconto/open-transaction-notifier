package main

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var config Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
}

func main() {
	duration, err := time.ParseDuration(config.OpenTransactionNotifier.TimeInterval)
	if err != nil {
		log.Fatalf("Error while parsing time duration %s", err)
	}

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkForNewTransactions()
		}
	}
}

var processedTransactions = make(map[string]bool)
var isFirstRun = true

func checkForNewTransactions() {
	log.Println("Checking for new transactions...")
	feed, err := GetTransactions()
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	newTransactions := make([]Transaction, 0) // Collect new transactions to notify

	for _, entry := range feed.Entries {
		link := entry.Link.Href
		shortURL := shortenURL(link)

		// Check if the transaction has already been processed
		if _, exists := processedTransactions[shortURL]; exists {
			continue
		}

		// Mark the transaction as processed
		processedTransactions[shortURL] = true

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

		// Find and request the first link with href ending with ".xml"
		document.Find("tbody a").EachWithBreak(func(index int, element *goquery.Selection) bool {
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

				// Collect new transactions
				newTransactions = append(newTransactions, extractTransactions(doc)...)

				// Break the loop after processing the first link
				return false
			}
			return true
		})
	}

	if !isFirstRun {
		notifyTransactions(newTransactions)
	} else {
		log.Println("Skipping notifications for the first run.")
	}

	isFirstRun = false // Update flag after first run
}

func extractTransactions(doc OwnershipDocument) []Transaction {
	var transactions []Transaction
	// Handle non-derivative transactions
	for _, transaction := range doc.NonDerivativeTable.NonDerivativeTransaction {
		date, err := time.Parse("2006-01-02", transaction.TransactionDate.Value)
		if err != nil {
			log.Fatalf("error parsing date: %v", err)
		}

		shares, err := strconv.ParseFloat(transaction.TransactionAmounts.TransactionShares.Value, 64)
		if err != nil {
			log.Fatalf("error parsing shares: %v", err)
		}

		pricePerShare, err := strconv.ParseFloat(transaction.TransactionAmounts.TransactionPricePerShare.Value, 64)
		if err != nil {
			log.Fatalf("error parsing price per share: %v", err)
		}

		var isDirector, isOfficer, isTenPercentOwner, isOther bool

		if doc.ReportingOwner.ReportingOwnerRelationship.IsDirector == "" {
			isDirector = false
		} else {
			isDirector = true
		}

		if doc.ReportingOwner.ReportingOwnerRelationship.IsOfficer == "" {
			isOfficer = false
		} else {
			isOfficer = true
		}

		if doc.ReportingOwner.ReportingOwnerRelationship.IsTenPercentOwner == "" {
			isTenPercentOwner = false
		} else {
			isTenPercentOwner = true
		}

		if doc.ReportingOwner.ReportingOwnerRelationship.IsOther == "" {
			isOther = false
		} else {
			isOther = true
		}
		//TODO Possible to move relationship logic up to here and determine relationship title here instead of passing down to notify
		tx := Transaction{
			Symbol: doc.Issuer.IssuerTradingSymbol,
			Issuer: doc.Issuer.IssuerName,
			Owner:  doc.ReportingOwner.ReportingOwnerId.RptOwnerName,
			Relationship: Relationship{
				IsDirector:        isDirector,
				IsOfficer:         isOfficer,
				IsTenPercentOwner: isTenPercentOwner,
				IsOther:           isOther,
				OtherText:         doc.ReportingOwner.ReportingOwnerRelationship.OtherText,
				OfficerTitle:      doc.ReportingOwner.ReportingOwnerRelationship.OfficerTitle,
			},
			Date:          date,
			Shares:        shares,
			PricePerShare: pricePerShare,
			Type:          transaction.TransactionCoding.TransactionCode,
		}

		transactions = append(transactions, tx)
	}
	return transactions
}

func notifyTransactions(transactions []Transaction) {
	for _, transaction := range transactions {
		method := config.OpenTransactionNotifier.NotificationMethod
		info := config.OpenTransactionNotifier.NotificationInfo
		if method == "" || info == "" {
			log.Fatalf("missing notification method or info")
		}
		Notify(method, info, transaction)
	}
}

func shortenURL(url string) string {
	lastSlashIndex := strings.LastIndex(url, "/")
	return url[:lastSlashIndex]
}
