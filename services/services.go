package services

import (
	"github.com/watsonserve/scaner/dao"
)

type Srv struct {
	daoIns     dao.Dao
	root       string
	searchStat int
	write      chan int64
}

func New(daoIns dao.Dao, root string) *Srv {
	// 1 doing 2 done 4 error 8 failed
	// 6 errList 10 failed_msg
	searchStat := 0
	write := make(chan int64)

	srv := &Srv{daoIns, root, searchStat, write}
	go srv.listen()
	return srv
}

func (s *Srv) listen() {
	for {
		<-s.write
		if 0 == s.searchStat {
			s.searchStat = 1
			go s.makeResourcesIndex()
		}
	}
}

func (s *Srv) makeResourcesIndex() {
	audioList, errList, err := search(s.root)

	if nil != err {
		s.searchStat = 10
		return
	}
	if nil != errList {
		s.searchStat = 6
		return
	}

	err = s.daoIns.SaveResources(audioList)
	if nil != err {
		s.searchStat = 0
	}
	s.searchStat = 2
}

func (s *Srv) MakeResourcesIndex() int {
	s.write <- 0
	return s.searchStat
}

func (s *Srv) GetMusicMeta(id string) (map[string]interface{}, error) {
	return s.daoIns.GetMusicMeta(id)
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
