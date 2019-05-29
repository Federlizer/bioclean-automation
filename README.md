# BioClean automation
This program automatically writes down the working hours, as per the BioClean spreadsheet format.
#### Usage
```
go get github.com/Federlizer/bioclean-automation
go install ~/go/src/github.com/Federlizer/bioclean-automation

bioclean-automation
```

A suggestion is to set up the program to fire once a day using cron jobs.

#### TODO
- [X] Create a config file to store defaults
- [ ] Make it possible to invoke another configuration path via flags
- [ ] Writing data for previous days
