package main

import (
	"crypto/tls"
	"net/http"
)

const baseUrl = "https://www.sec.gov/cgi-bin/browse-edgar"

const UA = "OpenTransactionNotifier torbenmconto@gmail.com"

var SecHttpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}
