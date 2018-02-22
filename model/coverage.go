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
	"encoding/json"

	"github.com/brocaar/lorawan"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Coverage struct {
	GatewayMac MacAddress
	DeviceAddr lorawan.DevAddr
	Time       CompactTime
	Frequency  float64
	RSSI       int16
	SNR        float64
	Size       uint16
	Payload    string
}

func (c *Coverage) Unmarshal(data []byte) error {
	var packet RxPacket

	if err := json.Unmarshal(data, &packet); err != nil {
		return err
	}

	phyPayload, err := getDecryptedPayload([]byte(packet.Data))
	if err != nil {
		return err
	}

	macPayload, ok := phyPayload.MACPayload.(*lorawan.MACPayload)
	if !ok {
		return errors.New("wrong payload type")
	}

	payload, ok := macPayload.FRMPayload[0].(*lorawan.DataPayload)
	if !ok {
		return errors.New("wrong payload type")
	}

	c.GatewayMac = packet.GatewayMac
	c.DeviceAddr = macPayload.FHDR.DevAddr
	c.Time = packet.Time
	c.Frequency = packet.Frequency
	c.RSSI = packet.RSSI
	c.SNR = packet.SNR
	c.Size = packet.Size
	c.Payload = string(payload.Bytes)

	return nil
}

func getDecryptedPayload(data []byte) (*lorawan.PHYPayload, error) {
	var phy lorawan.PHYPayload
	if err := phy.UnmarshalText(data); err != nil {
		return nil, err
	}

	var nwkSKey lorawan.AES128Key
	if err := nwkSKey.UnmarshalText([]byte(viper.GetString("lora.nwkskey"))); err != nil {
		return nil, err
	}

	ok, err := phy.ValidateMIC(nwkSKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid mic")
	}

	var appSKey lorawan.AES128Key
	if err := appSKey.UnmarshalText([]byte(viper.GetString("lora.appskey"))); err != nil {
		return nil, err
	}

	if err := phy.DecryptFRMPayload(appSKey); err != nil {
		return nil, err
	}

	return &phy, nil
}
