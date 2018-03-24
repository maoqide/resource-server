package utils
import (
	"testing"
)

func TestWalk(t *testing.T) {
	rootpath := "D:\\projects"

	root := FileNode{"projects", rootpath, []*FileNode{}}
	fileInfo, _ := os.Lstat(rootpath)
	walk(rootpath, fileInfo, &root)
	data, _ := json.Marshal(root)
	data, _ := json.MarshalIndent(root, "", "\t")
	fmt.Printf("%s", data)

	files := ListFileFromTime("/home/mao/test", []string{".txt"}, 1475118346)

	t.Logf("files: %v", files)
}