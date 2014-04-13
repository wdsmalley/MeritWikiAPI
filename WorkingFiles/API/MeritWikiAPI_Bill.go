// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"net/http"
	"model"
	"WikidbDAO"
	"fmt"
	"time"
	"encoding/json"
    "io/ioutil"
    "strings"
)

func sleepnow(){
	duration := time.Second
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	time.Sleep(duration);
	
}


func createUserHandler(w http.ResponseWriter, r *http.Request, t string) {
	fmt.Println("starting the create user test")
	//create a user
	u := model.NewUser("teststruct","gotest2@mail.com","mypass12")
	WikidbDAO.CreateUser(u)
	//wait to make sure the create is complete and committed to DB
	sleepnow()
	//find the user just created
	WikidbDAO.LookupUserNameByEmail(u)
	
}
func createPageHandler(w http.ResponseWriter, r *http.Request, t string) {
	fmt.Println("starting the create page test")
	var Message struct {
    Email string
    Title string
    Time int64
    }
	defer r.Body.Close()
    if body, err := ioutil.ReadAll(r.Body); err != nil {
        fmt.Println(w, "Couldn't read request body: %s", err)
    } else {
		fmt.Println("read body:")
		fmt.Println(body)
        dec := json.NewDecoder(strings.NewReader(string(body)))
        if err := dec.Decode(&Message); err != nil {
            fmt.Println(w, "Couldn't decode JSON: %s", err)
			os.Exit(1)
        } else {
            fmt.Println(w, "Value of email is: %s", Message.Email)
        }
    }
	u := model.NewUser("",Message.Email,"")
	WikidbDAO.LookupUserNameByEmail(u)
	fmt.Println(u.User_id())
	if u.User_id() > 0{
		p := model.NewPage("test","meritwiki.com/test",u.User_id())
		fmt.Println("found user, creating page")
		WikidbDAO.CreatePage(p)
		sleepnow()
	}else{
		http.Error(w, "could not create page, bad user email", http.StatusInternalServerError)
	}
	//wait to make sure the create is complete and committed to DB
	
}
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r,"")
	}
}

func main() {
	http.HandleFunc("/createuser/", makeHandler(createUserHandler))
	http.HandleFunc("/createpage/", makeHandler(createPageHandler))
	http.ListenAndServe(":8080", nil)
}
