package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"resource-server/entity"
	"resource-server/service"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const (
	HTTP_METHOD_GET    string = "GET"
	HTTP_METHOD_POST   string = "POST"
	HTTP_METHOD_PUT    string = "PUT"
	HTTP_METHOD_DELETE string = "DELETE"
)

var coreService = service.GetCoreService()

type Resource struct {
}

func (r Resource) Register(router *mux.Router) {

	//router = router.Path("/mserver").Subrouter()
	// example: ip:port/images?start=0&count=10
	router.HandleFunc("/images", r.ListImgHandler).Methods(HTTP_METHOD_GET).Queries("start", "{start:[0-9]+}", "count", "{count:[0-9]+}")
	router.HandleFunc("/nodes", r.GetNodesHandler).Methods(HTTP_METHOD_GET)

}

func (r Resource) ListImgHandler(writer http.ResponseWriter, request *http.Request) {

	logrus.Infof("hello")

	//parse queries as a map into request.Form
	request.ParseForm()

	start, _ := strconv.Atoi(request.Form.Get("start"))
	count, _ := strconv.Atoi(request.Form.Get("count"))

	images, err := service.ListImages(start, count)
	if err != nil {
		logrus.Errorf("ListImages error, err: %v", err)
		writer.Write([]byte("internal err."))
		return
	}

	resp := entity.ImagesResp{Start: start, Count: count, Total: 100, Images: images}

	respbytes, err := json.Marshal(resp)
	if err != nil {
		logrus.Errorf("marshal response error, err: %v", err)
		writer.Write([]byte("internal err."))
		return
	}

	writer.Write(respbytes)
}

//GetNodesHandler get nodes of a parents
func (r Resource) GetNodesHandler(writer http.ResponseWriter, request *http.Request) {
	logrus.Infof("getnodes")

	request.ParseForm()

	parent := bson.ObjectIdHex(request.Form.Get("parent"))
	logrus.Infof("request body: %s", string(parent))

	var documents []entity.FileNodeMgo
	coreService.QueryChild(&documents, parent)
	res, err := json.Marshal(documents)
	if err != nil {
		logrus.Errorf("getnodes error, err: %v", err)
		writer.Write([]byte("internal err."))
		return
	}

	writer.Write(res)

}
