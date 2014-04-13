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
func Extend(slice []Section, element Section) []Section {
    n := len(slice)
    if n == cap(slice) {
        // Slice is full; must grow.
        // We double its size and add 1, so if the size is zero we still grow.
        newSlice := make([]Section, len(slice), 2*len(slice)+1)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0 : n+1]
    slice[n] = element
    return slice
}
func Append(slice []Section, items ...Section) []Section {
    for _, item := range items {
        slice = Extend(slice, item)
    }
    return slice
}