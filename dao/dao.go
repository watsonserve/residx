package dao

import (
	"context"
	"time"

	"github.com/watsonserve/scaner/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dao interface {
	GetMusicMeta(id string) (map[string]interface{}, error)
	Find(cond map[string]interface{}, offset int64, limit int) ([]bson.M, int64, error)
	SaveResources(metas []*entities.AudioMeta) error
}

type daoIns struct {
	db *mongo.Database
}

func find(coll *mongo.Collection, conditions bson.D, offset int64, limit int) ([]bson.M, int64, error) {
	var results []bson.M
	cntOpts := options.Count().SetMaxTime(2 * time.Second)
	total, err := coll.CountDocuments(context.TODO(), conditions, cntOpts)
	if nil != err || 0 == total || 0 == limit {
		return results, total, err
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
func (d *daoIns) GetMusicMeta(id string) (map[string]interface{}, error) {
	coll := d.db.Collection("music")

	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&result)
	if mongo.ErrNoDocuments == err {
		err = nil
	}

	return result, err
}

/**
 * find by conditions
 */
func (d *daoIns) Find(cond map[string]interface{}, offset int64, limit int) ([]bson.M, int64, error) {
	conditions := make([]bson.E, 0)
	for key, value := range cond {
		conditions = append(conditions, bson.E{Key: key, Value: value})
	}

	coll := d.db.Collection("music")
	return find(coll, conditions, offset, limit)
}

/**
 * save
 */
func (d *daoIns) SaveResources(metas []*entities.AudioMeta) error {
	coll := d.db.Collection("music")

	docs := make([]interface{}, len(metas))
	for i, item := range metas {
		docs[i] = item
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := coll.InsertMany(context.TODO(), docs, opts)
	return err
}
