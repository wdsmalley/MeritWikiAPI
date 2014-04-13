package model

import (
	"time"
)

type PageRevision struct {
	page Page
	revision_id int
	revision_create_user User
	create_time time
}