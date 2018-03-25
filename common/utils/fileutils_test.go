package utils

import (
	"encoding/json"
	"os"
	"testing"

	"resource-server/entity"
)

func TestWalk(t *testing.T) {
	rootpath := "/home/mao/test"
	root := entity.FileNode{"test", rootpath, true, []*entity.FileNode{}}
	fileInfo, _ := os.Lstat(rootpath)
	Walk(rootpath, fileInfo, &root)
	// data, _ := json.Marshal(root)
	data, _ := json.MarshalIndent(root, "", "\t")
	t.Logf("tree: %s", string(data))
}
