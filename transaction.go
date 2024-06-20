package main

import "time"

type Relationship struct {
	IsDirector        bool
	IsOfficer         bool
	IsTenPercentOwner bool
	IsOther           bool
	OtherText         string
	OfficerTitle      string
}

type Transaction struct {
	Symbol        string
	Issuer        string
	Owner         string
	Relationship  Relationship
	Shares        float64
	PricePerShare float64
	Type          string // Acquisition or Disposition (A or D)
	Date          time.Time
}

// Transaction codes map
var transactionCodes = map[string]string{
	"A": "Grant, award, or other acquisition (such as under a bonus or incentive plan)",
	"C": "Conversion of derivative security (such as conversion of a convertible security)",
	"D": "Disposition to the issuer of issuer equity securities (e.g., a stock buyback)",
	"F": "Payment of exercise price or tax liability using portion of securities received from the exercise of a stock option",
	"G": "Gift of securities by or to the insider",
	"I": "Discretionary transaction, which is an order to the broker to buy or sell that is subject to broker’s discretion (used only for trust holdings)",
	"J": "Other acquisition or disposition (transaction type not covered by other codes)",
	"K": "Equity swap or similar derivative instrument",
	"L": "Small acquisition (a transaction not required to be reported during the period but voluntarily reported)",
	"M": "Exercise or conversion of derivative security (a call option, a warrant, or a convertible security)",
	"P": "Purchase of securities",
	"S": "Sale of securities",
	"U": "Disposition due to a tender of shares in a tender offer",
	"V": "Transaction voluntarily reported earlier than required",
	"W": "Acquisition or disposition by will or the laws of descent and distribution",
	"X": "Exercise of in-the-money or at-the-money option received in connection with a merger or acquisition",
}
