package main

import (
	"encoding/json"
	"os"
)

// Config holds configuration values described in a json file
type Config struct {
	WorksheetName string `json:"worksheetName"`

	DayOfTheWeekColumn string `json:"dayOfTheWeekColumn"`
	StartTimeColumn    string `json:"startTimeColumn"`
	EndTimeColumn      string `json:"endTimeColumn"`
	DescriptionColumn  string `json:"descriptionColumn"`

	EndOfWeekCellValue string `json:"endOfWeekCellValue"`
	EndOfWeekSkipCount int    `json:"endOfWeekSkipCount"`
	BeginningRow       int    `json:"beginningRow"`

	NormalInfo    Info `json:"normalInfo"`
	AlternateInfo Info `json:"alternateInfo"`

	PathToSpreadsheets string `json:"pathToSpreadsheets"`
}

func parseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	config := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return &config, err
	}

	return &config, nil
}
