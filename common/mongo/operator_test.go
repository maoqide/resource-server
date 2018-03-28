package mongo

import (
	"os"
	"testing"

	. "resource-server/common/utils"
	"resource-server/entity"

	"gopkg.in/mgo.v2/bson"
)

type TestStruct struct {
	ObjectId bson.ObjectId `bson:"_id" json:"_id"`
	Name     string        `bson:"name" json:"name"`
	Path     string        `bson:"path" json:"path"`
	Testkey  string        `bson:"testkey json:"testkey"`
}

func Init() {

	sessionMng := NewSessionManagerCustom("test", "../../mongo_config.toml")
	InjectSession(sessionMng, "")
}

func TestHandleInsert(t *testing.T) {

	Init()
	document1 := &entity.Resource{ObjectId: bson.NewObjectId(), Name: "qqq.jpg", Path: "/path/to/qqq.jpg"}
	document2 := &entity.Resource{ObjectId: bson.NewObjectId(), Name: "www.jpg", Path: "/path/to/www.jpg"}
	document3 := &entity.Resource{ObjectId: bson.NewObjectId(), Name: "eee.jpg", Path: "/path/to/eee.jpg"}
	document4 := &entity.Resource{ObjectId: bson.NewObjectId(), Name: "rrr.jpg", Path: "/path/to/rrr.jpg"}

	var docs []interface{}
	docs = append(docs, document1, document2, document3, document4)
	//document := &TestStruct{ObjectId: bson.NewObjectId(), Name: "hello3.jpg", Path: "/path/to/hello3.jpg", Testkey: "ttttest"}
	err := HandleInsert("testcoll", docs...)
	if err != nil {
		t.Errorf("err: %v", err)
		t.FailNow()
	}

}
func TestHandleQueryOne(t *testing.T) {

	Init()
	document := &entity.Resource{}
	var selector = bson.M{}
	selector[ParamID()] = selector["parent"]
	err := HandleQueryOne(&document, QueryStruct{"testcoll2", selector, 0, 0, "_id"})
	t.Logf("document: %v", document)
	if err != nil {
		t.Errorf("failed, err: %v", err)
	}

}

func TestHandleQueryAll(t *testing.T) {

	Init()
	documents := []entity.Resource{}
	var selector = bson.M{}
	selector["name"] = "hello2.jpg"
	total, err := HandleQueryAll(&documents, QueryStruct{"testcoll", selector, 0, 0, "_id"})
	t.Logf("total: %d, documents: %v", total, documents)
	if err != nil {
		t.Errorf("failed, err: %v", err)
	}

}

func TestHandleUpdateOne(t *testing.T) {
	Init()
	document := &entity.Resource{ObjectId: bson.ObjectIdHex("58f421c0e1382328c2bc7856"), Name: "helloupdate.jpg", Path: "/path/to/helloupdate.jpg"}
	var selector = bson.M{}
	selector[ParamID()] = bson.ObjectIdHex("58f421c0e1382328c2bc7856")
	created, err := HandleUpdateOne(&document, QueryStruct{"testcoll", selector, 0, 0, ""})
	t.Logf("create: %v", created)
	if err != nil {
		t.Errorf("failed, err: %v", err)
	}
}

// go test -v -run ^TestInsertFileNode$
func TestInsertFileNode(t *testing.T) {
	t.Logf("begin")

	Init()
	rootpath := "/home/mao/test"
	rootID := bson.NewObjectId()
	root := entity.FileNodeMgo{rootID, "test", rootpath, true, rootID}
	fileInfo, _ := os.Lstat(rootpath)
	WalkDB(rootpath, fileInfo, &root, HandleInsert, "testcoll2")
}
