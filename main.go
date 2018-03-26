package main

import (
	"net/http"

	"resource-server/api"
	// "resource-server/common/mongo"

	"resource-server/initializer"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func main() {

	initializer.Init()
	logrus.Infof("starting server...")

	// sessionMng := mongo.NewSessionManagerCustom("test", "mongo_config.toml")
	// mongo.InjectSession(sessionMng, "")

	r := api.Resource{}
	router := mux.NewRouter()
	r.Register(router)

	logrus.Fatal(http.ListenAndServe(":18080", router))
}
