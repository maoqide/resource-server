package service

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"

	"resource-server/common/mongo"
	"resource-server/common/utils"
	"resource-server/entity"
)

const (
	IMAGE_COLLECTION_NAME = "testcoll2"
)

// CoreService core service
type CoreService struct {
	collectionName string
}

var (
	coreService *CoreService
	coreOnce    sync.Once
)

// GetCoreService return singlton instance
func GetCoreService() *CoreService {
	coreOnce.Do(func() {
		coreService = &CoreService{IMAGE_COLLECTION_NAME}

		coreService.initialize()
	})

	return coreService
}

func (c *CoreService) initialize() {
	logrus.Printf("Once initialize COreService, %s", c.collectionName)
}

//ListImages list images
// mongoDB page, db.test.find(xxx).sort({"num": 1}).skip(10).limit(10)
func ListImages(start, count int) ([]entity.Resource, error) {

	documents := []entity.Resource{}
	var selector = bson.M{}

	_, err := mongo.HandleQueryAll(&documents, mongo.QueryStruct{IMAGE_COLLECTION_NAME, selector, start, count, mongo.ParamID()})
	if err != nil {
		logrus.Errorf("HandleQueryAll failed when list images, err: %v", err)
		return nil, err
	}

	return documents, nil

}

// QueryChild query child of parent
func (c *CoreService) QueryChild(documents interface{}, parent bson.ObjectId) error {
	// documents := []entity.FileNodeMgo{}
	var selector = bson.M{}
	selector["parent"] = parent
	_, err := mongo.HandleQueryAll(documents, mongo.QueryStruct{c.collectionName, selector, 0, 0, mongo.ParamID()})
	if err != nil {
		return err
	}
	return nil
}

// QueryRoot query root
func (c *CoreService) QueryRoot(documents interface{}) error {
	return c.QueryChild(documents, utils.ROOT_PARENT_ID)
}
