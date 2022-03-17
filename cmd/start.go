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
	_ "github.com/Todorov99/server/pkg/database"
	"github.com/Todorov99/server/pkg/server"
	"github.com/spf13/cobra"
)

var (
	port string
)

// startCmd represents the start command for starting the HTTP(S) Web Server
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start http web server",
	Long:  `Start is command line that starts a http web server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.HandleRequest(port)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "The port of the server")
}
