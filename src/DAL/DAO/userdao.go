// Copyright 2010 The go-pgsql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package WikidbDAO

import (
	"fmt"
	"os"
	"errors"
	"../../Model"
	"github.com/lxn/go-pgsql"
)

func DeleteUser(u *model.User) (error){
	conn, err := GetConnection()
	command := "select delete_user(@email);"
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	stmt, err := conn.Prepare(command, emailParam)
	defer stmt.Close()
	err = emailParam.SetValue(u.Email)
	if err != nil {
		return err
	}
	rs, err := stmt.Query()
	if err != nil  {
		err = errors.New("error deleting user")
		return err
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				err = errors.New("Error in user DAO")
				return err
			}

			if !hasRow {
				break
			}else{
				//get result from rs
				success, isNull, err := rs.Bool(0)
				if err != nil || !success || isNull{
					err = errors.New("Error in user DAO")
					return err
				}
				return nil
			}
		}
		
	}
	return err;
}
func LogonUser(inEmail string, inPW string) (bool){
	conn, err := GetConnection()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	command := "select logon_user(@email,@password)"
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	passwordParam := pgsql.NewParameter("@password", pgsql.Text)
	stmt, err := conn.Prepare(command, emailParam,passwordParam )
	defer stmt.Close()
	err = emailParam.SetValue(inEmail)
	if err != nil {
		os.Exit(1)
	}
	err = passwordParam.SetValue(inPW)
	if err != nil {
		os.Exit(1)
	}
	
	rs, err := stmt.Query()
	if err != nil {

		err = errors.New("User Logon failed in DAO")
		return false
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				err = errors.New("Invalid email or password!")
				return false
			}

			if !hasRow {
				break
			}else{
				//get result from rs
				isLoggedIn, isNull, err := rs.Bool(0)
				if err != nil ||isNull{
					os.Exit(1)
				}
				return isLoggedIn
			}
		}
		
	}
	return false;
}
func CreateUser(u *model.User) (error){
	//assume failure
	u.Resp.UserID = -1
	u.Resp.Message = "Failure Creating User"
	if (len(u.UserName) == 0 || len(u.Password) == 0 || len(u.Email) == 0 ){
		err := errors.New("Required Data Missing")
		return err
	}
	conn, err := GetConnection()
	command := "select create_wikiuser(@username, @email, @password)"
	usernameParam := pgsql.NewParameter("@username", pgsql.Text)
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	passwordParam := pgsql.NewParameter("@password", pgsql.Text)
	stmt, err := conn.Prepare(command, usernameParam, emailParam,passwordParam )
	defer stmt.Close()
	err = usernameParam.SetValue(u.UserName)
	if err != nil {
		return err
	}
	err = emailParam.SetValue(u.Email)
	if err != nil {
		return err
	}
	err = passwordParam.SetValue(u.Password)
	if err != nil {
		return err
	}
	
	rs, err := stmt.Query()
	if err != nil {
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
					os.Exit(1)
				}
				if isNull {
					user_id = 0
				}else{
					u.Resp.UserID = user_id
					u.Resp.Message = "Successfully Created User"  
				}
			}
		}
		
	}
	if u.Resp.UserID < 0{
		err = errors.New("User with that email already exists!")
		return err
	}else{
		return nil
	}
}
func LookupUserByEmail(u *model.User) (error){
	conn, err := GetConnection()	
	command := "select * from lookup_user(@email);"
	emailParam := pgsql.NewParameter("@email", pgsql.Text)
	stmt, err := conn.Prepare(command, emailParam)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	err = emailParam.SetValue(u.Email)
	if err != nil {
		return err
	}
	rs, err := stmt.Query()
	if err != nil {
		return err
	}
	defer rs.Close()
	for {
		hasRow, err := rs.FetchNext()
		if err != nil {
			return err
		}

		if !hasRow {
			break
		}else{
		    userid, isNull, err := rs.Int(0)
			if err != nil {
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

				return err
			}
			if isNull {
				user_rating = 0
			}else{
			    u.Rating = user_rating
			}
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
	
	_, err = stmt.Execute()
	if err != nil {
		fmt.Println(err)
		err = errors.New("User update failed in DAO")
		return err
	}
	return nil
}