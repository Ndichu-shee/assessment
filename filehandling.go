package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type CsvFields struct {
	ID               string
	UserID           string
	TimeReceived     time.Time
	TimeBegan        time.Time
	TimeFinished     time.Time
	DeletedBuilds    bool
	ExitBuildProcess int
	ImageSize        int
}

// NEW : instantiate a slice of CSVFields
var fields []CsvFields
func main() {
	// NEW read the data
	fields, err := readStats("stats.csv")
	checkError("Error in reading tis file", err)

	users := BuildRemoteService()
	fmt.Println("\nThese are the top 5 users:")
	fmt.Println(users)

	rate := SucccessRate()
	fmt.Printf("\nSuccess Rate: %.2f%%\n", rate)//round off the result to 2 decimal places

	inspectTime := time.Now() // update this to what should be passed

	executedbuilds := BuildTimeWindow(fields, inspectTime)

	fmt.Println(executedbuilds)
}

/* How many builds were executed in a time window ?
 BuildTimeWindow takes in slice of CsvFields, inspect parameter and returns an integer. Goes through the file checks  using After and Before methods to compare two times */
func BuildTimeWindow(fields []CsvFields, inspect time.Time) int {
	executedbuilds := 0

	for _, field := range fields {
		if inspect.After(field.TimeBegan) && inspect.Before(field.TimeFinished) {
			executedbuilds++
		}
	}

	return executedbuilds
}


/*  BuildRemoteService checks who are the top users and how many builds have they executed in the time window by going through the builds
if a user exists and has used the service it add one to their count if not it gives them  a new value of 1*/
func BuildRemoteService() int {
	ranking := 0  //initial rank of users is 0

	for _, data := range fields{//checks if the user is in the list
		if data.UserID == data.UserID{
			ranking++//if they are add one to their count
			fmt.Println(data.UserID)
		}else {
			ranking = 1 //else give them a new value
		}
	}
	return ranking
}
/*build success rate, and for builds that are not succeeding what are the top exit codes. succsessful build exit code is 0 */
func SucccessRate() float64{
	builds, successful := 0, 0
	for _, data := range fields{
		if data.ExitBuildProcess == 0{ //successful build output 0 count them
			successful++
		}
		fmt.Printf("%d top exit codes \n", data.ExitBuildProcess)
	}

	builds++ //counting the total

	successrate := (float64(successful)/float64(builds) * 100) //find the success rate
	return successrate
}

func checkError(msg string, err error)  {
	if err != nil{
		log.Fatal(msg,err)
	}
}




// NEW (return slice of CSVFields, and error if any)
// readStats returns a slice of CSVFields, and error if any
func readStats(filename string) ([]CsvFields, error) {
	//opens the stats.csv file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	} //if any errors break and return those errors



	//contents of the csv go to the lines variable
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}


	for _, line := range lines {
		// convert to proper types so that they can be stored in the struct(it has different formats)
		// handling errors
		received, err := time.Parse(time.RFC3339, line[2])
		if err != nil {
			return nil, err
		}
		began, _ := time.Parse(time.RFC3339, line[3])
		finished, _ := time.Parse(time.RFC3339, line[4])
		deleted, _ := strconv.ParseBool(line[5]) //converting to a boolean
		process, _ := strconv.Atoi(line[6])      //converting alphanumeric to an int
		size, _ := strconv.Atoi(line[7])
		//creating a new CsvFields objects from whatever we are reading from each line and each slice is a particular column
		data := CsvFields{
			ID:               line[0],
			UserID:           line[1],
			TimeReceived:     received,
			TimeBegan:        began,
			TimeFinished:     finished,
			DeletedBuilds:    deleted,
			ExitBuildProcess: process,
			ImageSize:        size,
		}

		//appending fields to the slice
		fields = append(fields, data)
	}

	// NEW (return fields and nil error)
	return fields, err
}