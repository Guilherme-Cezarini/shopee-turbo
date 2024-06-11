package user

import (
	"context"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	mongoDbCollection *mongo.Collection
}

type UserRepository interface {
	InsertUser(user *user.User) error
	GetUserByEmail(email string) (*user.User, error)
	InsertCodeInUser(email string, code string) error
}

func NewUserRepository(mongoClient *mongo.Database) UserRepository {
	mongoDBCollection := mongoClient.Collection(cfg.Env().MongoDbCollectionUsers)
	return &userRepository{
		mongoDbCollection: mongoDBCollection,
	}
}

// InsertCodeInUser inserts a code in the user
func (users *userRepository) InsertCodeInUser(email string, code string) error {
	filter := bson.D{{"email", email}}
	update := bson.D{{"$set", bson.D{{"code", code}}}}
	_, err := users.mongoDbCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil

}

// InsertUser inserts a new user in the database
func (users *userRepository) InsertUser(user *user.User) error {
	_, err := users.mongoDbCollection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByEmail gets a user by email
func (users *userRepository) GetUserByEmail(email string) (*user.User, error) {
	filter := bson.D{{"email", email}}
	var user user.User
	err := users.mongoDbCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
