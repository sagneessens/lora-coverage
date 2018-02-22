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
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/pkg/errors"
)

type CompactTime time.Time

type DataRate struct {
	LoRa string
	FSK  uint32
}

type MacAddress [8]byte

type Payload struct {
	lorawan.PHYPayload
}

type RxPacket struct {
	GatewayMac MacAddress  `json:"gateway mac"`
	Time       CompactTime `json:"time"`
	Frequency  float64     `json:"frequency"`
	IFChannel  uint8       `json:"IF channel"`
	RFChain    uint8       `json:"RF chain"`
	Crc        int8        `json:"crc"`
	Modulation string      `json:"modulation"`
	DataR      *DataRate   `json:"data rate"`
	CodingRate string      `json:"coding rate"`
	RSSI       int16       `json:"rssi"`
	SNR        float64     `json:"snr"`
	Size       uint16      `json:"size"`
	Data       string      `json:"data"`
}

//func (p Payload) String() string {
//	b, err := p.MarshalJSON()
//	if err != nil {
//		return ""
//	}
//	return string(b)
//}
//
//func (p *Payload) MarshalJSON() ([]byte, error) {
//	return p.PHYPayload.MarshalJSON()
//}
//
//func (p *Payload) UnmarshalJSON(data []byte) error {
//	dataStr, err := strconv.Unquote(string(data))
//	if err != nil {
//		return err
//	}
//	return p.UnmarshalText([]byte(dataStr))
//}

func (t CompactTime) String() string {
	return time.Time(t).UTC().Format(time.RFC3339Nano)
}

func (t CompactTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.String())), nil
}

func (t *CompactTime) UnmarshalJSON(data []byte) error {
	t2, err := time.Parse(`"`+time.RFC3339Nano+`"`, string(data))
	if err != nil {
		return err
	}
	*t = CompactTime(t2)
	return nil
}

func (d DataRate) String() string {
	if d.LoRa != "" {
		return d.LoRa
	}

	return string(d.FSK)
}

func (d DataRate) MarshalJSON() ([]byte, error) {
	if d.LoRa != "" {
		return []byte(`"` + d.LoRa + `"`), nil
	}
	return []byte(strconv.FormatUint(uint64(d.FSK), 10)), nil
}

func (d *DataRate) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseUint(string(data), 10, 32)
	if err != nil {
		d.LoRa = strings.Trim(string(data), `"`)
		return nil
	}
	d.FSK = uint32(i)
	return nil
}

func (m MacAddress) String() string {
	b := make([]byte, len(m))
	for i, v := range m {
		b[i] = byte(v)
	}
	return fmt.Sprintf("%X", b)
}

func (m MacAddress) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(m.String())), nil
}

func (m *MacAddress) UnmarshalJSON(data []byte) error {
	dataStr, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	mac, err := hex.DecodeString(dataStr)
	if err != nil {
		return err
	}
	if len(mac) == 8 {
		copy(m[:], mac)
	} else {
		return errors.New(fmt.Sprintf("Wrong mac address size: %v bytes", len(mac)))
	}
	return nil
}

func (m *Model) AddRxPk(p *RxPacket) error {
	return m.db.AddRxPk(p)
}
