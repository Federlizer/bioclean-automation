package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	configPath    = "/etc/bioclean/default.json"
	configPathDev = "./config/default.json"
)

func main() {
	conf, err := parseConfig(configPathDev)
	if err != nil {
		panic(err)
	}

	// We need the beginning of the day, otherwise the row calculation get's thrown off because of the time part of the struct
	t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

	// If it's a weekend day, skip it. We don't work on weekends.
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		fmt.Println("Not a working day, exiting.")
		return
	}

	f, err := excelize.OpenFile(conf.PathToSpreadsheets + getFileName(t))
	if err != nil {
		panic(err)
	}

	row := findRow(f, t, conf)

	// Depending on the day, write different data to the spreadsheet
	if t.Weekday() == time.Wednesday || t.Weekday() == time.Friday {
		writeWorkInfo(f, row, conf, conf.AlternateInfo)
	} else {
		writeWorkInfo(f, row, conf, conf.NormalInfo)
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
func findRow(f *excelize.File, now time.Time, config *Config) int {
	// startTime is the date at the top of the spreadsheet file (usually 15th of the month)
	var startTime time.Time
	var hours, _ = time.ParseDuration("24h")
	var row = config.BeginningRow

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
		dayOfTheWeek, err := f.GetCellValue(config.WorksheetName, config.DayOfTheWeekColumn+strconv.Itoa(row+1))
		if err != nil {
			panic(err)
		}

		// If it's the end of the week, skip more cells.
		// Distance between the weeks is defined on the top of the file.
		// The distance is defined by BioClean
		if dayOfTheWeek == config.EndOfWeekCellValue {
			row += config.EndOfWeekSkipCount
		} else {
			row++
		}
	}

	return row
}

// writeWorkInfo writes the work data to the spreadsheet, given the row of the spreadsheet.
// The columns are already defined in the config file.
func writeWorkInfo(f *excelize.File, row int, config *Config, info Info) {
	var err error
	var rowString string

	rowString = strconv.Itoa(row)

	err = f.SetCellValue(
		config.WorksheetName,
		config.StartTimeColumn+rowString,
		calculateDate(info.Start),
	)
	if err != nil {
		panic(err)
	}

	err = f.SetCellValue(
		config.WorksheetName,
		config.EndTimeColumn+rowString,
		calculateDate(info.End),
	)
	if err != nil {
		panic(err)
	}

	err = f.SetCellValue(
		config.WorksheetName,
		config.DescriptionColumn+rowString,
		info.Description,
	)
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

// calculateDate calculates the float value of the time given as an argument.
// Find more info about the calculation at https://www.myonlinetraininghub.com/excel-date-and-time
func calculateDate(t time.Time) float64 {
	var hours, minutes, seconds float64

	hours = float64(t.Hour() / 24)
	minutes = float64(t.Minute() / 1440)
	seconds = float64(t.Second() / 86400)

	return hours + minutes + seconds
}
