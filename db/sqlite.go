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

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Connection struct {
	database *sql.DB
}

func Connect() (*Connection, error) {
	dbFile := viper.GetString("database.dbfile")
	if len(dbFile) == 0 {
		return nil, errors.New("database file not found (did you run configure?)")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening database: %s", dbFile)
	}

	return &Connection{database: db}, nil
}

func (c *Connection) Init() error {
	if err := c.initRxPk(); err != nil {
		return err
	}

	return nil
}

func (c *Connection) Ping() error {
	return c.database.Ping()
}

func (c *Connection) Disconnect() error {
	err := c.database.Close()
	if err != nil {
		return errors.Wrap(err, "error closing the database")
	}

	return nil
}
