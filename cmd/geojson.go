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
	"github.com/spf13/cobra"
	"github.com/paulmach/go.geojson"
	"github.com/apex/log"
	"github.com/bullettime/lora-coverage/model"
	"github.com/spf13/viper"
	"github.com/bullettime/lora-coverage/db"
	"io/ioutil"
	"fmt"
)

var (
	callback = "eqfeed_callback"
	output = "data_geo.json"
)

// geojsonCmd represents the geojson command
var geojsonCmd = &cobra.Command{
	Use:   "geojson",
	Short: "Create a geo jsonp file from the data",
	Long: `lora-coverage geojson creates a geo jsonp file from the data currently in the database.

This command takes two arguments:
	1. gatewac mac (in hex) [eg. 008000000000b88d]
	2. datarate [eg. SF7BW125]
The arguments have to be entered in that order.`,
	Args: cobra.ExactArgs(2),
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

		points, err := dbModel.GetGeoJSonPoints(args[0], args[1])
		if err != nil {
			log.WithError(err).Fatal("getting geo json points")
		}

		featureCollection := geojson.NewFeatureCollection()
		for _, point := range points {
			featureCollection.AddFeature(point)
		}

		rawJson, err := featureCollection.MarshalJSON()
		if err != nil {
			log.WithError(err).Fatal("marshalling json")
		}

		jsonP := fmt.Sprintf("%s(%s);", callback, rawJson)

		ioutil.WriteFile(output, []byte(jsonP), 0644)
	},
}

func init() {
	RootCmd.AddCommand(geojsonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// geojsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// geojsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	geojsonCmd.Flags().StringVarP(&callback, "callback", "c", "eqfeed_callback", "name of the callback function")
	geojsonCmd.Flags().StringVarP(&output, "output", "o", "data_geo.json", "name of the output file")
}
