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
	"github.com/bullettime/lora-coverage/model"
	"github.com/pkg/errors"
)

var (
	createCoverageTable = `CREATE TABLE IF NOT EXISTS coverage(
gateway TEXT,
device TEXT,
time TEXT,
frequency REAL,
rssi INTEGER,
snr REAL,
size INTEGER,
payload TEXT,
create_time TEXT DEFAULT CURRENT_TIMESTAMP)`
	addCoverageRow = `INSERT INTO coverage(gateway, device, time, frequency, rssi, snr, size, payload) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
)

func (c *Connection) initCoverage() error {
	_, err := c.database.Exec(createCoverageTable)
	if err != nil {
		return errors.Wrap(err, "error initializing table 'coverage'")
	}
	return nil
}

func (c *Connection) AddCoverageRow(m *model.Coverage) error {
	_, err := c.database.Exec(addCoverageRow, m.GatewayMac.String(), m.DeviceAddr.String(), m.Time.String(),
		m.Frequency, m.RSSI, m.SNR, m.Size, m.Payload)
	if err != nil {
		return errors.Wrapf(err, "error adding coverage row: %s", m)
	}

	return nil
}
