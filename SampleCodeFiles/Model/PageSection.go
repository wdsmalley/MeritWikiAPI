package model
import (
	"time"
)

type PageSection struct {
	Page_id int
	Revision int
	Section_id int
	Section_title string
	Section_text string
	Created_by_user int
	create_time time
}
func NewPageSection( ipageid int, ititle string, itext string, icreatedbyuser int) *PageSection {
    p := new(PageSection)
	p.Page_id = ipageid
    p.Section_title = ititle
    p.Section_text = itext
	p.Created_by_user = icreatedbyuser
    return p
}