package model

import (
	//"time"
)

type Section struct {
	PageID int
	LatestRevisionID int
	SectionID int
	Title string
	Text string
	CreatedByUserID int
	LastUpdateByUserID int
}