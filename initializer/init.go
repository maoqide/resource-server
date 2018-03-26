package initializer

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2/bson"

	"resource-server/common/mongo"
	"resource-server/common/utils"
	"resource-server/entity"
)

// Init save file tree to mongoDB
func Init() {
	fmt.Println("init")
	createMgoSession()
	dumpFile2DB()
}

func createMgoSession() {
	sessionMng := mongo.NewSessionManagerCustom("test", "mongo_config.toml")
	mongo.InjectSession(sessionMng, "")
}

func dumpFile2DB() {

	// bakup & clear collection
	// TODO

	//walk file path and save to mongoDB
	rootpath := "/home/mao/test"
	rootID := bson.NewObjectId()
	root := entity.FileNodeMgo{rootID, "test", rootpath, true, rootID}
	fileInfo, _ := os.Lstat(rootpath)
	utils.WalkDB(rootpath, fileInfo, &root, mongo.HandleInsert, "testcoll2")

	//create index on 'parent'
	mongo.GetCollection("testcoll2", []string{"parent"})

}
