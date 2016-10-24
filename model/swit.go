package model

import "time"

type Swit struct {
	Text   string    "json:text"
	Time   time.Time "json:time"
	SwitId string    `json:"switId" bson:"switId,omitempty"`
	UserId string    `json:"userId" bson:"userId,omitempty"`
	Likes  []string  "json:likes"
}
