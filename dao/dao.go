package dao

import (
	"context"
	"errors"
	"time"

	"github.com/watsonserve/residx/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MUSIC_COLLECTION = "md_music"
const AUDIO_COLLECTION = "md_audio"

type Dao interface {
	GetMusic(rid string) ([]bson.M, error)
	Find(cond map[string]interface{}, offset int64, limit int) ([]bson.M, int64, error)
	SaveAttr(rId string, key string, value string) error
}

type daoIns struct {
	db *mongo.Database
}

func find(coll *mongo.Collection, conditions bson.D, offset int64, limit int) ([]bson.M, int64, error) {
	var results []bson.M
	cntOpts := options.Count().SetMaxTime(2 * time.Second)
	total, err := coll.CountDocuments(context.TODO(), conditions, cntOpts)
	if nil != err || 0 == total || 0 == limit {
		return make([]bson.M, 0), total, err
	}

	findOpts := options.Find().SetSkip(offset).SetLimit(int64(limit))
	cursor, _err := coll.Find(context.TODO(), conditions, findOpts)
	err = _err

	if nil == err {
		err = cursor.All(context.TODO(), &results)
	}

	return results, total, err
}

func New(db *mongo.Database) Dao {
	return &daoIns{
		db: db,
	}
}

/**
 * get one resource by ID
 */
func (d *daoIns) GetMusic(rid string) ([]bson.M, error) {
	coll := d.db.Collection(AUDIO_COLLECTION)

	var results []bson.M
	cursor, err := coll.Find(context.TODO(), bson.D{{Key: "rid", Value: rid}}, options.Find())

	if nil == err {
		err = cursor.All(context.TODO(), &results)
	}

	for _, item := range results {
		item["aid"] = item["_id"]
		delete(item, "_id")
	}
	return results, err
}

/**
 * find by conditions
 */
func (d *daoIns) Find(cond map[string]interface{}, offset int64, limit int) ([]bson.M, int64, error) {
	conditions := utils.MapToKvList(cond)
	coll := d.db.Collection(MUSIC_COLLECTION)
	results, count, err := find(coll, conditions, offset, limit)

	if nil != results {
		for _, item := range results {
			item["rid"] = item["_id"]
			delete(item, "_id")
		}
	}

	return results, count, err
}

func (d *daoIns) SaveAttr(rId string, key string, value string) error {
	coll := d.db.Collection(MUSIC_COLLECTION)
	opts := options.Update().SetUpsert(true)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}}
	result, err := coll.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: rId}}, update, opts)

	if nil == err && (0 == result.MatchedCount || 0 == result.UpsertedCount) {
		err = errors.New("none record updated")
	}

	return err
}
