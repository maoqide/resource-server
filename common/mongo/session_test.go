package mongo

import (
	"testing"
)

func TestGenerateURI(t *testing.T) {
	config, _ := NewSessionManagerCustom("dev", "../../mongo_config.toml").getConfig("test")

	uri := generateURI(config)

	t.Log("uri", uri)
	if len(uri) == 0 {
		t.Fail()
	}

}
