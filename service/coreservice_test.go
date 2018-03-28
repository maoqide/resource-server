package service

import (
	"resource-server/common/mongo"
	"resource-server/entity"
	"testing"
)

// func TestCoreservice(t *testing.T) {

// 	sessionMng := mongo.NewSessionManagerCustom("test", "../mongo_config.toml")
// 	mongo.InjectSession(sessionMng, "")

// 	imgs, err := ListImages()
// 	if err != nil {
// 		t.Errorf("err: %v", err)
// 		t.Fail()
// 	}
// 	t.Log(imgs)
// }

func TestQueryChild(t *testing.T) {
	sessionMng := mongo.NewSessionManagerCustom("test", "../mongo_config.toml")
	mongo.InjectSession(sessionMng, "")
	t.Log("asdasd")
	documents := []entity.FileNodeMgo{}
	cs := GetCoreService()
	t.Log(cs.collectionName)
	err := cs.queryRoot(&documents)
	if err != nil {
		t.Errorf("err: %v", err)
		t.Fail()
	}
	t.Logf("root %v", documents)
	documents2 := []entity.FileNodeMgo{}
	err = cs.queryChild(&documents2, documents[0].ObjectID)
	if err != nil {
		t.Errorf("err: %v", err)
		t.Fail()
	}
	t.Logf("child %v len %d", documents2, len(documents2))
	var entities = make([]entity.Entity, len(documents2))
	t.Logf(".................. %d", len(entities))

	for i, d := range documents2 {
		// entities = append(entities, &d)
		entities[i] = &d
	}

	t.Logf("entities %v len %d", entities, len(entities))
	jsonA, _ := Entities2JSON(entities)
	t.Logf("json %v", jsonA)
}
