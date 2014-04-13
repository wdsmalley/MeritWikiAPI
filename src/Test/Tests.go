package main

import (
//	"net/http"
	"fmt"
	"os"
//	"io/ioutil"
	"../DAL"
	"../Model"
	"../API"
	"../DAL/DAO"	
)

func main() {
	fmt.Println("Starting API")
	go api.StartAPI()
	fmt.Println("API Started")
	testversion := GetCurrentVersion()
	fmt.Println("Starting Tests",testversion)
	DAOTests(testversion)
	DALTests(testversion)
	TestAPI()
}
func DAOTests(testversion int){
	//DAO Tests
	fmt.Println("DAO tests")
	fmt.Println("")	

	if testversion >= 1{
		fmt.Println("Authentication tests")
		TestAuthenticateDAO("verify@login.com","test",true)						//should return success
		TestAuthenticateDAO("test@gmail.com","test",false)  					//should return error (does not exist)
		TestAuthenticateDAO("","test",false)  									//should return error (no email)
		TestAuthenticateDAO("test@gmail.com","",false)  						//should return error (no pw)
		TestAuthenticateDAO("","",false)  										//should return error (no email/pw)
		

		fmt.Println("")	
		fmt.Println("Create user tests")
		user := SetupUserObjectForTest("daotest", "daopw", "daotest@mail.com")
		TestCreateUserDAO(&user,true)											//should return success
		
		user = SetupUserObjectForTest("failtest", "fail", "verify@login.com")
		TestCreateUserDAO(&user,false)											//should return error (duplicate)

		user = SetupUserObjectForTest("failtest", "", "verify@login.com")
		TestCreateUserDAO(&user,false)											//should return error (no pw)

		user = SetupUserObjectForTest("", "fail", "verify@login.com")
		TestCreateUserDAO(&user,false)											//should return error (no username)

		user = SetupUserObjectForTest("failtest", "fail", "")
		TestCreateUserDAO(&user,false)											//should return error (no email)

		
		fmt.Println("")
		fmt.Println("Create page tests")
		var sections []model.Section
		sections = model.Append(sections,SetupSectionObjectForTest("daotestsection1","test section 1 text",1))
		sections = model.Append(sections,SetupSectionObjectForTest("daotestsection2","test section 2 text",1))
		page := SetupPageObjectForTest("daotest1",-34,sections)
		TestCreatePageDAO(&page,false)											//should return failure (no such user)
		page = SetupPageObjectForTest("daotest2",2,sections)
		TestCreatePageDAO(&page,true)											//should return success
	}
	if testversion >= 2{
		fmt.Println("")
		fmt.Println("Create tag tests")
		tag := SetupTagObjectForTest("newtagtest", "blah blah blah", -34)
		TestCreateTagDAO(&tag,false)											//should return error (no such user)
		
		tag = SetupTagObjectForTest("newtagtest", "blah blah blah", 2)
		TestCreateTagDAO(&tag,true)												//should return success	

		tag = SetupTagObjectForTest("oldtag", "blah blah blah", 3)
		TestCreateTagDAO(&tag,false)											//should return error (duplicate name)

		tag = SetupTagObjectForTest("", "blah blah blah", 3)
		TestCreateTagDAO(&tag,false)											//should return error (no name)

		tag = SetupTagObjectForTest("oldtag", "", 3)
		TestCreateTagDAO(&tag,false)											//should return error (no description)

		fmt.Println("")	
		fmt.Println("Delete user tests")
		user := SetupUserObjectForTest("", "", "delete@me.com")
		TestDeleteUserDAO(&user,false)											//should return error (does not exist)
		user = SetupUserObjectForTest("", "test1", "delete@mail.com")
		TestDeleteUserDAO(&user,true)											//should return success

		fmt.Println("")	
		fmt.Println("Delete tag tests")
		TestDeleteTagDAO("oldtag",true)											//should return success
		TestDeleteTagDAO("deletetag",false)										//should return error (does not exist)
		TestDeleteTagDAO("",false)												//should return error (no name)

		fmt.Println("")	
		fmt.Println("get all tags test")
		TestGetAllTagsDAO(true)													//should return array of tags
	}
}
func DALTests(testversion int){

	//DAL tests
	//User functions
	fmt.Println("")
	fmt.Println("DAL tests")
	if testversion >= 1{
		fmt.Println("")
		fmt.Println("Authentication tests")
		TestAuthenticateDAL("verify@login.com","test",true)						//should return success
		TestAuthenticateDAL("test@gmail.com","test",false)  					//should return error (does not exist)
		TestAuthenticateDAL("","test",false)  									//should return error (no email)
		TestAuthenticateDAL("test@gmail.com","",false)  						//should return error (no pw)
		TestAuthenticateDAL("","",false)  										//should return error (no email/pw)
		
		fmt.Println("")
		fmt.Println("Create user tests")
		TestCreateUserDAL("newtest","newtest@mail.com","test",true)				//should return success
		TestCreateUserDAL("failtest","verify@login.com","test",false)			//should return failure	(duplicate email)
		TestCreateUserDAL("failtest", "", "verify@login.com", false)			//should return error (no pw)
		TestCreateUserDAL("", "fail", "verify@login.com", false)				//should return error (no username)
		TestCreateUserDAL("failtest", "fail", "", false)						//should return error (no email)

		
		fmt.Println("")
		fmt.Println("update user tests")
		TestUpdateUserDAL(1,"DALUpdate1","","",true)							//should return success
		TestUpdateUserDAL(2,"","update@mail1.com","",true)						//should return success
		TestUpdateUserDAL(3,"","","update3",true)								//should return success
		TestUpdateUserDAL(-1,"newtest","newtest@mail.com","testupdate",false)	//should return error
		TestGetUserByID()
		
		//Page functions
		var sections []model.Section
		sections = model.Append(sections,SetupSectionObjectForTest("daltestsection1","test section 1 text",1))
		sections = model.Append(sections,SetupSectionObjectForTest("daltestsection2","test section 2 text",1))
		
		fmt.Println("")
		fmt.Println("create page tests")
		TestCreatePageDAL("testpage",-1,sections,false)								//should return error, user doesnt exist
		TestCreatePageDAL("testpage",2,sections,true)								//should return success
	}
	if testversion >= 2{
		fmt.Println("")
		fmt.Println("create tag tests")
		TestCreateTagDAL("newtagtest", "blah blah blah", -34,false)				//should return error (no such user)
		TestCreateTagDAL("DALtagtest", "blah blah blah", 2,true)				//should return success	
		TestCreateTagDAL("oldtag", "blah blah blah", 3,false)					//should return error (duplicate name)
		TestCreateTagDAL("", "blah blah blah", 3,false)							//should return error (no name)
		TestCreateTagDAL("newdaltag", "", 3,false)								//should return error (no description)
		
		fmt.Println("")
		fmt.Println("delete tag tests")
		TestDeleteTagDAL("daltag",true)											//should return success
		TestDeleteTagDAL("deletetag",false)										//should return error (does not exist)
		TestDeleteTagDAL("",false)												//should return error (no name)
		
		fmt.Println("")
		fmt.Println("get all tags test")
		TestGetAllTagsDAL(true)													//should return tags

	}
}
func TestAPI() {

	TestUpdatePage()
	TestGetPageByID()
	
	TestCreateSection()
	TestUpdateSection()
	TestGetSectionByID()
	TestGetRevisions()
}

