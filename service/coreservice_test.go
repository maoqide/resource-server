package service

import (
	"mediaserver/common/mongo"
	"testing"
)

func TestCoreservice(t *testing.T) {

	sessionMng := mongo.NewSessionManagerCustom("test", "../mongo_config.toml")
	mongo.InjectSession(sessionMng, "")

	imgs, err := ListImages()
	if err != nil {
		t.Errorf("err: %v", err)
		t.Fail()
	}
	t.Log(imgs)
}
