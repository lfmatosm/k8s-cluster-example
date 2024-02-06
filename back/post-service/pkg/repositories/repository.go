package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Save(collection string, data interface{}) (bool, error)
	List(collection string, result interface{}) error
}

type MongoRepository struct {
	client       *mongo.Client
	uri          string
	databaseName string
	db           *mongo.Database
}

func NewMongoRepository(uri, databaseName string) *MongoRepository {
	return &MongoRepository{
		uri:          uri,
		databaseName: databaseName,
	}
}

func (m *MongoRepository) connect() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.uri))
	if err != nil {
		return err
	}
	m.client = client
	m.db = client.Database(m.databaseName)
	return nil
}

func (m *MongoRepository) Save(collection string, data interface{}) (bool, error) {
	if m.client == nil {
		if err := m.connect(); err != nil {
			return false, err
		}
	}

	coll := m.db.Collection(collection)
	_, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		return false, err
	}
	return true, err
}

func (m *MongoRepository) List(collection string, result interface{}) error {
	if m.client == nil {
		if err := m.connect(); err != nil {
			return err
		}
	}

	coll := m.db.Collection(collection)
	cur, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	if err = cur.All(context.TODO(), result); err != nil {
		return err
	}

	return nil
}
