package main

import (
	"fmt"
	"time"

	"github.com/Federlizer/timeseddel"

	"github.com/Federlizer/timeseddel/bioclean"
)

var (
	normalWI = timeseddel.WorkInfo{
		Start: time.Date(0, time.January, 1, 17, 0, 0, 0, time.UTC),
		End:   time.Date(0, time.January, 1, 20, 0, 0, 0, time.UTC),
		Info:  "Picking up trash from 14, cleaning spots in the stairs in 14",
	}
)

func main() {
	date := time.Now()

	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		fmt.Println("Not a working day, skipping...")
		return
	}

	fmt.Println("Creating manager")
	manager, err := bioclean.NewBiocleanManager("/home/federlizer/Dropbox/Nikola Velichkov/")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Writing work info %v for date %v", normalWI, date)
	err = manager.WriteWorkInfo(date, normalWI)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")
}
