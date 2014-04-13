// Copyright 2010 The go-pgsql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package WikidbDAO

import (
	"fmt"
	"os"
	"github.com/lxn/go-pgsql"
	"../Model"
	"errors"
)
func CreatePageSection(p *model.Section) (error){
	conn, err := pgsql.Connect("dbname=wikidb user=postgres password=lipscomb", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		os.Exit(1)
	}
	defer conn.Close()
	
	command := "select * from insert_pagesection(@page_id, @sectionid,@sectiontitle, @sectiontext, @userid) "
	pageidParam := pgsql.NewParameter("@page_id", pgsql.Integer)
	sectionidParam := pgsql.NewParameter("@sectionid", pgsql.Integer)
	titleParam := pgsql.NewParameter("@sectiontitle", pgsql.Text)
	textParam := pgsql.NewParameter("@sectiontext", pgsql.Text)
	useridParam := pgsql.NewParameter("@userid", pgsql.Integer)
	stmt, err := conn.Prepare(command,pageidParam, sectionidParam, titleParam, textParam,useridParam )
	defer stmt.Close()
	
	err = pageidParam.SetValue(p.PageID)
	if err != nil {
		fmt.Println("error in setting PageID")
	}
	err = sectionidParam.SetValue(p.SectionID)
	if err != nil {
		fmt.Println("error in setting PageID")
	}
	err = titleParam.SetValue(p.Title)
	if err != nil {
		fmt.Println("error in setting PageTitle")
	}
	err = textParam.SetValue(p.Text)
	if err != nil {
		fmt.Println("error in setting section text")
	}	
	err = useridParam.SetValue(p.CreatedByUserID)
	if err != nil {
		fmt.Println("error in setting userid")
	}
	rs,err := stmt.Query()
	if err != nil {
		fmt.Println("error in creating PageSection")
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				err = errors.New("PageSection Creation failed")
				return err
			}

			if !hasRow {
				break
			}else{
				//get revisionid from rs
				page_revision, isNull, err := rs.Int(0)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if isNull {
					page_revision = 0
				}else{
					fmt.Println("set the revision to:",page_revision)
					p.LatestRevisionID = page_revision
					fmt.Println("revision set to:",p.LatestRevisionID)
				}
			}
		}
	}
	return nil
}
func UpdatePageSection(p *model.Section) (error){
	conn, err := pgsql.Connect("dbname=wikidb user=postgres password=lipscomb", pgsql.LogError)
	if err != nil {
		fmt.Println("error in connect")
		os.Exit(1)
	}
	defer conn.Close()
	
	command := "select update_section(@page_id, @sectionid, @sectiontitle, @sectiontext, @userid) "
	pageidParam := pgsql.NewParameter("@page_id", pgsql.Integer)
	sectionidParam := pgsql.NewParameter("@sectionid", pgsql.Integer)
	titleParam := pgsql.NewParameter("@sectiontitle", pgsql.Text)
	textParam := pgsql.NewParameter("@sectiontext", pgsql.Text)
	useridParam := pgsql.NewParameter("@userid", pgsql.Integer)
	
	pageidParam.SetValue(p.PageID)
	sectionidParam.SetValue(p.SectionID)
	titleParam.SetValue(p.Title)
	textParam.SetValue(p.Text)
	useridParam.SetValue(p.LastUpdateByUserID)
	
	stmt, err := conn.Prepare(command,pageidParam, sectionidParam, titleParam, textParam,useridParam )
	defer stmt.Close()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}else{
		fmt.Println("stmt",stmt)
	}
	rowsaffected,err := stmt.Execute()
	if err != nil {
		fmt.Println(rowsaffected)
		err = errors.New("error in updating page")
		return err
	}	
	return nil
}		
