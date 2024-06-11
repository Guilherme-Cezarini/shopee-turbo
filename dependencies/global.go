package dependencies

import (
	"context"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	hash2 "github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/service/hash"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Global struct {
	Mongo   mongoDep
	HashKit hash2.HashingInterface
}

type mongoDep struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var instance *Global

func LoadGlobalDependencies(ctx context.Context) (*Global, error) {
	if instance == nil {
		mongoClient, err := loadMongoDB(ctx)
		if err != nil {
			return nil, err
		}
		mongoDep := mongoDep{
			Client:   mongoClient.Client,
			Database: mongoClient.Database,
		}

		hashing, err := hash2.NewHashing(hash2.HashingConfig{
			First:  cfg.Env().HashingFirst,
			Second: cfg.Env().HashingSecond,
			Third:  cfg.Env().HashingThird,
			Fourth: cfg.Env().HashingFourth,
		})

		instance = &Global{
			Mongo:   mongoDep,
			HashKit: hashing,
		}
	}

	return instance, nil
}

func loadMongoDB(ctx context.Context) (*mongoDep, error) {
	clientOptions := options.Client().
		SetMinPoolSize(cfg.Env().MongoDbPoolMinSize).
		SetMaxPoolSize(cfg.Env().MongoDbPoolMaxSize).
		SetMaxConnIdleTime(cfg.Env().MongoDbMaxConnIdleTime).
		SetConnectTimeout(cfg.Env().MongoDbConnectTimeout).
		ApplyURI(cfg.Env().MongoDbURL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to mongodb")
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to mongodb")
	}

	return &mongoDep{
		Client:   client,
		Database: client.Database(cfg.Env().MongoDbDatabase),
	}, nil
}
