package entity

import (
	"gopkg.in/mgo.v2/bson"
)

// Resource file, image or video
type Resource struct {
	ObjectId bson.ObjectId `bson:"_id" json:"_id"`
	Name     string        `bson:"name" json:"name"`
	Path     string        `bson:"path" json:"path"`
	ModTime  int64         `bson:"modtime" json:"modtime"`
}

// FileNode single nodes of a file tree
type FileNode struct {
	// ObjectID  bson.ObjectId `bson:"_id" json:"_id"`
	Name      string      `bson:"name" json:"name"`
	Path      string      `bson:"path" json:"path"`
	IsDir     bool        `bson:"isDir" json:"isDir"`
	FileNodes []*FileNode `bson:"nodes" json:"nodes"`
}

// FileNodeMgo single nodes of a file tree
type FileNodeMgo struct {
	ObjectID bson.ObjectId `bson:"_id" json:"_id"`
	Name     string        `bson:"name" json:"name"`
	Path     string        `bson:"path" json:"path"`
	IsDir    bool          `bson:"isDir" json:"isDir"`
	Parent   bson.ObjectId `bson:"parent" json:"parent"`
	// Children []bson.ObjectId `bson:"chindren" json:"children"`
	// FileNodes []*FileNode     `bson:"nodes" json:"nodes"`
}
