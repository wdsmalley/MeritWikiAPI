// Copyright 2010 The go-pgsql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package WikidbDAO

import (
	"fmt"
	"os"
	"github.com/lxn/go-pgsql"
)
func GetConnection() (*pgsql.Conn, error){
	//TODO get connection info from config file
	conn, err := pgsql.Connect("dbname=wikidb user=postgres password=lipscomb", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		os.Exit(1)
	}
	return conn, err
}