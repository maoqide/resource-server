package entity

import (
	"gopkg.in/mgo.v2/bson"
)

// Resouce file, image or video
type Resource struct {
	ObjectId bson.ObjectId `bson:"_id" json:"_id"`
	Name     string        `bson:"name" json:"name"`
	Path     string        `bson:"path" json:"path"`
	ModTime  int64         `bson:"modtime" json:"modtime"`
}
