package main

import "time"

type Transaction struct {
	Symbol        string
	Owner         string
	Shares        float64
	PricePerShare float64
	Type          string // Acquisition or Disposition (A or D)
	Date          time.Time
}
