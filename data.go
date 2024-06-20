package main

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"net/url"
)

type Feed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Link Link `xml:"link"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

func GetTransactions() (Feed, error) {
	params := url.Values{
		"action": []string{"getcurrent"},
		"type":   []string{"4"},
		"count":  []string{"25"},
		"owner":  []string{"only"},
		"output": []string{"atom"},
	}

	encodedUrl := baseUrl + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, encodedUrl, http.NoBody)
	if err != nil {
		return Feed{}, err
	}
	req.Header.Set("User-Agent", UA)

	res, err := SecHttpClient.Do(req)
	if err != nil {
		return Feed{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Feed{}, err
	}

	var feed Feed
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&feed)
	if err != nil {
		return Feed{}, err
	}

	return feed, nil
}
