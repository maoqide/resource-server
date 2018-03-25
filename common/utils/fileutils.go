package utils

import (
	"os"
	"path/filepath"
	"sort"

	"resource-server/entity"

	"gopkg.in/mgo.v2/bson"
)

// Walk walk file tree, append FileNode to root
func Walk(path string, info os.FileInfo, node *entity.FileNode) {
	// 列出当前目录下的所有目录、文件
	files := listFiles(path)

	// 遍历这些文件
	for _, filename := range files {
		// 拼接全路径
		fpath := filepath.Join(path, filename)

		// 构造文件结构
		fio, _ := os.Lstat(fpath)

		// 将当前文件作为子节点添加到目录下
		isDir := fio.IsDir()
		child := entity.FileNode{filename, fpath, isDir, []*entity.FileNode{}}
		node.FileNodes = append(node.FileNodes, &child)

		if isDir {
			// 如果遍历的当前文件是个目录，则进入该目录进行递归
			Walk(fpath, fio, &child)
		}
	}
	return
}

// WalkDB walk file tree, append FileNode to root, save to DB
func WalkDB(path string, info os.FileInfo, node *entity.FileNodeMgo, add2DB func(string, ...interface{}) error, collName string) {

	add2DB(collName, &node)
	// 列出当前目录下的所有目录、文件
	files := listFiles(path)
	// 遍历这些文件
	for _, filename := range files {
		// 拼接全路径
		fpath := filepath.Join(path, filename)

		// 构造文件结构
		fio, _ := os.Lstat(fpath)

		// 将当前文件作为子节点添加到目录下
		isDir := fio.IsDir()
		child := entity.FileNodeMgo{bson.NewObjectId(), filename, fpath, isDir, node.ObjectID}
		// node.Children = append(node.Children, child.ObjectID)

		if isDir {
			// 如果遍历的当前文件是个目录，则进入该目录进行递归
			WalkDB(fpath, fio, &child, add2DB, collName)
		} else {
			add2DB(collName, &child)
		}
	}
	return
}

func listFiles(dirname string) []string {
	f, _ := os.Open(dirname)

	names, _ := f.Readdirnames(-1)
	f.Close()

	sort.Strings(names)

	return names
}
