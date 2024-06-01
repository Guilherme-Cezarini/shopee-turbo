package cfg

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/verifymycontent/go-env"
	"time"
)

type Environment struct {
	Env                        string        `env:"ENV"`
	LogLevel                   string        `env:"LOG_LEVEL"`
	LogJSON                    bool          `env:"LOG_JSON"`
	LogComponent               string        `env:"LOG_COMPONENT"`
	HttpServerPort             string        `env:"HTTP_SERVER_PORT"`
	TracingEnabled             bool          `env:"TRACING_ENABLED"`
	HashingFirst               string        `env:"HASHING_FIRST"`
	HashingSecond              string        `env:"HASHING_SECOND"`
	HashingThird               string        `env:"HASHING_THIRD"`
	HashingFourth              string        `env:"HASHING_FOURTH"`
	MongoDbURL                 string        `env:"MONGODB_URL"`
	MongoDbDatabase            string        `env:"MONGODB_DATABASE"`
	MongoDbUsersCollection     string        `env:"MONGODB_USERS_COLLECTION"`
	MongoDbPoolMinSize         uint64        `env:"MONGODB_POOL_MIN_SIZE"`
	MongoDbPoolMaxSize         uint64        `env:"MONGODB_POOL_MAX_SIZE"`
	MongoDbMaxConnIdleTime     time.Duration `env:"MONGODB_MAX_CONN_IDLE_TIME"`
	MongoDbConnectTimeout      time.Duration `env:"MONGODB_CONNECT_TIMEOUT"`
	MongoDbCollectionUsers     string        `env:"MONGODB_COLLECTION_USERS"`
	MongoAccessTokenCollection string        `env:"MONGODB_ACCESS_TOKEN_COLLECTION"`
	MercadoLivreAuthUrl        string        `env:"MERCADO_LIVRE_AUTH_URL"`
	MercadoLivreClientId       string        `env:"MERCADO_LIVRE_CLIENT_ID"`
	MercadoLivreClientSecret   string        `env:"MERCADO_LIVRE_CLIENT_SECRET"`
	MercadoLivreRedirectUri    string        `env:"MERCADO_LIVRE_REDIRECT_URI"`
	Extras                     env.EnvSet
	Loaded                     bool
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msgf("error loading .env file")
	}

	// Load environment singleton
	Env()
}

var environmentSingleton *Environment

func Env() *Environment {
	if environmentSingleton == nil {
		environmentSingleton = new(Environment)
		es, err := env.UnmarshalFromEnviron(environmentSingleton)
		if err != nil {
			log.Fatal().Err(err).Msgf("error loading environment")
		}
		environmentSingleton.Extras = es
	}
	return environmentSingleton
}
