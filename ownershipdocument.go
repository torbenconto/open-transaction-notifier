package main

import "encoding/xml"

// OwnershipDocument was generated 2024-06-16 16:32:35 by https://xml-to-go.github.io/ in Ukraine.
type OwnershipDocument struct {
	XMLName        xml.Name `xml:"ownershipDocument"`
	Text           string   `xml:",chardata"`
	SchemaVersion  string   `xml:"schemaVersion"`
	DocumentType   string   `xml:"documentType"`
	PeriodOfReport string   `xml:"periodOfReport"`
	Issuer         struct {
		Text                string `xml:",chardata"`
		IssuerCik           string `xml:"issuerCik"`
		IssuerName          string `xml:"issuerName"`
		IssuerTradingSymbol string `xml:"issuerTradingSymbol"`
	} `xml:"issuer"`
	ReportingOwner struct {
		Text             string `xml:",chardata"`
		ReportingOwnerId struct {
			Text         string `xml:",chardata"`
			RptOwnerCik  string `xml:"rptOwnerCik"`
			RptOwnerName string `xml:"rptOwnerName"`
		} `xml:"reportingOwnerId"`
		ReportingOwnerAddress struct {
			Text                     string `xml:",chardata"`
			RptOwnerStreet1          string `xml:"rptOwnerStreet1"`
			RptOwnerStreet2          string `xml:"rptOwnerStreet2"`
			RptOwnerCity             string `xml:"rptOwnerCity"`
			RptOwnerState            string `xml:"rptOwnerState"`
			RptOwnerZipCode          string `xml:"rptOwnerZipCode"`
			RptOwnerStateDescription string `xml:"rptOwnerStateDescription"`
		} `xml:"reportingOwnerAddress"`
		ReportingOwnerRelationship struct {
			Text              string `xml:",chardata"`
			IsDirector        string `xml:"isDirector"`
			IsOfficer         string `xml:"isOfficer"`
			IsTenPercentOwner string `xml:"isTenPercentOwner"`
			IsOther           string `xml:"isOther"`
			OfficerTitle      string `xml:"officerTitle"`
			OtherText         string `xml:"otherText"`
		} `xml:"reportingOwnerRelationship"`
	} `xml:"reportingOwner"`
	Aff10b5One         string `xml:"aff10b5One"`
	NonDerivativeTable struct {
		Text                     string `xml:",chardata"`
		NonDerivativeTransaction []struct {
			Text          string `xml:",chardata"`
			SecurityTitle struct {
				Text  string `xml:",chardata"`
				Value string `xml:"value"`
			} `xml:"securityTitle"`
			TransactionDate struct {
				Text  string `xml:",chardata"`
				Value string `xml:"value"`
			} `xml:"transactionDate"`
			TransactionCoding struct {
				Text                string `xml:",chardata"`
				TransactionFormType string `xml:"transactionFormType"`
				TransactionCode     string `xml:"transactionCode"`
				EquitySwapInvolved  string `xml:"equitySwapInvolved"`
				FootnoteId          struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"footnoteId"`
			} `xml:"transactionCoding"`
			TransactionAmounts struct {
				Text              string `xml:",chardata"`
				TransactionShares struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"transactionShares"`
				TransactionPricePerShare struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"transactionPricePerShare"`
				TransactionAcquiredDisposedCode struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"transactionAcquiredDisposedCode"`
			} `xml:"transactionAmounts"`
			PostTransactionAmounts struct {
				Text                            string `xml:",chardata"`
				SharesOwnedFollowingTransaction struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"sharesOwnedFollowingTransaction"`
			} `xml:"postTransactionAmounts"`
			OwnershipNature struct {
				Text                      string `xml:",chardata"`
				DirectOrIndirectOwnership struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"directOrIndirectOwnership"`
			} `xml:"ownershipNature"`
		} `xml:"nonDerivativeTransaction"`
	} `xml:"nonDerivativeTable"`
	DerivativeTable struct {
		Text                  string `xml:",chardata"`
		DerivativeTransaction struct {
			Text          string `xml:",chardata"`
			SecurityTitle struct {
				Text  string `xml:",chardata"`
				Value string `xml:"value"`
			} `xml:"securityTitle"`
			ConversionOrExercisePrice struct {
				Text  string `xml:",chardata"`
				Value string `xml:"value"`
			} `xml:"conversionOrExercisePrice"`
			TransactionDate struct {
				Text  string `xml:",chardata"`
				Value string `xml:"value"`
			} `xml:"transactionDate"`
			TransactionCoding struct {
				Text                string `xml:",chardata"`
				TransactionFormType string `xml:"transactionFormType"`
				TransactionCode     string `xml:"transactionCode"`
				EquitySwapInvolved  string `xml:"equitySwapInvolved"`
			} `xml:"transactionCoding"`
			TransactionAmounts struct {
				Text              string `xml:",chardata"`
				TransactionShares struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"transactionShares"`
				TransactionPricePerShare struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"transactionPricePerShare"`
				TransactionAcquiredDisposedCode struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"transactionAcquiredDisposedCode"`
			} `xml:"transactionAmounts"`
			ExerciseDate struct {
				Text       string `xml:",chardata"`
				FootnoteId struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"footnoteId"`
			} `xml:"exerciseDate"`
			ExpirationDate struct {
				Text  string `xml:",chardata"`
				Value string `xml:"value"`
			} `xml:"expirationDate"`
			UnderlyingSecurity struct {
				Text                    string `xml:",chardata"`
				UnderlyingSecurityTitle struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"underlyingSecurityTitle"`
				UnderlyingSecurityShares struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"underlyingSecurityShares"`
			} `xml:"underlyingSecurity"`
			PostTransactionAmounts struct {
				Text                            string `xml:",chardata"`
				SharesOwnedFollowingTransaction struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"sharesOwnedFollowingTransaction"`
			} `xml:"postTransactionAmounts"`
			OwnershipNature struct {
				Text                      string `xml:",chardata"`
				DirectOrIndirectOwnership struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value"`
				} `xml:"directOrIndirectOwnership"`
			} `xml:"ownershipNature"`
		} `xml:"derivativeTransaction"`
	} `xml:"derivativeTable"`
	Footnotes struct {
		Text     string `xml:",chardata"`
		Footnote struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"footnote"`
	} `xml:"footnotes"`
	OwnerSignature struct {
		Text          string `xml:",chardata"`
		SignatureName string `xml:"signatureName"`
		SignatureDate string `xml:"signatureDate"`
	} `xml:"ownerSignature"`
}
