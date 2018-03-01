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
	"database/sql"

	"github.com/bullettime/lora-coverage/model"
	"github.com/pkg/errors"
	"github.com/paulmach/go.geojson"
)

var (
	createCoverageTable = `CREATE TABLE IF NOT EXISTS coverage(
gateway TEXT NOT NULL,
device TEXT NOT NULL,
time TEXT NOT NULL,
frequency REAL NOT NULL,
datarate STRING NOT NULL,
power INTEGER,
rssi INTEGER NOT NULL,
snr REAL NOT NULL,
size INTEGER NOT NULL,
payload TEXT NOT NULL,
lat REAL,
lon REAL,
create_time TEXT DEFAULT CURRENT_TIMESTAMP,
UNIQUE(gateway, device, time, payload))`
	addCoverageRow = `INSERT INTO coverage(gateway, device, time, frequency, datarate, power, rssi, snr, size, 
payload, lat, lon) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	getGeoJsonPoints = `SELECT rssi, lat, lon FROM coverage WHERE gateway=? AND datarate=?`
)

func (c *Connection) initCoverage() error {
	_, err := c.database.Exec(createCoverageTable)
	if err != nil {
		return errors.Wrap(err, "error initializing table 'coverage'")
	}
	return nil
}

func (c *Connection) AddCoverageRow(m *model.Coverage) error {
	power := getNullPower(m.Power)
	latitude := getNullLatLon(m.Latitude)
	longitude := getNullLatLon(m.Longitude)
	_, err := c.database.Exec(addCoverageRow, m.GatewayMac.String(), m.DeviceAddr.String(), m.Time.String(),
		m.Frequency, m.DataRate.String(), power, m.RSSI, m.SNR, m.Size, m.Payload, latitude, longitude)
	if err != nil {
		return errors.Wrapf(err, "error adding coverage row: %s", m)
	}

	return nil
}

func getNullLatLon(value float64) sql.NullFloat64 {
	if value == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: value,
		Valid:   true,
	}
}

func getNullPower(value int8) sql.NullInt64 {
	if value == 127 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(value),
		Valid: true,
	}
}

func (c *Connection) GetGeoJSonPoints(gateway string, datarate string) ([]*geojson.Feature, error) {
	var points []*geojson.Feature

	rows, err := c.database.Query(getGeoJsonPoints, gateway, datarate)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving geo json points")
	}
	defer rows.Close()

	for rows.Next() {
		var rssi int
		var lat, lon float64

		if err := rows.Scan(&rssi, &lat, &lon); err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}

		feature := geojson.NewPointFeature([]float64{lat,lon})
		feature.SetProperty("rssi", rssi)

		points = append(points, feature)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error in rows")
	}

	return points, nil
}
