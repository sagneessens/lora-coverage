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
	DataRate   DataRate
	Power      int8
	RSSI       int16
	SNR        float64
	Size       uint16
	Payload    string
	Latitude   float64
	Longitude  float64
}

var (
	InvalidCrcError          = errors.New("invalid crc")
	InvalidMicError          = errors.New("invalid mic")
	InvalidMacPayloadError   = errors.New("invalid mac payload")
	InvalidFramePayloadError = errors.New("invalid frame payload")
	InvalidPayloadError      = errors.New("invalid payload")
)

func (c *Coverage) Unmarshal(data []byte) error {
	var packet RxPacket

	if err := json.Unmarshal(data, &packet); err != nil {
		return err
	}

	if packet.Crc < 0 {
		return InvalidCrcError
	}

	phyPayload, err := getDecryptedPayload([]byte(packet.Data))
	if err != nil {
		return err
	}

	macPayload, ok := phyPayload.MACPayload.(*lorawan.MACPayload)
	if !ok {
		return InvalidMacPayloadError
	}

	payload, ok := macPayload.FRMPayload[0].(*lorawan.DataPayload)
	if !ok {
		return InvalidFramePayloadError
	}

	c.GatewayMac = packet.GatewayMac
	c.DeviceAddr = macPayload.FHDR.DevAddr
	c.Time = packet.Time
	c.Frequency = packet.Frequency
	c.DataRate = packet.DataR
	c.RSSI = packet.RSSI
	c.SNR = packet.SNR
	c.Size = packet.Size
	c.Payload = hex.EncodeToString(payload.Bytes[:])

	lat, lon, err := getLocation(payload.Bytes[:])
	if err != nil {
		return err
	}

	c.Latitude = lat
	c.Longitude = lon

	pwr, err := getPower(payload.Bytes[:])
	c.Power = pwr
	if err != nil {
		return err
	}

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
		return nil, InvalidMicError
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

func getLocation(data []byte) (float64, float64, error) {
	if !isValidPayload(data) {
		return 0, 0, InvalidPayloadError
	}

	multiplier := float64(10000)

	latitude := float64(uint32(data[2])|uint32(data[1])<<8|uint32(data[0])<<16) / multiplier
	longitude := float64(uint32(data[5])|uint32(data[4])<<8|uint32(data[3])<<16) / multiplier

	return latitude, longitude, nil
}

func getPower(data []byte) (int8, error) {
	if !isValidPayload(data) || len(data) < 7 {
		return 127, InvalidPayloadError
	}

	power := int8(data[6])

	return power, nil
}

func isValidPayload(data []byte) bool {
	if len(data) < 6 || len(data) > 7 {
		return false
	} else {
		return true
	}
}
