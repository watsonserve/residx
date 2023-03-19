package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/watsonserve/goutils"
	"github.com/watsonserve/residx/entities"
	"github.com/watsonserve/residx/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MUSIC_COLLECTION = "md_music"
const AUDIO_COLLECTION = "md_audio"

type Saver struct {
	musicCOll *mongo.Collection
	audioCOll *mongo.Collection
	opts      *options.UpdateOptions
}

func ConnDB(dbAddr string, dbName string) (*Saver, error) {
	clientOpts := options.Client().ApplyURI("mongodb://" + dbAddr)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if nil != err {
		return nil, err
	}
	db := client.Database(dbName)
	if nil == db {
		client.Disconnect(context.TODO())
		return nil, errors.New("Connect to mongodb failed")
	}

	return &Saver{
		musicCOll: db.Collection(MUSIC_COLLECTION),
		audioCOll: db.Collection(AUDIO_COLLECTION),
		opts:      options.Update().SetUpsert(true),
	}, nil
}

func (s *Saver) Close() {
	s.audioCOll.Database().Client().Disconnect(context.TODO())
}

func (s *Saver) updateOne(coll *mongo.Collection, data bson.M) (interface{}, error) {
	filter := utils.MapToKvList(data)
	that, err := coll.UpdateOne(context.TODO(), filter, data, s.opts)
	return that.UpsertedID, err
}

/**
* save
 */
func (s *Saver) saveResource(song *entities.Song) error {
	m, err := utils.ToMap(song)
	if nil != err {
		return err
	}

	musicInfo := utils.MapFilter(m, []string{"title", "album", "artist"})
	id, merr := s.updateOne(s.musicCOll, musicInfo)
	if nil != merr {
		return merr
	}
	audioInfo := utils.MapFilter(m, []string{"url", "hash", "sample_rate", "bit_rate", "channels", "duration"})
	audioInfo["rid"] = id
	_, aerr := s.updateOne(s.audioCOll, audioInfo)
	return aerr
}

func main() {
	options := []goutils.Option{
		{
			Opt:       'h',
			Option:    "help",
			HasParams: false,
		},
	}
	helpTxt := goutils.GenHelp(options, "dbAddr dbName directory")
	argMap, params := goutils.GetOptions(options)
	if 3 != len(params) {
		argMap["help"] = ""
	}
	_, hasHelp := argMap["help"]
	if hasHelp {
		fmt.Println(helpTxt)
		return
	}
	dbAddr := params[0]
	dbName := params[1]
	dir := params[2]

	saver, err := ConnDB(dbAddr, dbName)
	if nil != err {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	list, errs, err := search(dir)
	if nil != err {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	for _, fErr := range errs {
		fmt.Printf("%s: %s\n", fErr.Filename, fErr.Error())
	}
	if 0 == len(list) {
		return
	}

	for _, item := range list {
		saver.saveResource(item)
	}

	saver.Close()
}
