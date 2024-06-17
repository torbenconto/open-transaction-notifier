package main

import "time"

type Transaction struct {
	Symbol        string
	Owner         string
	Shares        int
	PricePerShare float64
	Type          string // Acquisition or Disposition (A or D)
	Date          time.Time
}
