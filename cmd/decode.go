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
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/brocaar/lorawan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode LoRaWAN payload",
	Long: `lora-coverage decode will decode the data from a LoRaWAN payload with the configured keys.

This command expects one argument, the payload in base64 string format.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := args[0]
		log.WithField("data", data).Info("decoding")

		var phy lorawan.PHYPayload
		if err := phy.UnmarshalText([]byte(data)); err != nil {
			log.WithError(err).Fatal("unmarshalling data")
		}

		var nwkSKey lorawan.AES128Key
		if err := nwkSKey.UnmarshalText([]byte(viper.GetString("lora.nwkskey"))); err != nil {
			log.WithError(err).Fatal("creating network session key")
		}

		ok, err := phy.ValidateMIC(nwkSKey)
		if err != nil {
			log.WithError(err).Fatal("validating mic")
		}
		if !ok {
			log.Fatal("invalid mic")
		}

		phyJSON, err := json.MarshalIndent(phy, "", "  ")
		if err != nil {
			log.WithError(err).Fatal("marshalling json")
		}

		fmt.Printf("LoRaWAN Packet:\n%s\n", phyJSON)

		var appSKey lorawan.AES128Key
		if err := appSKey.UnmarshalText([]byte(viper.GetString("lora.appskey"))); err != nil {
			log.WithError(err).Fatal("creating app session key")
		}

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

		fmt.Printf("Payload: %X\n", pl.Bytes)
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
