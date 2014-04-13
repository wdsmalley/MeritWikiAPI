package model

import (
	//"errors"
	"encoding/json"
)

type Response struct {
	ResponseObject interface{}
	Error string
}

func (r * Response) ToJsonString() (string){
	jsonData, err := json.Marshal(r)
	
	if err != nil {
		return ""
	}

	return string(jsonData)
}