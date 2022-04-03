/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/Todorov99/sensorcli/pkg/util"
	"github.com/Todorov99/sensorcli/pkg/writer"
	"github.com/spf13/cobra"
)

var (
	format         string
	deltaDuration  int64
	sensorGroups   []string
	totalDuration  float64
	file           string
	webHook        string
	generateReport bool
	reportType     string
)

// cpuCmd represents the cpu command
var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		interuptSignal := make(chan os.Signal, 1)

		signal.Notify(interuptSignal, os.Interrupt)
		defer func() {
			signal.Stop(interuptSignal)
			cancel()
		}()

		go func() {
			select {
			case <-interuptSignal:
				cancel()
			case <-ctx.Done():
			}
		}()

		err := terminateForTotalDuration(ctx)
		if err != nil {
			cmdLogger.Error(err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpuCmd)

	cpuCmd.Flags().StringVar(&format, "format", "JSON", "The data could be printed either JSON or YAML format.")
	cpuCmd.Flags().Int64Var(&deltaDuration, "delta_duration", 3, "The period of which you will get your sensor data.")
	cpuCmd.Flags().StringVar(&file, "output_file", "", "Writing the output into CSV file.")
	cpuCmd.Flags().Float64Var(&totalDuration, "total_duration", 60.0, "Terminating the whole program after specified duration")
	cpuCmd.Flags().StringVar(&webHook, "web_hook_url", "", "Expose to current port.")
	cpuCmd.Flags().StringVar(&reportType, "reportType", "xlsx", "The type of the report file that has to be generated")
	cpuCmd.Flags().BoolVar(&generateReport, "generateReport", false, "generate xslx report file")

	cpuCmd.Flags().StringSliceVar(&sensorGroups, "sensor_group", []string{""}, "There are three main sensor groups: CPU_TEMP, CPU_USAGE and MEMORY_USAGE. Each senosr group could have system file that will hold specific information")
}

func getSensorGroupsWithSystemFile(sensorflag []string) map[string]string {
	sensorGroupWithSysFile := make(map[string]string)

	for _, group := range sensorflag {
		sysFile := ""
		splitArgs := strings.Split(group, "=")

		if len(splitArgs) > 1 && splitArgs[1] != "" {
			sysFile = splitArgs[1]
		}

		sensorGroupWithSysFile[splitArgs[0]] = sysFile
	}

	return sensorGroupWithSysFile
}

func terminateForTotalDuration(ctx context.Context) error {
	appTerminaitingDuration := time.After(getTotalDurationInSeconds())
	device, err := loadDeviceConfig()
	if err != nil {
		return err
	}
	groups := getSensorGroupsWithSystemFile(sensorGroups)
	cpu := NewCpu(groups)
	reportWriter := writer.New("measurement_" + time.Now().Format(time.RFC3339) + "." + reportType)

	for {
		select {
		case <-ctx.Done():
			cmdLogger.Error(ctx.Err())
			return ctx.Err()
		case <-appTerminaitingDuration:
			return nil
		default:
			multipleSensorsData, err := cpu.GetMeasurements(ctx, device)
			if err != nil {
				cmdLogger.Error(err)
				return err
			}

			err = getMeasurementsInDeltaDuration(ctx, reportWriter, generateReport, multipleSensorsData, getDeltaDurationInSeconds())
			if err != nil {
				return err
			}
		}
	}

}

func getMultipleSensorsMeasurements(ctx context.Context, groups map[string]string) ([]sensor.Measurment, error) {
	var multipleSensorsData []sensor.Measurment
	for group, sysFile := range groups {
		var currentSensorGroupData []sensor.Measurment

		currentSensorGroupData, err := getSensorMeasurements(ctx, group, sysFile)
		if err != nil {
			return nil, err
		}

		multipleSensorsData = append(multipleSensorsData, currentSensorGroupData...)
	}

	return multipleSensorsData, nil
}

func getMeasurementsInDeltaDuration(ctx context.Context, reportWriter writer.ReportWriter, generateReport bool, sensorData []sensor.Measurment, deltaDuration time.Duration) error {
	cmdLogger.Info("Getting measurements in delta duration...")

	measurementDuration := time.After(deltaDuration)
	done := make(chan bool)
	sensorsData := sendSensorData(sensorData, done)

	defer func() {
		close(done)
	}()

	for {
		select {
		case data := <-sensorsData:
			if webHook != "" {
				go func() {
					webHookURL(webHook, util.ParseDataAccordingToFormat("JSON", data))
				}()
			}

			if generateReport {
				go func() {

					if reportType == "xlsx" {
						err := reportWriter.WritoToXslx(sensorData)
						if err != nil {
							cmdLogger.Errorf("error during writing in XLSX file: %w", err)
						}
					}

					if reportType == "csv" {
						var sensorsData []string
						sensorsData = append(sensorsData, util.ParseDataAccordingToFormat(format, data))
						err := reportWriter.WriteOutputToCSV(sensorsData)
						if err != nil {
							cmdLogger.Errorf("error during writing in CSV file: %w", err)
						}
					}
				}()
			}

			fmt.Println(util.ParseDataAccordingToFormat(format, data))
		case <-measurementDuration:
			done <- true
			return nil
		case <-ctx.Done():
			done <- true
			cmdLogger.Error(ctx.Err())
			return ctx.Err()
		}
	}

}

func sendSensorData(sensorsInfo []sensor.Measurment, done chan bool) <-chan sensor.Measurment {
	cmdLogger.Info("Sending sensor data...")

	out := make(chan sensor.Measurment)

	go func() {
		for _, currentSensorInfo := range sensorsInfo {
			out <- currentSensorInfo
		}

		if <-done {
			close(out)
			return
		}
	}()

	return out
}

func getDeltaDurationInSeconds() time.Duration {
	return time.Duration(deltaDuration) * time.Second
}

func getTotalDurationInSeconds() time.Duration {
	return time.Duration(totalDuration) * time.Second
}

func webHookURL(url string, data string) {
	var json = []byte(data)
	http.Post(url, "application/json", bytes.NewBuffer(json))
}
