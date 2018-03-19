package service

import (
	"os"

	"resource-server/common/mongo"
	"resource-server/entity"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

//return a slice of Image, length of the slice
func genImages(imgpaths []string) (images []entity.Resource, length int) {

	for _, imgpath := range imgpaths {

		info, err := os.Stat(imgpath)
		if err != nil {
			logrus.Warnf("get FileInfo of [%s] failed, err: %v", imgpath, err)
			continue
		}
		image := entity.Resource{ObjectId: bson.NewObjectId(), Name: info.Name(), Path: imgpath, ModTime: info.ModTime().Unix()}
		images = append(images, image)
	}

	return images, len(images)
}

func save2db(collname string, images []entity.Resource) error {

	doc := make([]interface{}, len(images))
	for i, image := range images {
		doc[i] = image
	}
	return mongo.HandleInsert(collname, doc...)

}
