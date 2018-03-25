package utils

import (
	"encoding/json"
	"os"
	"testing"
)

func TestWalk(t *testing.T) {
	rootpath := "/home/mao/test"

	root := FileNode{"test", rootpath, true, []*FileNode{}}
	fileInfo, _ := os.Lstat(rootpath)
	walk(rootpath, fileInfo, &root)
	// data, _ := json.Marshal(root)
	data, _ := json.MarshalIndent(root, "", "\t")

	// files := ListFileFromTime("/home/mao/test", []string{".txt"}, 1475118346)

	t.Logf("tree: %s", string(data))
}
