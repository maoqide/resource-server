package mongo

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"time"

	"resource-server/common/utils"
	"resource-server/entity"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

const (
	MONGO_URI_PREFIX = "mongodb://"
	MONGODB_CFG_TOML = "mongo_config.toml"
)

type SessionManager struct {
	defaultConfigKey string
	configMap        map[string]*entity.MgoConfig
	sessions         map[string]*mgo.Session
	accessLock       *sync.RWMutex
}

func NewSessionManager(defaultConfigKey string) *SessionManager {

	return NewSessionManagerCustom(defaultConfigKey, MONGODB_CFG_TOML)
}

func NewSessionManagerCustom(defaultConfigKey string, mongotoml string) *SessionManager {

	sessionManager := &SessionManager{
		defaultConfigKey: defaultConfigKey,
		configMap:        make(map[string]*entity.MgoConfig),
		sessions:         make(map[string]*mgo.Session),
		accessLock:       &sync.RWMutex{},
	}
	sessionManager.setConfig(mongotoml)
	return sessionManager
}

func (s *SessionManager) setConfig(mongotoml string) error {

	cm := make(map[string]*entity.MgoConfig)
	_, err := toml.DecodeFile(mongotoml, &cm)

	if err != nil {
		logrus.Errorf("parse toml config file error, err: %v", err)
		return err
	}

	s.configMap = cm
	return nil
}

func (s *SessionManager) getConfig(tomlkey string) (*entity.MgoConfig, bool) {

	config, ok := s.configMap[tomlkey]
	return config, ok
}

func (s *SessionManager) GetDefault() (*mgo.Session, bool, error) {

	return s.Get(s.defaultConfigKey)
}

// return session, exist, err
func (s *SessionManager) Get(tomlkey string) (*mgo.Session, bool, error) {

	//get config from server.config.toml
	config, exist := s.getConfig(tomlkey)
	if !exist {
		logrus.Errorf("config %v not exist when get session, ", config)
		return nil, false, errors.New("session config not exists.")
	}
	logrus.Infof("got mongodb config, config: %v", config)

	//get uri and timeout from config
	uri := generateURI(config)
	timeout := config.Timeout

	//sessionid
	sessionid := strings.TrimPrefix(MONGO_URI_PREFIX, uri)

	//check if sessionid existed
	s.accessLock.RLock()
	existSession := s.sessions[sessionid]
	s.accessLock.RUnlock()

	//copy and return if session existed
	if existSession != nil {
		return existSession.Copy(), true, nil
	}

	//connect to mongoDB
	logrus.Infof("get new session to mongodb, uri: %s, timeout: %d", uri, timeout)

	s.accessLock.Lock()
	newSession, err := mgo.DialWithTimeout(uri, time.Duration(timeout)*time.Second)
	s.accessLock.Unlock()

	if err != nil {
		logrus.Errorf("connect to db %s failed, err: %v", uri, err)
		return nil, false, err
	}

	return newSession, false, nil
}

//generate mongodb uri
//mongodb://[username:password@]host1[:port1][,host2[:port2],...][/[database][?options]]
//example: mongodb://user:passwd@127.0.0.1:27017/media
func generateURI(config *entity.MgoConfig) string {

	//start with "mongodb://"
	uriBuffer := bytes.NewBufferString(MONGO_URI_PREFIX)

	//add username and password
	if utils.ValidateStr(config.Username) {

		uriBuffer.WriteString(config.Username)
		uriBuffer.WriteString(":")
		uriBuffer.WriteString(config.Password)
		uriBuffer.WriteString("@")
	}

	//add addrs
	if len(config.Addrs) == 0 {
		uriBuffer.WriteString("localhost")
	} else {

		for i := 0; i < len(config.Addrs)-2; i++ {
			uriBuffer.WriteString(config.Addrs[i])
			uriBuffer.WriteString(",")
		}
		uriBuffer.WriteString(config.Addrs[len(config.Addrs)-1])
	}
	logrus.Debugf("Addrs: %v, length: %d", config.Addrs, len(config.Addrs))

	//add database
	if utils.ValidateStr(config.Database) {
		uriBuffer.WriteString("/")
		uriBuffer.WriteString(config.Database)
	}

	uri := uriBuffer.String()
	logrus.Infof("URI: %s", uri)
	return uri

}

// close session
func (s *SessionManager) Close(sessionid string) {

	logrus.Infof("closing mongo session, sessionid: %s", sessionid)
	s.accessLock.Lock()

	if session, exist := s.sessions[sessionid]; exist {
		session.Close()
		delete(s.sessions, sessionid)
	}
	s.accessLock.Unlock()

}

// close all sessions
func (s *SessionManager) CloseAll() {

	logrus.Infof("closing all sessions, total: %d", len(s.sessions))
	s.accessLock.Lock()
	for sessionid, session := range s.sessions {
		session.Close()
		delete(s.sessions, sessionid)
	}
	s.accessLock.Unlock()
}
