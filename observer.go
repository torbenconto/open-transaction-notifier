package main

import (
	"bytes"
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Observer struct {
	links        []string
	processed    map[string]bool
	firstRun     bool
	transactions []Transaction
	config       *Config
}

func NewObserver(config *Config) *Observer {
	return &Observer{
		config: config,
	}
}

func (o *Observer) Start() {
	log.Println("Scanning SEC filings...", o.firstRun)
	o.Fetch()

	transactions := make([]Transaction, 0)

	for _, link := range o.links {
		split := strings.Split(link, "/")

		// Remove the last path (html index)
		indexUrl := strings.Join(split[:len(split)-1], "/")

		// Check if the link is already processed
		if _, exists := o.processed[indexUrl]; exists {
			log.Printf("skipping: %s\n", indexUrl)
			continue
		}

		// Mark the link as processed
		o.processed[indexUrl] = true

		if !o.firstRun {
			log.Printf("Processing: %s\n", indexUrl)
			// Get index of transaction
			req, err := http.NewRequest(http.MethodGet, indexUrl, http.NoBody)
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

			document.Find("td a").EachWithBreak(func(index int, element *goquery.Selection) bool {
				href, exists := element.Attr("href")

				if exists && strings.HasSuffix(href, ".xml") {
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

					transactions = append(transactions, doc.ExtractTransactions()...)

					return false
				}
				return true
			})
		}
	}

	if !o.firstRun {
		notifier := NewNotifier(o.config)
		for _, t := range transactions {
			notifier.Notify(t)
		}
	}

	o.firstRun = false
}

func (o *Observer) Fetch() {
	params := url.Values{
		"action": []string{"getcurrent"},
		"type":   []string{"4"},
		"count":  []string{"50"},
		"owner":  []string{"only"},
		"output": []string{"atom"},
	}

	encodedUrl := baseUrl + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, encodedUrl, http.NoBody)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", UA)

	res, err := SecHttpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var feed Feed
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&feed)
	if err != nil {
		log.Fatalf("Failed to decode response body: %v", err)
	}

	uniqueLinks := make(map[string]string)
	for _, link := range feed.GetLinks() {
		u, err := url.Parse(link)
		if err != nil {
			log.Printf("Failed to parse URL: %v", err)
			continue
		}
		parts := strings.Split(path.Clean(u.Path), "/")
		if len(parts) < 2 {
			continue
		}
		lastTwoParts := strings.Join(parts[len(parts)-2:], "/")
		uniqueLinks[lastTwoParts] = link
	}

	links := make([]string, 0, len(uniqueLinks))
	for _, link := range uniqueLinks {
		links = append(links, link)
	}

	o.links = links
}

func (o *Observer) Init() {
	o.processed = make(map[string]bool)
	o.firstRun = true
}
