// The MIT License (MIT)
//
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

package db

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/bullettime/lora-coverage/model"
	"github.com/spf13/viper"
)

var (
	tearDown  = true
	message   logMessage
	messages  []logMessage
	rxpackets []model.RxPacket
)

type logMessage struct {
	Fields    interface{} `json:"fields"`
	Level     string      `json:"level"`
	TimeStamp string      `json:"timestamp"`
	Message   string      `json:"message"`
}

func TestRxPk(t *testing.T) {
	// setup
	jsonFile, err := os.Open("../lora.json")
	if err != nil {
		t.Fatal("error opening json testfile:", err)
	}
	defer jsonFile.Close()
	lineScanner := bufio.NewScanner(jsonFile)

	for lineScanner.Scan() {
		err := json.Unmarshal(lineScanner.Bytes(), &message)
		if err != nil {
			t.Error("error unmarshaling line:", lineScanner.Text(), "(", err, ")")
		} else {
			messages = append(messages, message)
		}

	}

	for _, m := range messages {
		if m.Message == "PUSH_DATA: RXPK" {
			var p model.RxPacket
			fieldsJson, err := json.Marshal(m.Fields)
			if err != nil {
				t.Error("error marshaling fields to json:", m.Fields, "(", err, ")")
			}
			json.Unmarshal(fieldsJson, &p)
			rxpackets = append(rxpackets, p)
		}
	}

	viper.Set("database.dbfile", "test_rxpk.sqlite")
	db, err := Connect()
	if err != nil {
		t.Fatal("error opening database:", err)
	}
	defer db.Disconnect()

	// tests
	t.Run("TestInitRxPk", func(t *testing.T) {
		if err := db.initRxPk(); err != nil {
			t.Error("error initRxPk:", err)
		}
	})

	t.Run("TestAddRxPk", func(t *testing.T) {
		for _, packet := range rxpackets {
			if err := db.AddRxPk(&packet); err != nil {
				t.Error("error AddRxPk:", err)
			}
		}
	})

	// tear-down
	if tearDown {
		if err := os.Remove(viper.GetString("database.dbfile")); err != nil {
			t.Log("error deleting database:", err)
		}
	}
}
