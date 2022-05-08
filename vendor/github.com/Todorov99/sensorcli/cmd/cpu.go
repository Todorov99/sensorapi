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
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Todorov99/sensorcli/pkg/client"
	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/Todorov99/sensorcli/pkg/util"
	"github.com/Todorov99/sensorcli/pkg/writer"
	"github.com/spf13/cobra"
)

var (
	format        string
	deltaDuration int64
	sensorGroups  []string
	totalDuration float64
	//file           string
	webHook        string
	mailHook       string
	email          string
	username       string
	password       string
	generateReport bool
	reportType     string
	configDirPath  string
	rootCAPath     string
)

// cpuCmd represents the cpu command
var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Start collecting hardware metrics",
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

	cpuCmd.Flags().StringVar(&format, "format", "JSON", "The data could be printed either JSON or YAML format")
	cpuCmd.Flags().Int64Var(&deltaDuration, "delta_duration", 3, "The period between each measurement")
	cpuCmd.Flags().Float64Var(&totalDuration, "total_duration", 60.0, "The period in which the hardware measurements will be collected")
	cpuCmd.Flags().StringVar(&username, "username", "", "The username of the user used for remote authentication to the sensor API")
	cpuCmd.Flags().StringVar(&password, "password", "", "The password of the user used for remote authentication to the sensor API")
	cpuCmd.Flags().StringVar(&webHook, "web_hook_url", "", "Base URL to the sensor API")
	cpuCmd.Flags().StringVar(&mailHook, "mail_hook_url", "", "Base URL to the mailsender API")
	cpuCmd.Flags().StringVar(&email, "email", "", "The email to which the final result should be send")
	cpuCmd.Flags().StringVar(&rootCAPath, "rootCAPath", "", "The path to the root Certificate Authoritate (CA) used for the TLS client config. If no CA is provided the verification of the certificates is skipped")
	cpuCmd.Flags().StringVar(&reportType, "reportType", "xlsx", "The type of the report file that has to be generated. Possible values xlsx, csv.")
	cpuCmd.Flags().StringVar(&configDirPath, "configDirPath", "", "The path to the configuration file for the device for which the hardware measurements will be collected")
	cpuCmd.Flags().BoolVar(&generateReport, "generateReport", false, "Flag that shows whether to generate report file")
	cpuCmd.Flags().StringSliceVar(&sensorGroups, "sensor_group", []string{""}, "There are three main sensor groups: CPU_TEMP, CPU_USAGE and MEMORY_USAGE. Each senosor group could have system file that will hold specific information")
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
	device, err := loadDeviceConfig(configDirPath)
	if err != nil {
		return err
	}
	groups := getSensorGroupsWithSystemFile(sensorGroups)
	cpu := NewCpu(groups)
	reportFile := "measurement_" + time.Now().Format(sensor.TimeFormat) + "." + reportType
	reportWriter := writer.New(reportFile)

	if mailHook != "" && email == "" {
		return fmt.Errorf("email for sending the result has not been specified")
	}

	for {
		select {
		case <-ctx.Done():
			cmdLogger.Error(ctx.Err())
			return ctx.Err()
		case <-appTerminaitingDuration:
			if mailHook != "" {
				mailSenderClient := client.NewMailSenderTLSClient(mailHook, rootCAPath)
				sender := client.MailSender{
					Subject: "Measurements from the CLI",
					To: []string{
						email,
					},
					Body: "Measurements started from the CLI has finished",
				}
				if generateReport {
					err := mailSenderClient.SendWithAttachments(ctx, sender, []string{reportFile})
					if err != nil {
						return err
					}
				} else {
					err := mailSenderClient.Send(ctx, sender)
					if err != nil {
						return err
					}
				}
			}
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
	errChan := make(chan error)

	sensorsData := sendSensorData(sensorData, done)

	defer func() {
		close(done)
	}()

	apiClient := client.NewAPITLSClient(ctx, webHook, rootCAPath)

	if generateReport {

		go func() {
			if reportType == "xlsx" {
				err := reportWriter.WritoToXslx(sensorData)
				if err != nil {
					errChan <- fmt.Errorf("error during writing in XLSX file: %w", err)
					return
				}
			}

			if reportType == "csv" {
				var sensorsData []string
				sensorsData = append(sensorsData, util.ParseDataAccordingToFormat(format, sensorData))
				err := reportWriter.WriteOutputToCSV(sensorsData)
				if err != nil {
					errChan <- fmt.Errorf("error during writing in CSV file: %w", err)
					return
				}
			}
		}()
	}

	for {
		select {
		case data := <-sensorsData:
			if webHook != "" {
				go func() {
					resp := apiClient.SendMetrics(ctx, username, password, data)
					if resp.Err != nil {
						errChan <- resp.Err
						return
					}
				}()
			}

			fmt.Println(util.ParseDataAccordingToFormat(format, data))
		case <-measurementDuration:
			done <- true
			return nil
		case err := <-errChan:
			done <- true
			cmdLogger.Error(err)
			return err
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
