package services

import (
	"go.mongodb.org/mongo-driver/bson"
	"github.com/watsonserve/scaner/dao"
)

type Srv struct {
	daoIns     dao.Dao
	root       string
	searchStat int
}

func New(daoIns dao.Dao, root string) *Srv {
	// 1 doing 2 done 4 error 8 failed
	// 6 errList 10 failed_msg
	searchStat := 0

	return &Srv{daoIns, root, searchStat}
}

func (s *Srv) GetMusicMeta(id string) ([]bson.M, error) {
	return s.daoIns.GetMusic(id)
}

func (s *Srv) Find(cond map[string]interface{}, offset int64, limit int) (map[string]interface{}, error) {
	result, total, err := s.daoIns.Find(cond, offset, limit)
	if nil != err {
		return nil, err
	}
	ret := make(map[string]interface{})
	ret["list"] = result
	ret["total"] = total
	return ret, err
}

func (s *Srv) SaveAttr(rId string, key string, value string) error {
	return s.daoIns.SaveAttr(rId, key, value)
}