package mongo

import (
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var operator *Operator

const paramId = "_id" //mongo id parameter

type Operator struct {
	sessionManager *SessionManager
	customDB       string //custom database name instead of default database in mongodb session uri,
}

type QueryStruct struct {
	CollectionName string
	Selector       bson.M
	Skip           int
	Limit          int
	Sort           string
}

func ParamID() string {
	return paramId
}

//handle insert to mongoDB, one or more
func HandleInsert(collName string, document ...interface{}) error {

	//validate global varible 'operator'
	if !operatorValidate() {
		logrus.Errorf("operator not injected, must InjectOperator() first.")
		return errors.New("MONGO.OPERATOR NOT INJECTED")
	}

	//get mongo session
	session, existed, err := operator.sessionManager.GetDefault()
	if err != nil {
		logrus.Errorf("get mongo session failed, err: %v", err)
		return err
	}

	//if already exist one, close when operation completed
	if existed {
		defer session.Close()
	}

	//get mongo collection
	coll := getCollection(collName, session)

	//insert document to collection
	if err = coll.Insert(document...); err != nil {
		logrus.Errorf("insert document to collection error, err: %v", err)
		return err
	}

	logrus.Debugf("HandleInsert document [v%] into collection [s%] success.", document, collName)
	return nil

}

func HandleQueryOne(document interface{}, queryStruct QueryStruct) error {

	//validate global varible 'operator'
	if !operatorValidate() {
		logrus.Errorf("operator not injected, must InjectOperator() first.")
		return errors.New("MONGO.OPERATOR NOT INJECTED")
	}

	//get mongo session
	session, existed, err := operator.sessionManager.GetDefault()
	if err != nil {
		logrus.Errorf("get mongo session failed, err: %v", err)
		return err
	}

	//if already exist one, close when operation completed
	if existed {
		defer session.Close()
	}

	// compose mongo Query
	query := composeQuery(session, queryStruct, true)

	//get one document
	logrus.Debugf("HandleQueryOne success, document: %v, queryStruct: %v", document, queryStruct)
	return query.One(document)
}

// return total number of documents returned in result set, err
func HandleQueryAll(documents interface{}, queryStruct QueryStruct) (int, error) {

	//validate global varible 'operator'
	if !operatorValidate() {
		logrus.Errorf("operator not injected, must InjectOperator() first.")
		return 0, errors.New("MONGO.OPERATOR NOT INJECTED")
	}

	//get mongo session
	session, existed, err := operator.sessionManager.GetDefault()
	if err != nil {
		logrus.Errorf("get mongo session failed, err: %v", err)
		return 0, err
	}

	//if already exist one, close when operation completed
	if existed {
		defer session.Close()
	}

	//compose mongo Query
	query := composeQuery(session, queryStruct, false)

	// get all documents
	if err = query.All(documents); err != nil {

		return 0, err
	}

	logrus.Debugf("HandleQueryAll success, document: %v, queryStruct: %v", documents, queryStruct)
	return query.Count()
}

// return created, err
func HandleUpdateOne(document interface{}, queryStruct QueryStruct) (bool, error) {

	//validate global varible 'operator'
	if !operatorValidate() {
		logrus.Errorf("operator not injected, must InjectOperator() first.")
		return false, errors.New("MONGO.OPERATOR NOT INJECTED")
	}

	//get mongo session
	session, existed, err := operator.sessionManager.GetDefault()
	if err != nil {
		logrus.Errorf("get mongo session failed, err: %v", err)
		return false, err
	}

	//if already exist one, close when operation completed
	if existed {
		defer session.Close()
	}

	//get mongo collection
	coll := getCollection(queryStruct.CollectionName, session)

	//update document by id
	info, err := coll.UpsertId(queryStruct.Selector[paramId], document)
	if err != nil {
		return false, err
	}

	created := false
	if info != nil && info.UpsertedId != nil {

		created = (info.Updated == 0)
	}

	logrus.Debugf("HandleUpdateOne success, document: %v, queryStruct: %v", document, queryStruct)
	return created, nil
}

// HandleDelete remove records of collection, Removeall if one is false,
// empty collection if selector is nil and one is false
func HandleDelete(collName string, selector bson.M, one bool) error {

	// validate global varible 'operator'
	if !operatorValidate() {
		logrus.Errorf("operator not injected, must InjectOperator() first.")
		return errors.New("MONGO.OPERATOR NOT INJECTED")
	}

	// get mongo session
	session, existed, err := operator.sessionManager.GetDefault()
	if err != nil {
		logrus.Errorf("get mongo session failed, err: %v", err)
		return err
	}

	// if already exist one, close when operation completed
	if existed {
		defer session.Close()
	}

	// if len(selector) == 0 {
	// 	logrus.Errorf("error handledelete mongo collection [s%], selector is empty", collName)
	// 	return errors.New("MONGO DELETE COLLECTION FAILED, SELECTOR CANNOT BE EMPTY")
	// }

	// get mongo collection
	coll := getCollection(collName, session)

	logrus.Debugf("HandleDelete success, collName: %s, selector: %v, one: %b", collName, selector, one)
	if one {
		return coll.Remove(selector)
	}
	_, err = coll.RemoveAll(selector)
	return err

}

func HandleUpdateAllBySelector(collName string, selector bson.M, document interface{}) error {

	// validate global varible 'operator'
	if !operatorValidate() {
		logrus.Errorf("operator not injected, must InjectOperator() first.")
		return errors.New("MONGO.OPERATOR NOT INJECTED")
	}

	// get mongo session
	session, existed, err := operator.sessionManager.GetDefault()
	if err != nil {
		logrus.Errorf("get mongo session failed, err: %v", err)
		return err
	}

	// if already exist one, close when operation completed
	if existed {
		defer session.Close()
	}

	// get mongo collection
	coll := getCollection(collName, session)

	// update
	change := bson.M{"$set": document}
	_, err = coll.UpdateAll(selector, change)
	logrus.Debugf("HandleUpdateAllBySelector success, collName: %s, selector: %v, document: %v", document, selector, document)
	return err
}

func HandleUpdateBySelector(collName string, selector bson.M, document interface{}) error {

	// validate global varible 'operator'
	if !operatorValidate() {
		logrus.Errorf("operator not injected, must InjectOperator() first.")
		return errors.New("MONGO.OPERATOR NOT INJECTED")
	}

	// get mongo session
	session, existed, err := operator.sessionManager.GetDefault()
	if err != nil {
		logrus.Errorf("get mongo session failed, err: %v", err)
		return err
	}

	// if already exist one, close when operation completed
	if existed {
		defer session.Close()
	}

	// get mongo collection
	coll := getCollection(collName, session)

	// update
	change := bson.M{"$set": document}

	logrus.Debugf("HandleUpdateBySelector success, collName: %s, selector: %v, document: %v", document, selector, document)
	return coll.Update(selector, change)
}

// conpose mongo Query from QueryStruct
func composeQuery(session *mgo.Session, queryStruct QueryStruct, one bool) *mgo.Query {

	logrus.Debugf("composeQuery, queryStruct: %v", queryStruct)

	//get mongo collection
	coll := getCollection(queryStruct.CollectionName, session)

	//create a mongo Query
	query := coll.Find(queryStruct.Selector)

	// if select one from _id parameter
	if one {
		return query
	}

	// Number of documents to skip in result set
	query.Skip(queryStruct.Skip)

	// Maximum number of documents in the result set
	query.Limit(queryStruct.Limit)

	// sort from comma separate list in querystruct
	if len(queryStruct.Sort) == 0 {
		query.Sort(paramId)
	}
	query.Sort(strings.Split(queryStruct.Sort, ",")...)

	return query
}

//get collection, if operator.customDB == "", mongo will use default database in session uri
func getCollection(collName string, session *mgo.Session) *mgo.Collection {

	return session.DB(operator.customDB).C(collName)
}

//get collection, and create an index
func getCollectionWithIndex(collName string, session *mgo.Session, key []string) *mgo.Collection {

	collection := session.DB(operator.customDB).C(collName)
	collection.EnsureIndex(mgo.Index{Key: key})
	return collection
}

// GetCollection get mongoDB collection, optional EnsureIndex
func GetCollection(collName string, index []string) (*mgo.Collection, error) {

	//validate global varible 'operator'
	if !operatorValidate() {
		logrus.Errorf("operator not injected, must InjectOperator() first.")
		return nil, errors.New("MONGO.OPERATOR NOT INJECTED")
	}

	//get mongo session
	session, existed, err := operator.sessionManager.GetDefault()
	if err != nil {
		logrus.Errorf("get mongo session failed, err: %v", err)
		return nil, err
	}

	//if already exist one, close when operation completed
	if existed {
		defer session.Close()
	}

	if len(index) > 0 {
		return getCollectionWithIndex(collName, session, index), nil
	}
	return getCollection(collName, session), nil
}

// validate if operator has been injected.
func operatorValidate() bool {

	if operator == nil {
		return false
	}

	return true
}

func InjectSession(sessionMng *SessionManager, customDB string) {

	operator = &Operator{sessionManager: sessionMng, customDB: customDB}
}
