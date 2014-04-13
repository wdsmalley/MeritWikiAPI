package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"encoding/json"
	//"errors"
	//"time"
	//"net/url"

	"./DAL"
	"./Model"
)

func main() {
	http.HandleFunc("/help", Help_PageHandler)
	
	http.HandleFunc("/authenticate", Authenticate_PostHandler)
	
	http.HandleFunc("/user/create", User_Create_PostHandler)
	http.HandleFunc("/user/update", User_Create_PostHandler)
	http.HandleFunc("/user/get", User_Get_Handler)
	
	http.HandleFunc("/page/create", Page_Create_PostHandler)
	http.HandleFunc("/page/update", Page_Create_PostHandler)
	http.HandleFunc("/page/get", Page_Get_Handler)

	http.HandleFunc("/section/create", Section_Create_PostHandler)
	http.HandleFunc("/section/update", Section_Update_PostHandler)
	http.HandleFunc("/section/get", Section_Get_Handler)
	http.HandleFunc("/section/getrevisions", Section_GetRevisions_Handler)
	
	http.ListenAndServe(":8080", nil)
}

func Authenticate_PostHandler(w http.ResponseWriter, r * http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	
	user, err1 := DAL.Authenticate(username, password)

	if err1 == nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:user, Error:""}).ToJsonString()) 
	} else {
		fmt.Fprintf(w, (&model.Response{ResponseObject:user, Error:err1.Error()}).ToJsonString())
	}
}

func Help_PageHandler(w http.ResponseWriter, r * http.Request) {
	body, _ := ioutil.ReadFile("help.html")
	fmt.Fprintf(w, string(body))
}

func User_Create_PostHandler(w http.ResponseWriter, r * http.Request) {
	if r.Method == "POST" {
		var user * model.User
		
		decoder := json.NewDecoder(r.Body)
		err1 := decoder.Decode(&user)
		
		if err1 == nil {
			retval, err2 := DAL.CreateUser(user)

			if err2 == nil {
				fmt.Fprintf(w, (&model.Response{ResponseObject:retval, Error:""}).ToJsonString()) 
			} else {
				fmt.Fprintf(w, (&model.Response{ResponseObject:retval, Error:err2.Error()}).ToJsonString())
			}
		} else {
			fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:"Malformed Json Data"}).ToJsonString())
		}
	}
}

func User_Update_PostHandler(w http.ResponseWriter, r * http.Request) {
	if r.Method == "POST" {
		
	}
}

func User_Get_Handler(w http.ResponseWriter, r * http.Request) {
	id := r.FormValue("id")
	
	user_id, err0 := strconv.Atoi(id)
	
	if err0 != nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:"Invalid User ID"}).ToJsonString()) 
	}
	
	user, err1 := DAL.GetUserByID(user_id)

	if err1 == nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:user, Error:""}).ToJsonString()) 
	} else {
		fmt.Fprintf(w, (&model.Response{ResponseObject:user, Error:err1.Error()}).ToJsonString())
	}
}

func Page_Create_PostHandler(w http.ResponseWriter, r * http.Request) {
	if r.Method == "POST" {
		var page * model.Page
		decoder := json.NewDecoder(r.Body)
		err1 := decoder.Decode(&page)
		if err1 == nil {
			err2 := DAL.CreatePage(page)

			if err2 == nil {
				fmt.Fprintf(w, (&model.Response{ResponseObject:page, Error:"Page creation succesful"}).ToJsonString()) 
			} else {
				http.Error(w, err2.Error(), 412)
				fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:err2.Error()}).ToJsonString())
			}
		} else {
			fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:"Malformed Json Data"}).ToJsonString())
		}
	}
}

func Page_Update_PostHandler(w http.ResponseWriter, r * http.Request) {
	if r.Method == "POST" {
		
	}
}

func Page_Get_Handler(w http.ResponseWriter, r * http.Request) {
	id := r.FormValue("id")
	
	page_id, err0 := strconv.Atoi(id)
	
	page, err1 := DAL.GetPageByID(page_id)

	if err0 != nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:"Invalid Page ID"}).ToJsonString()) 
	}
	
	if err1 == nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:page, Error:""}).ToJsonString()) 
	} else {
		fmt.Fprintf(w, (&model.Response{ResponseObject:page, Error:err1.Error()}).ToJsonString())
	}
}

func Section_Create_PostHandler(w http.ResponseWriter, r * http.Request) {
	if r.Method == "POST" {
		
	}
}

func Section_Update_PostHandler(w http.ResponseWriter, r * http.Request) {
	if r.Method == "POST" {
		
	}
}

func Section_Get_Handler(w http.ResponseWriter, r * http.Request) {
	id := r.FormValue("id")
	
	section_id, err0 := strconv.Atoi(id)
	
	section, err1 := DAL.GetSectionByID(section_id)

	if err0 != nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:"Invalid Section ID"}).ToJsonString()) 
	}
	
	if err1 == nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:section, Error:""}).ToJsonString()) 
	} else {
		fmt.Fprintf(w, (&model.Response{ResponseObject:section, Error:err1.Error()}).ToJsonString())
	}
}

func Section_GetRevisions_Handler(w http.ResponseWriter, r * http.Request) {
	id := r.FormValue("id")
	
	section_id, err0 := strconv.Atoi(id)
	
	if err0 != nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:"Invalid Section ID"}).ToJsonString()) 
	}
	
	revisions, err1 := DAL.GetSectionRevisionsByID(section_id)

	if err1 == nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:revisions, Error:""}).ToJsonString()) 
	} else {
		fmt.Fprintf(w, (&model.Response{ResponseObject:revisions, Error:err1.Error()}).ToJsonString())
	}
}