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

package model

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

var (
	tearDown = true
	message  logMessage
	messages []logMessage
)

type logMessage struct {
	Fields    interface{} `json:"fields"`
	Level     string      `json:"level"`
	TimeStamp string      `json:"timestamp"`
	Message   string      `json:"message"`
}

func TestRxPk(t *testing.T) {
	// setup
	jsonFile, err := os.Open("../test.json")
	if err != nil {
		t.Fatal("error opening json testfile:", err)
	}
	defer jsonFile.Close()
	lineScanner := bufio.NewScanner(jsonFile)

	// tests
	t.Run("A", func(t *testing.T) {
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
				var p RxPacket
				fieldsJson, err := json.Marshal(m.Fields)
				if err != nil {
					t.Error("error marshaling fields to json:", m.Fields, "(", err, ")")
				}
				t.Log(string(fieldsJson))
				json.Unmarshal(fieldsJson, &p)
				t.Log(&p)
				t.Log("Gateway Mac:", p.GatewayMac)
				t.Log("Time:", p.Time)
				t.Log("Frequency:", p.Frequency)
				t.Log("IFChannel:", p.IFChannel)
				t.Log("RFChain:", p.RFChain)
				t.Log("CRC:", p.Crc)
				t.Log("Modulation:", p.Modulation)
				t.Log("Data rate:", p.DataR)
				t.Log("Coding rate:", p.CodingRate)
				t.Log("RSSI:", p.RSSI)
				t.Log("SNR:", p.SNR)
				t.Log("Size:", p.Size)
				t.Log("data:", p.Data)
			}
		}
	})

	// tear-down
	if tearDown {
	}
}
