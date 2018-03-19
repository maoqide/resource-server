package service

import (
	"testing"

	"mediaserver/common/mongo"
	"mediaserver/common/utils"
)

func TestSave2db(t *testing.T) {

	sessionMng := mongo.NewSessionManagerCustom("test", "../mongo_config.toml")
	mongo.InjectSession(sessionMng, "")

	files := utils.ListFileFromTime("/home/mao/test", nil, 1475118346)
	images, _ := genImages(files)
	err := save2db("testcoll", images)
	if err != nil {
		t.Fail()
	}

}
