package user

import (
	"context"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	mongoDbCollection *mongo.Collection
}

type UserRepository interface {
	InsertUser(user *user.User) error
}

func NewUserRepository(mongoClient *mongo.Database) UserRepository {
	mongoDBCollection := mongoClient.Collection(cfg.Env().MongoDbCollectionUsers)
	return &userRepository{
		mongoDbCollection: mongoDBCollection,
	}
}

// InsertUser inserts a new user in the database
func (users *userRepository) InsertUser(user *user.User) error {
	_, err := users.mongoDbCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
