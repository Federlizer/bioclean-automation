package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// Constants defined by BioClean on their spreadsheets
const (
	worksheetName = "Timeseddel"

	dayOfTheWeekColumn = "B"
	startTimeColumn    = "C"
	endTimeColumn      = "D"
	hoursColumn        = "E"
	descriptionColumn  = "F"

	endOfWeekCellValue = "Total"
	endOfWeekSkipCount = 4

	beginningRow = 11

	path = "/root/Dropbox/Nikola Velichkov/"
)

// info holds the data that has to be written in the spreadsheets
// Calculate start and end values here https://www.myonlinetraininghub.com/excel-date-and-time
type info struct {
	start       float64
	end         float64
	description string
}

// Current default data to be written in the spreadsheets
var (
	normalInfo = info{
		start:       0.708333333,
		end:         0.8125,
		description: "Cleaning first two floors of stairs in 14 A, picking up trash from 14, cleaning main stairway in 46.",
	}
	altInfo = info{
		start:       0.708333333,
		end:         0.833333333,
		description: "Cleaning all stairs in 14, picking up trash in 14.",
	}
)

func main() {
	// We need the beginning of the day, otherwise the row calculation get's thrown off because of the time part of the struct
	t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

	// If it's a weekend day, skip it. We don't work on weekends.
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		fmt.Println("Not a working day, exiting.")
		return
	}

	f, err := excelize.OpenFile(path + getFileName(t))
	if err != nil {
		panic(err)
	}

	row := findRow(f, t)

	// Depending on the day, write different data to the spreadsheet
	if t.Weekday() == time.Wednesday || t.Weekday() == time.Friday {
		writeWorkInfo(f, row, altInfo)
	} else {
		writeWorkInfo(f, row, normalInfo)
	}
}

// getFileName calculates what the filename of the spreatsheet should be according to the date
// of which the program is executed.
func getFileName(now time.Time) string {
	var value strings.Builder

	if now.Month() < 10 {
		value.WriteString("0" + strconv.Itoa(int(now.Month())))
	} else {
		value.WriteString(strconv.Itoa(int(now.Month())))
	}

	value.WriteString(" Salary - ")
	value.WriteString(now.Month().String() + " ")
	value.WriteString(strconv.Itoa(now.Year()))
	value.WriteString(".xlsx")

	return value.String()
}

// findRow finds correct row corresponding to the date that has been passed as an argument.
// returns the row number that the program can use to write data on.
func findRow(f *excelize.File, now time.Time) int {
	// startTime is the date at the top of the spreadsheet file (usually 15th of the month)
	var startTime time.Time
	var hours, _ = time.ParseDuration("24h")
	var row = beginningRow

	// Since the spreadsheet files are from the 15th of the previous month,
	// until the 14th of the current month, calculate startTime accordingly
	if now.Day() < 15 {
		startTime = time.Date(now.Year(), now.Month()-1, 15, 0, 0, 0, 0, time.UTC)
	} else {
		startTime = time.Date(now.Year(), now.Month(), 15, 0, 0, 0, 0, time.UTC)
	}

	// Increment the row relatively to the days skipped
	for startTime.Before(now) {
		startTime = startTime.Add(hours)
		dayOfTheWeek, err := f.GetCellValue(worksheetName, dayOfTheWeekColumn+strconv.Itoa(row+1))
		if err != nil {
			panic(err)
		}

		// If it's the end of the week, skip more cells.
		// Distance between the weeks is defined on the top of the file.
		// The distance is defined by BioClean
		if dayOfTheWeek == endOfWeekCellValue {
			row += endOfWeekSkipCount
		} else {
			row++
		}
	}

	return row
}

// writeWorkInfo writes the work data to the spreadsheet, given the row of the spreadsheet.
// The columns are already defined at the top of the file.
func writeWorkInfo(f *excelize.File, row int, i info) {
	rowString := strconv.Itoa(row)

	err := f.SetCellValue(worksheetName, startTimeColumn+rowString, i.start)
	if err != nil {
		panic(err)
	}

	err = f.SetCellValue(worksheetName, endTimeColumn+rowString, i.end)
	if err != nil {
		panic(err)
	}

	err = f.SetCellValue(worksheetName, descriptionColumn+rowString, i.description)
	if err != nil {
		panic(err)
	}

	fmt.Println("Data written successfully. Saving file...")

	err = f.Save()
	if err != nil {
		panic(err)
	}

	fmt.Println("File saved.")
}
