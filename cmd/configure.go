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
	"encoding/hex"

	"os"

	"github.com/apex/log"
	"github.com/segmentio/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	Database databaseConfig `yaml:"database"`
	Lora     loraConfig     `yaml:"lora"`
}

type databaseConfig struct {
	DBFile string `yaml:"dbfile"`
}

type loraConfig struct {
	NwkSKey string `yaml:"nwkskey"`
	AppSKey string `yaml:"appskey"`
}

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure lora-coverage",
	Long: `lora-coverage configure creates a yaml configuration file for the coverage tool.

Various different values for settings that are needed to use this tool are asked.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			newDBFile  string
			newNwkSKey string
			newAppSKey string
		)

		newDBFile = prompt.String("database location [eg. coverage.db]")
		if len(newDBFile) == 0 {
			newDBFile = viper.GetString("database.dbfile")
		}

		for len(newNwkSKey) != 16*2 {
			newNwkSKey = prompt.StringRequired("network session key in hex (16 bytes) [eg. 0102030405060708090A0B0C0D0E0F11]")
			_, err := hex.DecodeString(newNwkSKey)
			if err != nil {
				log.WithError(err).Error("non hex network session key")
			}
		}

		for len(newAppSKey) != 16*2 {
			newAppSKey = prompt.StringRequired("application session key in hex (16 bytes) [eg. 0102030405060708090A0B0C0D0E0F11]")
			_, err := hex.DecodeString(newAppSKey)
			if err != nil {
				log.WithError(err).Error("non hex application session key")
			}
		}

		newConfig := &yamlConfig{
			Database: databaseConfig{
				DBFile: newDBFile,
			},
			Lora: loraConfig{
				NwkSKey: newNwkSKey,
				AppSKey: newAppSKey,
			},
		}

		output, err := yaml.Marshal(newConfig)
		if err != nil {
			log.WithError(err).Fatal("failed generating yaml config")
		}

		if len(viper.ConfigFileUsed()) == 0 {
			viper.SetConfigFile(cfgFile)
		}

		f, err := os.Create(viper.ConfigFileUsed())
		if err != nil {
			log.WithError(err).Fatal("failed creating log file")
		}

		defer f.Close()

		f.Write(output)
		log.WithField("path", viper.ConfigFileUsed()).Debug("new configuration file saved")
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
