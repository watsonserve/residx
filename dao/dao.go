package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dao interface {
	GetMusicMeta(id string) (map[string]interface{}, error)
}

type daoIns struct {
	db *mongo.Database
}

func New(db *mongo.Database) Dao {
	return &daoIns{
		db: db,
	}
}

func (d *daoIns) GetMusicMeta(id string) (map[string]interface{}, error) {
	coll := d.db.Collection("music")

	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&result)
	if mongo.ErrNoDocuments == err {
		err = nil
	}

	return result, err
}
