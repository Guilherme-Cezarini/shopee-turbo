package access_token

import (
	"context"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/cfg"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/access_token"
	"github.com/Guilherme-Cezarini/gerenciamento-catalogo-ml-gin/models/responses"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type accessTokenRepository struct {
	mongoDbCollection *mongo.Collection
}

type AccessTokenRepository interface {
	CreateAccessToken(userId string, accessTokenResponse *responses.AccessTokenResponse) error
}

func NewAccessTokenRepository(mongoClient *mongo.Database) AccessTokenRepository {
	mongoDBCollection := mongoClient.Collection(cfg.Env().MongoDbAccessTokenCollection)
	return &accessTokenRepository{
		mongoDbCollection: mongoDBCollection,
	}
}

// CreateAccessToken creates a new access token
func (r *accessTokenRepository) CreateAccessToken(userId string, accessTokenResponse *responses.AccessTokenResponse) error {
	accessTokenModel := access_token.AccessToken{
		AccessTokenResponse: accessTokenResponse,
		UserId:              userId,
	}
	_, err := r.mongoDbCollection.InsertOne(context.Background(), accessTokenModel)
	if err != nil {
		log.Debug().Err(err).Msg("could not insert access token")
		return err
	}

	return nil
}
