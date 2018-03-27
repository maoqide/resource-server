package initializer

import (
	"fmt"
	"os"
	"strconv"
	"time"

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

func copy2coll(src, dest string) (int, error) {
	var documents []interface{}
	total, err := mongo.HandleQueryAll(&documents, mongo.QueryStruct{src, nil, 0, 0, "_id"})
	fmt.Printf("total: %d", total)
	fmt.Printf("documents: %v", documents)
	if err != nil || total == 0 {
		return 0, err
	}
	err = mongo.HandleInsert(dest, documents...)
	if err != nil {
		return total, err
	}
	return total, nil
}

func dumpFile2DB() {

	// bakup & clear collection
	var collName = "testcoll2"

	destColl := collName + "_bak_" + strconv.FormatInt(time.Now().Unix(), 10)
	fmt.Printf("bak collection to %s \n", destColl)
	_, err := copy2coll(collName, destColl)
	if err != nil {
		panic(err)
	}

	fmt.Println("empty collection")
	err = mongo.HandleDelete(collName, nil, false)
	if err != nil {
		panic(err)
	}

	//walk file path and save to mongoDB
	fmt.Println("walk file tree")
	rootpath := "/home/mao/test"
	rootID := bson.NewObjectId()
	root := entity.FileNodeMgo{rootID, "test", rootpath, true, rootID}
	fileInfo, _ := os.Lstat(rootpath)
	utils.WalkDB(rootpath, fileInfo, &root, mongo.HandleInsert, collName)

	//create index on 'parent'
	mongo.GetCollection(collName, []string{"parent"})
}
