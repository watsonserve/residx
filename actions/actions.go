package actions

import (
	"encoding/json"
	"io"
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
	router.Set("/autocomplete", a.searchMusic)
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
		if "POST" != req.Method {
			httpCode = http.StatusMethodNotAllowed
			ret.Code = -1
			ret.Msg = "request method must be POST"
			break
		}

		conditions := map[string]interface{}{}
		content, err := io.ReadAll(req.Body)
		if nil == err {
			err = json.Unmarshal(content, &conditions)
		}
		if nil != err {
			httpCode = http.StatusBadRequest
			ret.Code = -1
			ret.Msg = "Request body format must be json"
			break
		}

		offset := conditions["offset"]
		limit := conditions["limit"]
		delete(conditions, "offset")
		delete(conditions, "limit")
		result, err := a.srv.Find(conditions, offset.(int64), limit.(int))
		ret.Data = result

		if nil != err {
			ret.Code = -1
			ret.Msg = err.Error()
		}

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

/**
 * search
 */
func (d *daoIns) autoComplete(res http.ResponseWriter, req *http.Request) {

}
