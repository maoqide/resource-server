package utils

import (
	"testing"
)

func TestListFileInfo(t *testing.T) {

	//infos := listFileInfo("/home/mao/test", nil)

	//for _, i := range infos {
	//t.Logf("name: %s, time: %d", i.Name(), i.ModTime().Unix())
	//}
}

func TestListFileFromTime(t *testing.T) {

	files := ListFileFromTime("/home/mao/test", []string{".txt"}, 1475118346)

	t.Logf("files: %v", files)
}
