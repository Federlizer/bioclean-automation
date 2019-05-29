package main

import "time"

// Info holds the data that's to be written on the spreadsheet
type Info struct {
	Start       time.Time
	End         time.Time
	Description string
}
