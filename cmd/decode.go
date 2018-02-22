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
	"fmt"
	"github.com/apex/log"
	"github.com/brocaar/lorawan"
	"github.com/spf13/cobra"
)

var (
	nwkSKey = [16]byte{200, 243, 24, 62, 142, 24, 220, 221, 149, 70, 183, 25, 178, 79, 176, 124}
	appSKey = [16]byte{53, 225, 222, 151, 84, 90, 143, 185, 215, 10, 186, 32, 243, 199, 190, 181}
	appKey = [16]byte{228, 251, 27, 233, 86, 88, 213, 168, 184, 166, 7, 93, 169, 105, 117, 121}
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := args[0]
		log.WithField("data", data).Info("decoding")

		var phy lorawan.PHYPayload
		if err := phy.UnmarshalText([]byte(data)); err != nil {
			log.WithError(err).Fatal("unmarshalling data")
		}

		log.WithFields(log.Fields{
			"MType": phy.MHDR.MType,
			"Major": phy.MHDR.Major,
		}).Info("test")

		ok, err := phy.ValidateMIC(nwkSKey)
		if err != nil {
			log.WithError(err).Fatal("validating mic")
		}
		if !ok {
			log.Fatal("invalid mic")
		}

		phyJSON, err := phy.MarshalJSON()
		if err != nil {
			log.WithError(err).Fatal("marshalling json")
		}

		log.WithField("phyJSON", phyJSON).Info("phy")

		if err := phy.DecryptFRMPayload(appSKey); err != nil {
			log.WithError(err).Fatal("decrypting payload")
		}
		macPL, ok := phy.MACPayload.(*lorawan.MACPayload)
		if !ok {
			log.Fatal("*MACPayload expected")
		}

		pl, ok := macPL.FRMPayload[0].(*lorawan.DataPayload)
		if !ok {
			log.Fatal("*DataPayload expected")
		}

		log.WithFields(log.Fields{
			"bytes": fmt.Sprint(pl.Bytes),
			"string": fmt.Sprintf("%s", pl.Bytes),
		}).Info("frame payload")
	},
}

func init() {
	RootCmd.AddCommand(decodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
