package DAL

import (
	"errors"
	"../DAO"
	"../Model"
)

func Authenticate(username string, password string) (* model.User, error) {
	if username == "demo" {
		return &model.User{}, nil
	} else {
		return nil, errors.New("DAL Test Error")
	}
}

func CreateUser(user * model.User) (error){
	err := WikidbDAO.CreateUser(user)
	if user.UserID >= 1 {
		return  nil
	} else {
		return  err
	}
}
func UpdateUser(user * model.User) (error){
	if user.UserID >= 1 {
		err := WikidbDAO.UpdateUser(user)
		if err != nil{
			return  errors.New("DAO error updating user")
		}
	} else {
		return  errors.New("Error updating User")
	}
	return nil
}
func CreatePage(page * model.Page) (error){
	WikidbDAO.CreatePage(page)

	if page.PageID >= 1 {
		for i := range page.Sections {
			page.Sections[i].PageID = page.PageID 
			page.Sections[i].CreatedByUserID = page.Created_by_user
			page.Sections[i].SectionID = i+1
			WikidbDAO.CreatePageSection(&page.Sections[i])
		}
		return nil
	}else{
		return errors.New("Error Creating Page in DAO")
	}
}
func UpdatePage(page * model.Page) (error){
	if page.PageID >= 1 {
		err := WikidbDAO.UpdatePage(page)
		if err != nil{
			return err
		}else{
			return nil
		}
	}else{
		return errors.New("Invalid page id for update")
	}
}
func UpdateSection(sections []model.Section) (error){
	for i := range sections {
		if sections[i].PageID >= 1 {
			err := WikidbDAO.UpdatePageSection(&sections[i])
			if err != nil{
				return err
			}
		}else{
			return errors.New("Invalid page id for update")
		}
	}
	return nil
}
func GetUserByID(user_id int) (* model.User, error){
	if user_id == 0 {
		return &model.User{}, nil
	} else {
		return nil, errors.New("DAL Test Error")
	}
}

func GetPageByID(page_id int) (* model.Page, error){
	if page_id == 0 {
		return &model.Page{}, nil
	} else {
		return nil, errors.New("DAL Test Error")
	}
}

func GetSectionRevisionsByID(section_id int) ([]model.Revision, error){
	if section_id == 0 {
		return []model.Revision{model.Revision{}},nil
	} else {
		return nil, errors.New("DAL Test Error")
	}
}

func GetSectionByID(section_id int) (* model.Section, error){
	if section_id == 0 {
		return &model.Section{}, nil
	} else {
		return nil, errors.New("DAL Test Error")
	}
}