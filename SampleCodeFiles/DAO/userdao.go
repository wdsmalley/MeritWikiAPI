// Copyright 2010 The go-pgsql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package WikidbDAO

import (
	"fmt"
	"os"
	"errors"
	"../Model"
	"github.com/lxn/go-pgsql"
)

func DeleteUser(u *model.User){
	conn, err := GetConnection()
	command := "select delete_user(@username);"
	usernameParam := pgsql.NewParameter("@username", pgsql.Text)
	stmt, err := conn.Prepare(command, usernameParam)
	defer stmt.Close()
	err = usernameParam.SetValue(u.UserName)
	if err != nil {
		os.Exit(1)
	}
	
	rs, err := stmt.Query()
	if err != nil {
		fmt.Println("error in query")
		fmt.Println(rs)
		os.Exit(1)
	}
	fmt.Println("successful delete of User")
}
func CreateUser(u *model.User) (error){
	conn, err := GetConnection()
	command := "select create_wikiuser(@username, @email, @password)"
	usernameParam := pgsql.NewParameter("@username", pgsql.Text)
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	passwordParam := pgsql.NewParameter("@password", pgsql.Text)
	stmt, err := conn.Prepare(command, usernameParam, emailParam,passwordParam )
	defer stmt.Close()
	err = usernameParam.SetValue(u.UserName)
	if err != nil {
		os.Exit(1)
	}
	err = emailParam.SetValue(u.Email)
	if err != nil {
		os.Exit(1)
	}
	err = passwordParam.SetValue(u.Password)
	if err != nil {
		os.Exit(1)
	}
	
	rs, err := stmt.Query()
	fmt.Println(err)
	if err != nil {
		fmt.Println("error in query")
		fmt.Println(rs)
		err = errors.New("User Creation failed in DAO")
		return err
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				err = errors.New("User with that email already exists!")
				return err
			}

			if !hasRow {
				break
			}else{
				//get user_id from rs
				user_id, isNull, err := rs.Int(0)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if isNull {
					user_id = 0
				}else{
					u.UserID =  user_id
				}
			}
		}
		
	}
	return nil
}
func LookupUserByEmail(u *model.User) (error){
	conn, err := GetConnection()	
	command := "select * from lookup_user(@email);"
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	stmt, err := conn.Prepare(command, emailParam)
	if err != nil {
		fmt.Println("error in prepare")
		return err
	}
	defer stmt.Close()
	
	err = emailParam.SetValue(u.Email)
	if err != nil {
		return err
	}
	fmt.Println("running the query")
	rs, err := stmt.Query()
	if err != nil {
		fmt.Println("error in query: " )
		fmt.Println(rs)
		return err
	}
	defer rs.Close()
	for {
		hasRow, err := rs.FetchNext()
		if err != nil {
			fmt.Println("error in fetch:")
			return err
		}

		if !hasRow {
			break
		}else{
		    userid, isNull, err := rs.Int(0)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if isNull {
				userid = 0
			}else{
			    u.UserID = userid
			}
			//////
			username, isNull, err := rs.String(1)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if isNull {
				username = "(null)"
			}else{
			    u.UserName = username
			}
			//////
			user_email, isNull, err := rs.String(2)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if isNull {
				user_email = "(null)"
			}else{
			    u.Email = user_email 
			}
			//////
			user_status, isNull, err := rs.Int(3)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if isNull {
				user_status = 0
			}else{
			    u.Status = user_status
			}
			//////
			user_rating, isNull, err := rs.Int(4)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if isNull {
				user_rating = 0
			}else{
			    u.Rating = user_rating
			}
		   fmt.Println("userid: ", u.UserID)	
		   fmt.Println("username: ", u.UserName)
		   fmt.Println("useremail: ", u.Email)
		   fmt.Println("userstatus: ", u.Status)
		   fmt.Println("userrating: ", u.Rating)
		}

	}
	return nil
}		
func UpdateUser(u *model.User) (error){
	conn, err := GetConnection()
	command := "select update_wikiuser(@user_id, @username, @email, @password)"
	useridParam := pgsql.NewParameter("@user_id", pgsql.Integer)
	usernameParam := pgsql.NewParameter("@username", pgsql.Text)
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	passwordParam := pgsql.NewParameter("@password", pgsql.Text)
	stmt, err := conn.Prepare(command, useridParam, usernameParam, emailParam,passwordParam )
	defer stmt.Close()
	fmt.Println(u.UserName)
	fmt.Println(u.Email)
	err = useridParam.SetValue(u.UserID)
	if err != nil {
		os.Exit(1)
	}
	err = usernameParam.SetValue(u.UserName)
	if err != nil {
		os.Exit(1)
	}
	err = emailParam.SetValue(u.Email)
	if err != nil {
		os.Exit(1)
	}
	err = passwordParam.SetValue(u.Password)
	if err != nil {
		os.Exit(1)
	}
	
	rowsaffected, err := stmt.Execute()
	if err != nil {
		fmt.Println(rowsaffected)
		fmt.Println("error in query")
		err = errors.New("User Creation failed in DAO")
		return err
	}
	return nil
}