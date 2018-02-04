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
	"github.com/pkg/errors"
	"github.com/bullettime/lora-coverage/model"
)

var (
	createTable = `CREATE TABLE IF NOT EXISTS rxpk(
gatewaymac BLOB,
time TEXT,
frequency REAL,
ifchannel BLOB,
rfchain BLOB,
crc INTEGER,
modulation TEXT,
datarate TEXT,
codingrate TEXT,
rssi INTEGER,
snr REAL,
size INTEGER,
data TEXT,
create_time TEXT DEFAULT CURRENT_TIMESTAMP)`
	addRxPk = `INSERT INTO rxpk(gatewaymac, time, frequency, ifchannel, rfchain, crc, modulation, datarate, codingrate, 
rssi, snr, size, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

func (c *Connection) initRxPk() error {
	_, err := c.database.Exec(createTable)
	if err != nil {
		return errors.Wrap(err, "error initializing table 'rxpk'")
	}
	return nil
}

func (c *Connection) AddRxPk(p *model.RxPacket) error {
	_, err := c.database.Exec(addRxPk, p.GatewayMac.String(), p.Time.String(), p.Frequency, p.IFChannel, p.RFChain,
		p.Crc, p.Modulation, p.DataR.String(), p.CodingRate, p.RSSI, p.SNR, p.Size, p.Data)
	if err != nil {
		return errors.Wrapf(err, "error adding rxpk: %s", p)
	}
	return nil
}
