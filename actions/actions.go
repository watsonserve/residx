package actions

import (
	"encoding/json"
	"net/http"

	"github.com/watsonserve/goengine"
	"github.com/watsonserve/scaner/dao"
	"github.com/watsonserve/scaner/entities"
	"github.com/watsonserve/scaner/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func sendJSON(res http.ResponseWriter, httpCode int, body *entities.StdJSONPacket) {
	jsonData, err := json.MarshalIndent(body, "", "")
	if nil != err {
		httpCode = http.StatusInternalServerError
		jsonData = []byte("{\"code\": -1, \"msg\": \"" + err.Error() + "\", \"data\": null}")
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(httpCode)
	res.Write(jsonData)
}

type action struct {
	srv *services.Srv
}

func New(db *mongo.Database, root string) goengine.HttpAction {
	_dao := dao.New(db)

	return &action{srv: services.New(_dao, root)}
}

func (a *action) Bind(router *goengine.HttpRoute) {
	router.Set("/scan-music", a.scanMusic)
	router.Set("/search-music", a.searchMusic)
	router.Set("/music-meta", a.getMusicMeta)
}

func (a *action) searchMusic(res http.ResponseWriter, req *http.Request) {
	httpCode := 200
	ret := &entities.StdJSONPacket{
		Code: 0,
		Msg:  "",
		Data: nil,
	}

	for {
		if "GET" != req.Method {
			httpCode = http.StatusMethodNotAllowed
			ret.Code = -1
			ret.Msg = "request method must be GET"
			break
		}

		query := req.URL.Query()
		attr := query.Get("attr")
		value := query.Get("value")
		if "" != musicId {
			httpCode = http.StatusBadRequest
			ret.Code = -1
			ret.Msg = "id is required"
			break
		}

		meta, err := a.srv.GetMusicMeta(musicId)

		if nil != err {
			ret.Code = -1
			ret.Msg = err.Error()
			break
		}

		ret.Data = meta
		break
	}

	sendJSON(res, httpCode, ret)
}

func (a *action) getMusicMeta(res http.ResponseWriter, req *http.Request) {
	httpCode := 200
	ret := &entities.StdJSONPacket{
		Code: 0,
		Msg:  "",
		Data: nil,
	}

	for {
		if "GET" != req.Method {
			httpCode = http.StatusMethodNotAllowed
			ret.Code = -1
			ret.Msg = "request method must be GET"
			break
		}

		query := req.URL.Query()
		musicId := query.Get("id")
		if "" != musicId {
			httpCode = http.StatusBadRequest
			ret.Code = -1
			ret.Msg = "id is required"
			break
		}

		meta, err := a.srv.GetMusicMeta(musicId)

		if nil != err {
			ret.Code = -1
			ret.Msg = err.Error()
			break
		}

		ret.Data = meta
		break
	}

	sendJSON(res, httpCode, ret)
}

func (a *action) scanMusic(res http.ResponseWriter, req *http.Request) {
	a.srv.MakeResourcesIndex()
	headers := res.Header()
	headers.Set("Content-Type", "application/json")
	res.Write([]byte("{}"))
}
