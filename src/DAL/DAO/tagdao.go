// Copyright 2010 The go-pgsql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package WikidbDAO

import (
	"errors"
	"../../Model"
	"github.com/lxn/go-pgsql"
)
func GetAllTags() ([]model.Tag){
	tags := []model.Tag{}
	conn, err := GetConnection()
	command := "select * from get_all_tags()"
	stmt, err := conn.Prepare(command)
	defer stmt.Close()
	rs, err := stmt.Query()
	if err != nil {
		return nil
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				return nil
			}

			if !hasRow {
				break
			}else{
				t := model.Tag{}
				//get tag from rs
				tag_id, isNull, err := rs.Int(0)
				if err != nil {
					return nil
				}
				if isNull {
					tag_id = 0
				}else{
					t.TagID = tag_id
				}
				
				tag_name, isNull, err := rs.String(1)
				if err != nil {
					return nil
				}
				if isNull {
					tag_name = ""
				}else{
					t.TagName = tag_name
				}
				
				tag_desc, isNull, err := rs.String(2)
				if err != nil {
					return nil
				}
				if isNull {
					tag_desc = ""
				}else{
					t.TagDesc = tag_desc
				}
				tags = model.AppendTag(tags,t)
			}
		}
		
	}

	return tags;
}
func CreateTag(t *model.Tag) (error){
	//assume failure
	t.TagID = -1
	if (len(t.TagName) == 0 || len(t.TagDesc) == 0){
		err := errors.New("Missing parameters")
		return err
	}
	conn, err := GetConnection()
	command := "select create_tag(@tagname, @desc, @userid)"
	tagnameParam := pgsql.NewParameter("@tagname", pgsql.Text)
	descParam := pgsql.NewParameter("@desc", pgsql.Text)
	useridParam := pgsql.NewParameter("@userid", pgsql.Integer)
	stmt, err := conn.Prepare(command, tagnameParam, descParam,useridParam )
	defer stmt.Close()
	err = tagnameParam.SetValue(t.TagName)
	if err != nil {
		return err
	}
	err = descParam.SetValue(t.TagDesc)
	if err != nil {
		return err
	}
	err = useridParam.SetValue(t.UserID)
	if err != nil {
		return err
	}
	
	rs, err := stmt.Query()
	if err != nil {
		err = errors.New("Tag Creation failed in DAO")
		return err
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				err = errors.New("Tag with that name already exists!")
				return err
			}

			if !hasRow {
				break
			}else{
				//get tag_id from rs
				tag_id, isNull, err := rs.Int(0)
				if err != nil {
					return err
				}
				if isNull {
					tag_id = 0
				}else{
					t.TagID = tag_id
				}
			}
		}
		
	}
	if t.TagID < 0{
		err = errors.New("User with that email already exists!")
		return err
	}else{
		return nil
	}
}
func DeleteTag(name string) (bool){
	//assume failure
	result := false
	if len(name) == 0 {
		return result	
	}
	conn, err := GetConnection()
	command := "select delete_tag(@tagname)"
	tagnameParam := pgsql.NewParameter("@tagname", pgsql.Text)
	stmt, err := conn.Prepare(command, tagnameParam)
	defer stmt.Close()
	err = tagnameParam.SetValue(name)
	if err != nil {
		return result
	}
	rs, err := stmt.Query()
	if err != nil {
		return false
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				return false
			}
			if !hasRow {
				break
			}else{
				//get result from rs
				result, isNull, err := rs.Bool(0)
				if err != nil {
					return false
				}
				if isNull {
					return  false
				}else{
					return result
				}
			}
		}
		
	}
	return false
}