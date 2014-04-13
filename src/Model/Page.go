package model

type Page struct {
	PageID int
	Title string
	Url string
	RevisionID int
	Created_by_user int
	Created_by_username string
	Email string
	//create_time time
	Sections []Section
	//revisions []Revision
}
