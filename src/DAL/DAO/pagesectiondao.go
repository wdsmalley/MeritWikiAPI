// Copyright 2010 The go-pgsql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package WikidbDAO

import (
	"fmt"
	"os"
	"github.com/lxn/go-pgsql"
	"../../Model"
	"errors"
)
func CreatePageSection(p *model.Section) (error){
	conn, err := GetConnection()
	if err != nil {
		fmt.Println(err)
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
		err = errors.New("PageSection Creation failed")
		return err
	}
	err = sectionidParam.SetValue(p.SectionID)
	if err != nil {
		err = errors.New("PageSection Creation failed")
		return err
	}
	err = titleParam.SetValue(p.Title)
	if err != nil {
		err = errors.New("PageSection Creation failed")
		return err
	}
	err = textParam.SetValue(p.Text)
	if err != nil {
		err = errors.New("PageSection Creation failed")
		return err
	}
	err = useridParam.SetValue(p.CreatedByUserID)
	if err != nil {
		err = errors.New("PageSection Creation failed")
		return err
	}
	rs,err := stmt.Query()
	if err != nil {
		err = errors.New("PageSection Creation failed")
		return err
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
					return(err)
				}
				if isNull {
					page_revision = 0
				}else{
					p.LatestRevisionID = page_revision
				}
			}
		}
	}
	return nil
}
func UpdatePageSection(p *model.Section) (error){
	conn, err := GetConnection()
	if err != nil {
		fmt.Println(err)
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
		return err
	}
	rowsaffected,err := stmt.Execute()
	if err != nil || rowsaffected < 1{
		err = errors.New("error in updating page")
		return err
	}	
	return nil
}		
