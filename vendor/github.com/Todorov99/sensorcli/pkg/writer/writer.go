package writer

import (
	"encoding/csv"
	"os"
)

var (
	csvFileFullPath string = "./"
)

// WriteOutputToCSV writes sensor measurement data into CSV file
func WriteOutputToCSV(data []string, csvFileName string) chan error {
	done := make(chan error)

	go func() {

		defer func() {
			close(done)
		}()

		fileName := csvFileFullPath + csvFileName + ".csv"

		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			done <- err
		}

		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		writer.Comma = '|'

		writingErr := writer.Write(data)
		if writingErr != nil {
			done <- err

		}
	}()

	return done
}
