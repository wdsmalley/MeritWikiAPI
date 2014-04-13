// Copyright 2010 The go-pgsql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"time"
)

import (
	"github.com/lxn/go-pgsql"
	//"database/sql"
)


func sleepnow(){
	duration := time.Second
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
}
func createUser(){
	conn, err := pgsql.Connect("dbname=wikidb user=postgres password=lipscomb", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		sleepnow()
		os.Exit(1)
	}
	defer conn.Close()
	command := "select create_wikiuser(@username, @email, @password)"
	usernameParam := pgsql.NewParameter("@username", pgsql.Text)
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	passwordParam := pgsql.NewParameter("@password", pgsql.Text)
	stmt, err := conn.Prepare(command, usernameParam, emailParam,passwordParam )
	defer stmt.Close()
	err = usernameParam.SetValue("Gotest")
	if err != nil {
		os.Exit(1)
	}
	err = emailParam.SetValue("gotest@mail.com")
	if err != nil {
		os.Exit(1)
	}
	err = passwordParam.SetValue("mypass12")
	if err != nil {
		os.Exit(1)
	}
	rs, err := stmt.Query()
	if err != nil {
		fmt.Println("error in query")
		sleepnow()
		os.Exit(1)
	}else{
		fmt.Println(rs)
	}
	//defer rs.Close()
	fmt.Println("successful insert")
	sleepnow()
}
func lookupUserNameByEmail(inEmail string){
	conn, err := pgsql.Connect("dbname=wikidb user=postgres", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		sleepnow()
		os.Exit(1)
	}
	defer conn.Close()
	command := "select lookup_username(@email);"
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	stmt, err := conn.Prepare(command, emailParam)
	if err != nil {
		fmt.Println("error in prepare")
		sleepnow()
		os.Exit(1)
	}
	defer stmt.Close()
	
	err = emailParam.SetValue(inEmail)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("running the query")
	rs, err := stmt.Query()
	if err != nil {
		fmt.Println("error in query")
		sleepnow()
		os.Exit(1)
	}else{
		fmt.Println(rs)
	}
	defer rs.Close()

	for {
		hasRow, err := rs.FetchNext()
		if err != nil {
			fmt.Println("error in fetch:")
			sleepnow()
			os.Exit(1)
		}

		if !hasRow {
			break
		}else{
			username, isNull, err := rs.String(0)
			if err != nil {
				fmt.Println(err)
				sleepnow()
				os.Exit(1)
			}
			if isNull {
				username = "(null)"
			}
		   fmt.Println("username: ", username)
		}

	}
	sleepnow()
}		
func main() {
	fmt.Println("starting the test")
	createUser()
	lookupUserNameByEmail("gotest@mail.com")
	//lookupUserNameByEmail("bill2@test.com")
	//lookupUserNameByEmail("bill3@test.com")
	
	
}