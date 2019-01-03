package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	User struct {
		ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Username   string        `json:"username" bson:"username"`
		Password   string        `json:"password,omitempty" bson:"password"`
		Token      string        `json:"token,omitempty" bson:"-"`
		Name       string        `json:"name,omitempty" bson:"name"`
		Rank       string        `json:"rank,omitempty" bson:"rank"`
		Status     string        `json:"status,omitempty" bson:"status"`
		Level      string        `json:"level,omitempty" bson:"level"`
		Created_at time.Time     `json:"created_at,omitempty" bson:"created_at"`
	}
)