//testing setup functions
func SetupUserObjectForTest(username string, password string, email string) model.User{
	user := model.User{}
	user.UserName = username
	user.Password = password
	user.Email = email
	return user
}
func SetupTagObjectForTest(tagname string, tagdesc string, taguser int) model.Tag{
	tag := model.Tag{}
	tag.TagName = tagname
	tag.TagDesc = tagdesc
	tag.UserID = taguser
	return tag
}
func SetupPageObjectForTest(title string, userid int, sections []model.Section) model.Page{
	page := model.Page{}
	page.Title = title
	page.Created_by_user = userid
	page.Sections = sections
	return page
}
func SetupSectionObjectForTest(title string, text string, userid int) model.Section{
	section := model.Section{}
	section.Title = title
	section.Text = text
	section.CreatedByUserID = userid
	return section
}
//*****************************************************************
//dao test functions
func TestAuthenticateDAO(email string, pw string,expected bool) {
	r := WikidbDAO.LogonUser(email,pw) 	
	if r == expected{
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
}

func TestCreateUserDAO(user * model.User,expected bool){
	WikidbDAO.CreateUser(user)
	if (user.Resp.UserID  >= 1 && expected) || (user.Resp.UserID < 1 && !expected)   {
		fmt.Println("test passed: ", user.Email)
	}else{
		fmt.Println("test failed: ",user.Email)
		os.Exit(1)
	}
}
func TestDeleteUserDAO(user * model.User,expected bool){
	err := WikidbDAO.DeleteUser(user)
	if (err == nil && expected) || (err != nil && !expected)   {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
}
func TestCreatePageDAO(page * model.Page,expected bool){
	WikidbDAO.CreatePage(page)
	if (page.PageID >= 1 && expected) || (page.PageID < 1 && !expected)   {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
}

func TestCreateTagDAO(tag * model.Tag,expected bool){
	WikidbDAO.CreateTag(tag)
	if (tag.TagID  >= 1 && expected) || (tag.TagID < 1 && !expected)   {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
}
func TestDeleteTagDAO(name string,expected bool){
	result := WikidbDAO.DeleteTag(name)
	if (result == expected) {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
}
func TestGetAllTagsDAO(expected bool){
	tags := WikidbDAO.GetAllTags()
	if (len(tags) >0 && expected) || (tags == nil && !expected) {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}	
}
//**********************************************************************
//DAL test functions
//test the DAL functionality to authenicate a uset by email and password
func TestAuthenticateDAL(email string, pw string,expected bool) {
	r := DAL.Authenticate(email,pw)
	
	if r == expected{
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
	
	
}

//test the DAL functionalty to create a user
func TestCreateUserDAL(username string, email string, pw string, expected bool) {
	user := model.User{}
	user.UserName = username
	user.Password = pw
	user.Email = email
	err := DAL.CreateUser(&user)
	if (expected == true && err == nil && user.Resp.UserID > 0) || (expected == false && err != nil && user.Resp.UserID <= 0)  {
	//if user.Resp.UserID > 0 || err == nil {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
}
 
//test the dal functionality to update a user (username, password, or email) by id
func TestUpdateUserDAL(userid int, username string, email string, pw string, expected bool) {
	user := model.User{}
	user.UserID = userid
	user.UserName = username
	user.Password = pw
	user.Email = email
	err := DAL.UpdateUser(&user)
	if (expected == true && err == nil) || (expected == false && err != nil)  {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}	
}

func TestGetUserByID() {
	
}

func TestCreateTagDAL(name string, desc string, userid int,expected bool){
	tag := model.Tag{}
	tag.TagName = name
	tag.TagDesc = desc
	tag.UserID = userid
	err := DAL.CreateTag(&tag)
	if (expected == true && err == nil) || (expected == false && err != nil)  {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}	
}
func TestDeleteTagDAL(name string,expected bool){
	result := DAL.DeleteTag(name)
	if (result == expected) {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
}
func TestGetAllTagsDAL(expected bool){
	tags := DAL.GetAllTags()
	if (tags != nil && expected) || (tags == nil && !expected)  {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}
}
func TestCreatePageDAL(title string, userid int,sections []model.Section, expected bool) {
	page := model.Page{}
	page.Title = title
	page.Created_by_user = userid
	page.Sections = sections
	err := DAL.CreatePage(&page)
	if (expected == true && err == nil) || (expected == false && err != nil)  {
		fmt.Println("test passed")
	}else{
		fmt.Println("test failed")
		os.Exit(1)
	}	
}

func TestUpdatePage() {
	
}

func TestGetPageByID() {
	
}

func TestCreateSection() {
	
}

func TestUpdateSection() {
	
}

func TestGetSectionByID() {
	
}

func TestGetRevisions() {
	
}
func GetCurrentVersion() int {
	conn, err := WikidbDAO.GetConnection()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	command := "select version from appversion where name = 'wikidb';"
	stmt, err := conn.Prepare(command)
	defer stmt.Close()

	rs, err := stmt.Query()
	if err != nil {
		return 0
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				return 0
			}

			if !hasRow {
				break
			}else{
				//get result from rs
				version, isNull, err := rs.Int(0)
				if err != nil ||isNull{
					os.Exit(1)
				}
				return version
			}
		}
		
	}
	return 0;
}
func GetUserIDForTest() int {
	conn, err := WikidbDAO.GetConnection()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	command := "select max(user_id) from wikiuser;"
	stmt, err := conn.Prepare(command)
	defer stmt.Close()

	rs, err := stmt.Query()
	if err != nil {
		return 0
	}else{
		for {
			hasRow, err := rs.FetchNext()
			if err != nil {
				return 0
			}

			if !hasRow {
				break
			}else{
				//get result from rs
				userid, isNull, err := rs.Int(0)
				if err != nil ||isNull{
					os.Exit(1)
				}
				return userid
			}
		}
		
	}
	return 0;
}