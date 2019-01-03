package model

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	Pokemons struct {
		ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Name    string        `json:"name" bson:"name"`
		Element string        `json:"element" bson:"element"`
		Weight  string        `json:"weight" bson:"weight"`
	}
)
