package writer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/xuri/excelize/v2"
)

var (
	sheetName string = "Measurements"
)

type ReportWriter interface {
	WriteOutputToCSV(data []string) error
	WritoToXslx(measurements []sensor.Measurment) error
}

type reportWriter struct {
	file string
	mx   sync.RWMutex
}

func New(file string) ReportWriter {
	return &reportWriter{
		file: file,
	}
}

// WriteOutputToCSV writes sensor measurement data into CSV file
func (r *reportWriter) WriteOutputToCSV(data []string) error {
	r.mx.RLock()
	defer r.mx.RUnlock()
	file, err := os.OpenFile(r.file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Comma = '|'

	writingErr := writer.Write(data)
	if writingErr != nil {
		return err
	}

	return err
}

// WritoToXslx writes sensor measurement data into Xslx file
func (r *reportWriter) WritoToXslx(measurements []sensor.Measurment) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	var f *excelize.File
	f, err := excelize.OpenFile(r.file)
	if errors.Is(err, os.ErrNotExist) {
		f = excelize.NewFile()
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	defer f.Close()
	activeSheetIndex := f.NewSheet(sheetName)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	startingIndex := len(rows) + 1
	if startingIndex == 0 {
		f.SetCellValue(sheetName, "A1", "measuredAt")
		f.SetCellValue(sheetName, "B1", "value")
		f.SetCellValue(sheetName, "C1", "sensorID")
		f.SetCellValue(sheetName, "D1", "deviceID")
		f.SetCellValue(sheetName, "E1", "unit")
		startingIndex = 2
	}

	for _, m := range measurements {
		timeStampAxis := fmt.Sprintf("A%d", startingIndex)
		valueAxis := fmt.Sprintf("B%d", startingIndex)
		sensorAxis := fmt.Sprintf("C%d", startingIndex)
		deviceAxis := fmt.Sprintf("D%d", startingIndex)
		unitAxis := fmt.Sprintf("E%d", startingIndex)

		f.SetCellValue(sheetName, timeStampAxis, m.MeasuredAt)
		f.SetCellValue(sheetName, valueAxis, m.Value)
		f.SetCellValue(sheetName, sensorAxis, m.SensorID)
		f.SetCellValue(sheetName, deviceAxis, m.DeviceID)
		f.SetCellValue(sheetName, unitAxis, m.Unit)
		startingIndex++
	}
	f.SetActiveSheet(activeSheetIndex)
	return f.SaveAs(r.file)
}
