package api

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"strconv"
	"encoding/json"
	//"errors"
	//"time"
	//"net/url"

	"../DAL"
	"../Model"
)

func StartAPI() {
	
	http.HandleFunc("/help", Help_PageHandler)

	http.HandleFunc("/authenticate", Authenticate_PostHandler)

	http.HandleFunc("/user/create", User_Create_PostHandler)
	http.HandleFunc("/user/update", User_Update_PostHandler)
	http.HandleFunc("/user/get", User_Get_Handler)

	http.HandleFunc("/page/create", Page_Create_PostHandler)
	http.HandleFunc("/page/update", Page_Update_PostHandler)
	http.HandleFunc("/page/get", Page_Get_Handler)

	http.HandleFunc("/section/create", Section_Create_PostHandler)
	http.HandleFunc("/section/update", Section_Update_PostHandler)
	http.HandleFunc("/section/get", Section_Get_Handler)
	http.HandleFunc("/section/getrevisions", Section_GetRevisions_Handler)

	http.ListenAndServe(":8080",nil)
}

func Authenticate_PostHandler(w http.ResponseWriter, r * http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	
	isValid := DAL.Authenticate(username, password)

	if isValid {
		fmt.Fprintf(w, (&model.Response{ResponseObject:isValid, Error:""}).ToJsonString()) 
	} else {
		fmt.Fprintf(w, (&model.Response{ResponseObject:isValid, Error:"Invalid username or password"}).ToJsonString())
	}
}

func Help_PageHandler(w http.ResponseWriter, r * http.Request) {
	body, _ := ioutil.ReadFile("help.html")
	fmt.Fprintf(w, string(body))
}

func User_Create_PostHandler(w http.ResponseWriter, r * http.Request) {
	u := model.User{}
	resp := &u.Resp
	PostHandler(w, r, &u, &resp, func(obj interface{}) (error) { return DAL.CreateUser(obj.(*model.User)) })
}

func User_Update_PostHandler(w http.ResponseWriter, r * http.Request) {
	u := model.User{}
	resp := &u.Resp
	PostHandler(w, r, &u,&resp, func(obj interface{}) (error) { return DAL.UpdateUser(obj.(*model.User)) })
}

func User_Get_Handler(w http.ResponseWriter, r * http.Request) {
	GetByIDHandler(w, r, func(id int) (interface{}, error) { return DAL.GetUserByID(id) })
}

func Page_Create_PostHandler(w http.ResponseWriter, r * http.Request) {
	PostHandler(w, r, &model.Page{},&model.Page{}, func(obj interface{}) (error) { return DAL.CreatePage(obj.(*model.Page)) })
}

func Page_Update_PostHandler(w http.ResponseWriter, r * http.Request) {
	PostHandler(w, r, &model.Page{},&model.Page{}, func(obj interface{}) (error) { return DAL.UpdatePage(obj.(*model.Page)) })
}

func Page_Get_Handler(w http.ResponseWriter, r * http.Request) {
	GetByIDHandler(w, r, func(id int) (interface{}, error) { return DAL.GetPageByID(id) })
}

func Section_Create_PostHandler(w http.ResponseWriter, r * http.Request) {
	PostHandler(w, r, &model.Section{},&model.Section{}, func(obj interface{}) (error) { return DAL.CreateSection(obj.(*model.Section)) })
}

func Section_Update_PostHandler(w http.ResponseWriter, r * http.Request) {
	PostHandler(w, r, &[]model.Section{},&[]model.Section{}, func(obj interface{}) (error) { return DAL.UpdateSection(obj.([]model.Section)) })
}

func Section_Get_Handler(w http.ResponseWriter, r * http.Request) {	
	GetByIDHandler(w, r, func(id int) (interface{}, error) { return DAL.GetSectionByID(id) })
}

func Section_GetRevisions_Handler(w http.ResponseWriter, r * http.Request) {
	GetByIDHandler(w, r, func(id int) (interface{}, error) { return DAL.GetSectionRevisionsByID(id) })
}

type PostDelegate func(interface{}) error

func PostHandler(w http.ResponseWriter, r * http.Request, req interface{},resp interface{}, del PostDelegate) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		err1 := decoder.Decode(&req)
		if err1 == nil {
			err2 := del(req);
			if err2 == nil {
				fmt.Fprintf(w, (&model.Response{ResponseObject:resp, Error:""}).ToJsonString()) 
				
			} else {
				fmt.Println("error",err2)
				fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:err2.Error()}).ToJsonString())
			}
		} else {
			fmt.Println("error",err1)
			os.Exit(1)
			fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:"Malformed Json Data"}).ToJsonString())
		}	
	}
}

type GetByIDDelegate func(int) (interface{}, error)

func GetByIDHandler(w http.ResponseWriter, r * http.Request, del GetByIDDelegate) {
	temp := r.FormValue("id")
	
	id, err0 := strconv.Atoi(temp)

	if err0 != nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:nil, Error:"Invalid ID"}).ToJsonString()) 
		return
	}
	
	obj, err1 := del(id)
	
	if err1 == nil {
		fmt.Fprintf(w, (&model.Response{ResponseObject:obj, Error:""}).ToJsonString()) 
	} else {
		fmt.Fprintf(w, (&model.Response{ResponseObject:obj, Error:err1.Error()}).ToJsonString())
	}
}
