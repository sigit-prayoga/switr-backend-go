package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	Swit struct {
		Text   string        "json:text"
		Time   time.Time     "json:time"
		SwitId bson.ObjectId `json:"switId" bson:"switId,omitempty"`
		UserId bson.ObjectId `json:"userId" bson:"userId,omitempty"`
		Likes  []string      "json:likes"
	}
)
