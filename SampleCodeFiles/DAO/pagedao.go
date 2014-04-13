// Copyright 2010 The go-pgsql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package WikidbDAO

import (
	"fmt"
	"os"
	"errors"
	"github.com/lxn/go-pgsql"
	"../Model"
)
func DeletePage(iPageID int){
conn, err := pgsql.Connect("dbname=wikidb user=postgres password=lipscomb", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		os.Exit(1)
	}
	defer conn.Close()
	command := "select delete_page(@pageid);"
	pageidParam := pgsql.NewParameter("@pageid", pgsql.Integer)
	stmt, err := conn.Prepare(command, pageidParam)
	defer stmt.Close()
	err = pageidParam.SetValue(iPageID)
	if err != nil {
		os.Exit(1)
	}
	
	rs, err := stmt.Query()
	if err != nil {
		fmt.Println("error in query")
		fmt.Println(rs)
		os.Exit(1)
	}
	fmt.Println("successful delete of Page")
}
func CreatePage(p *model.Page) (error){
	conn, err := pgsql.Connect("dbname=wikidb user=postgres password=lipscomb", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		os.Exit(1)
	}
	defer conn.Close()
	command := "select * from create_page(@title, @url, @userid)"
	titleParam := pgsql.NewParameter("@title", pgsql.Text)
	urlParam := pgsql.NewParameter("@url", pgsql.Text)
	useridParam := pgsql.NewParameter("@userid", pgsql.Integer)
	stmt, err := conn.Prepare(command, titleParam, urlParam,useridParam )
	defer stmt.Close()
	err = titleParam.SetValue(p.Title)
	if err != nil {
		os.Exit(1)
	}
	err = urlParam.SetValue(p.Url)
	if err != nil {
		os.Exit(1)
	}
	err = useridParam.SetValue(p.Created_by_user)
	if err != nil {
		os.Exit(1)
	}
	
	rs, err := stmt.Query()
	if err != nil {
		fmt.Println("error in creating page")
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				err = errors.New("Page Creation failed")
				return err
			}

			if !hasRow {
				break
			}else{
				//get page_id from rs
				page_id, isNull, err := rs.Int(0)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if isNull {
					page_id = 0
				}else{
					p.PageID = page_id
				}
			}
		}
	}
	return nil
}
func UpdatePage(p *model.Page) (error){
	conn, err := pgsql.Connect("dbname=wikidb user=postgres password=lipscomb", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		os.Exit(1)
	}
	defer conn.Close()
	command := "select update_page(@pageid, @title, @userid)"
	pageidParam := pgsql.NewParameter("@pageid", pgsql.Integer)
	titleParam := pgsql.NewParameter("@title", pgsql.Text)
	useridParam := pgsql.NewParameter("@userid", pgsql.Integer)
	
	err = pageidParam.SetValue(p.PageID)
	if err != nil {
		os.Exit(1)
	}
	err = titleParam.SetValue(p.Title)
	if err != nil {
		os.Exit(1)
	}
	err = useridParam.SetValue(p.Created_by_user)
	if err != nil {
		os.Exit(1)
	}
	stmt, err := conn.Prepare(command, pageidParam,titleParam, useridParam )
	defer stmt.Close()
	rowsaffected, err := stmt.Execute()
	if err != nil {
		fmt.Println(rowsaffected)
		err = errors.New("error in updating page")
		return err
	}
	return nil
	
}
func LookupPageByTitle(p *model.Page){
	conn, err := pgsql.Connect("dbname=wikidb user=postgres", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		os.Exit(1)
	}
	defer conn.Close()
	command := "select * from lookup_page_bytitle(@title);"
	titleParam := pgsql.NewParameter("@title", pgsql.Text)
	stmt, err := conn.Prepare(command, titleParam)
	if err != nil {
		fmt.Println("error in prepare")
		os.Exit(1)
	}
	defer stmt.Close()
	
	err = titleParam.SetValue(p.Title)
	
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("running the query")
	rs, err := stmt.Query()
	if err != nil {
		fmt.Println("error in query: " )
		fmt.Println(rs)
		os.Exit(1)
	}
	defer rs.Close()

	for {
		hasRow, err := rs.FetchNext()
		if err != nil {
			fmt.Println("error in fetch:")
			os.Exit(1)
		}

		if !hasRow {
			break
		}else{
			//get page_id from rs
			page_id, isNull, err := rs.Int(0)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if isNull {
				page_id = 0
			}else{
			    p.PageID = page_id
			}
			
			//get page title from rs
			page_title, isNull, err := rs.String(1)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if isNull {
				page_title = "(null)"
			}else{
			    p.Title = page_title
			}
			
			//get url from rs
			url, isNull, err := rs.String(2)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if isNull {
				url = "(null)"
			}else{
			    p.Url = url
			}
			//get revisionid from rs
			revisionid, isNull, err := rs.Int(3)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if isNull {
				revisionid = 0
			}else{
			    p.RevisionID = revisionid
			}
			//get created_by_username from rs
			created_by_username, isNull, err := rs.String(4)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if isNull {
				created_by_username = "(null)"
			}else{
			    p.Created_by_username = created_by_username
			}
		}

	}
}		
