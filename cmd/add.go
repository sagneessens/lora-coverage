// Copyright © 2018 Sven Agneessens <sven.agneessens@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/apex/log"
	"github.com/bullettime/lora-coverage/db"
	"github.com/bullettime/lora-coverage/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var message logMessage

type logMessage struct {
	Fields    interface{} `json:"fields"`
	Level     string      `json:"level"`
	TimeStamp string      `json:"timestamp"`
	Message   string      `json:"message"`
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add data from lora-logger json file",
	Long: `Process the data in a json file produced by the lora-logger tool.
It will only select the rx packets and add this data to a new or the existing database.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		database, err := db.Connect()
		if err != nil {
			log.WithError(err).WithField("database", viper.GetString("database.dbfile")).Fatal("connecting to database")
		}
		defer database.Disconnect()

		if err := database.Init(); err != nil {
			log.WithError(err).Fatal("initializing database")
		}

		myModel := model.New(database)

		jsonFileName := args[0]
		ctx := log.WithField("data-file", jsonFileName)

		jsonFile, err := os.Open(jsonFileName)
		if err != nil {
			ctx.WithError(err).Fatal("opening json file (arg)")
		}
		defer jsonFile.Close()

		lineScanner := bufio.NewScanner(jsonFile)
		for lineScanner.Scan() {
			err := json.Unmarshal(lineScanner.Bytes(), &message)
			if err != nil {
				log.WithError(err).WithField("line", string(lineScanner.Bytes())).Error("unmarshalling line")
			} else {
				if message.Message == "PUSH_DATA: RXPK" {
					var packet model.RxPacket

					fieldsJson, err := json.Marshal(message.Fields)
					if err != nil {
						log.WithError(err).WithField("fields", message.Fields).Error("marshalling fields")
					} else {
						if err := json.Unmarshal(fieldsJson, &packet); err != nil {
							log.WithError(err).WithField("fields", string(fieldsJson)).Error("unmarshalling fields")
						} else {
							myModel.AddRxPk(&packet)
						}
					}
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
