package dependencies

import (
	"context"
	"crypto/tls"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

package dependencies

import (
"context"
"crypto/tls"

"github.com/pkg/errors"
"github.com/rs/zerolog/log"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

type Global struct {
	Mongo   mongoDep

}

type mongoDep struct {
	Client      *mongo.Client
	Database    *mongo.Database

}

var instance *Global

func LoadGlobalDependencies(ctx context.Context) (*Global, error) {
	if instance == nil {
		mongoClient, err := loadMongoDB(ctx)
		if err != nil {
			return nil, err
		}
		mongoDep := mongoDep{
			Client:      mongoClient.Client,
			Database:    mongoClient.Database,
			Transaction: nil,
		}
		instance = &Global{
			Mongo:   mongoDep,
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

	client, err := mongodb.NewClient(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &mongoDep{
		Client:      client,
		Database:    client.Database(cfg.Env().MongoDbDatabase),

	}, nil
}

/*func BuildMongoTransaction(mongoClient *mongo.Client) *mongodb.MongoTransaction {
	txnOpts := options.Transaction().
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetReadConcern(readconcern.Local()).
		SetReadPreference(readpref.Primary())

	return mongodb.NewMongoTransaction(mongoClient, txnOpts)
}*/

// newMongoClient returns the mongo connection
func newMongoClient(ctx context.Context) *mongo.Client {
	connectTimeout := cfg.Env().MongoDbConnectTimeout
	maxConnIdleTime := cfg.Env().MongoDbMaxConnIdleTime
	maxPoolSize := cfg.Env().MongoDbPoolMaxSize
	minPoolSize := cfg.Env().MongoDbPoolMinSize

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	mongoClientOptions := &mongoOptions.ClientOptions{
		ConnectTimeout:  &connectTimeout,
		MaxConnIdleTime: &maxConnIdleTime,
		MaxPoolSize:     &maxPoolSize,
		MinPoolSize:     &minPoolSize,
		TLSConfig:       tlsConfig,
	}

	mongoClientOptions.ApplyURI(cfg.Env().MongoDbURL)

	/*if cfg.Env().Env == "SDX" {
		loggerOptions := mongoOptions.
			Logger().
			SetComponentLevel(mongoOptions.LogComponentCommand, mongoOptions.LogLevelDebug)

		mongoClientOptions.SetLoggerOptions(loggerOptions)
	}*/

	mongoClient, err := mongodb.NewClient(ctx, mongoClientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to mongodb")
	}
	return mongoClient
}

func BuildFeatureFlags(ctx context.Context) flag.FeatureFlag {
	env := cfg.Env().Env
	ff, err := flag.New(ctx, flag.InstanceConfig{
		Environment:           env,
		StrapiToken:           cfg.Env().StrapiToken,
		StrapiURL:             cfg.Env().StrapiURL,
		ApplicationIdentifier: cfg.ApplicationName,
	})
	if err != nil {
		log.Fatal().Err(errors.WithStack(err)).Msg("error creating feature flag client")
	}
	return ff
}

