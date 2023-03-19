package actions

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/watsonserve/goengine"
	"github.com/watsonserve/residx/dao"
	"github.com/watsonserve/residx/entities"
	"github.com/watsonserve/residx/services"
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

func New(db *mongo.Database) goengine.HttpAction {
	_dao := dao.New(db)

	return &action{srv: services.New(_dao)}
}

func (a *action) Bind(router *goengine.HttpRoute) {
	router.Set("/save-attr", a.saveAttr)
	router.Set("/search", a.searchMusic)
	router.Set("/meta", a.getMusicMeta)
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

		offset := conditions["offset"].(float64)
		limit := conditions["limit"].(float64)
		delete(conditions, "offset")
		delete(conditions, "limit")
		result, err := a.srv.Find(conditions, int64(offset), int(limit))
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
		if "" == musicId {
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

type attr struct {
	Rid   string `json:"rid"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

/**
 * save attr
 */
func (a *action) saveAttr(res http.ResponseWriter, req *http.Request) {
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

		form := &attr{}
		content, err := io.ReadAll(req.Body)
		if nil == err {
			err = json.Unmarshal(content, &form)
		}
		if nil != err {
			httpCode = http.StatusBadRequest
			ret.Code = -1
			ret.Msg = "Request body format must be json"
			break
		}

		err = a.srv.SaveAttr(form.Rid, form.Key, form.Value)
		if nil != err {
			ret.Code = -1
			ret.Msg = err.Error()
		}
		break
	}

	sendJSON(res, httpCode, ret)
}
