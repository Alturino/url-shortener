package request

import (
	"encoding/json"
)

type UrlRequest struct {
	Url string
}

func (u *UrlRequest) String() string {
	json, _ := json.Marshal(u)
	return string(json)
}
