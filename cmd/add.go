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

type logMessage struct {
	Fields    json.RawMessage `json:"fields"`
	Level     string          `json:"level"`
	TimeStamp string          `json:"timestamp"`
	Message   string          `json:"message"`
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add data from lora-logger json file",
	Long: `lora-coverage add will process the data in a json file produced by the lora-logger tool.

This command takes one argument:
	- file name from the json file [eg. lora-log.json]
It will select the rx packets and add this data to a new or the existing database.`,
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

		dbModel := model.New(database)

		addDataFromGatewayLogger(args[0], dbModel)
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

func addDataFromGatewayLogger(fileName string, dbModel *model.Model) {
	ctx := log.WithField("data-file", fileName)

	jsonFile, err := os.Open(fileName)
	if err != nil {
		ctx.WithError(err).Fatal("opening json file (arg)")
	}
	defer jsonFile.Close()

	var message logMessage

	lineScanner := bufio.NewScanner(jsonFile)
	for lineScanner.Scan() {
		err := json.Unmarshal(lineScanner.Bytes(), &message)
		if err != nil {
			log.WithError(err).WithField("line", string(lineScanner.Bytes())).Error("unmarshalling line")
		} else {
			if message.Message == "PUSH_DATA: RXPK" {
				var coverage model.Coverage

				if err := coverage.Unmarshal(message.Fields); err != nil {
					ctx := log.WithField("fields", string(message.Fields))
					switch err {
					case model.InvalidCrcError:
						ctx.Debug(model.InvalidCrcError.Error())
					case model.InvalidMicError:
						ctx.Debug(model.InvalidMicError.Error())
					case model.InvalidMacPayloadError:
						ctx.Debug(model.InvalidMacPayloadError.Error())
					case model.InvalidFramePayloadError:
						ctx.Debug(model.InvalidFramePayloadError.Error())
					case model.InvalidPayloadError:
						ctx.Warn("invalid payload (no location data and/or power)")
						addCoverageRow(dbModel, &coverage)
					default:
						ctx.WithError(err).Error("error unmarshalling fields")
					}
				} else {
					addCoverageRow(dbModel, &coverage)
				}
			}
		}
	}
}

func addCoverageRow(dbModel *model.Model, row *model.Coverage) {
	err := dbModel.AddCoverageRow(row)
	if err != nil {
		log.WithError(err).Warn("Did you already scan this file?")
	}
}
