package modules

import (
	"context"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type KV struct {
	Key  string
	Body string
}

// RedisStorageSDK is a struct that provides methods for interacting with a Redis database.
// It includes a Redis client, which is used to retrieve data from
// any API and store it in the Redis database.
//
// The RedisStorageSDK also provides methods for building a resource tree of a Kubernetes cluster,
// using the Redis store to retreive the data, which is synced by the data-sync server.
// which represents the cluster's resources and their relationships. The resource tree is stored
// in the Redis database and can be retrieved and manipulated using the RedisStorageSDK's methods.
//
// Note: Before using a RedisStorageSDK, you must call the NewRedisStorageSDK function to initialize
// the Redis client and the Kubernetes client.
type RedisStorageSDK struct {
	Client redis.UniversalClient
}

// NewRedisStorageSDK creates a new RedisStorageSDK with a Redis client that connects to the specified Redis address.
// The Redis address can be a single address or a comma-separated list of addresses.
func NewRedisStorageSDK(redisAddr string) RedisStorageSDK {
	sdk := RedisStorageSDK{
		Client: redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: strings.Split(redisAddr, ","),
		}),
	}
	return sdk
}

// Get retrieves the value associated with the specified key from the Redis database. It unmarshals the value into
// an 'any' type and returns it. If the key does not exist in the database, the method logs a warning and returns
// nil. If an error occurs while retrieving or unmarshalling the value, the method logs an error and returns the error.
func (sdk *RedisStorageSDK) Get(ctx context.Context, key string) (any, error) {
	val, err := sdk.Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Warn().Msgf("%s has not been added in the cache yet: %v", key, err)
		} else {
			log.Error().Msgf("Error getting key %v from cache: %v", key, err)
			return nil, err
		}
	}

	var data any
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(val, &data)
	if err != nil {
		log.Error().Msgf("Error unmarshalling body from redis lookup: %v", err)
		return nil, err
	}
	return data, nil
}

// GetNoUnmarshal retrieves the value associated with the specified key from the Redis database and returns it as a byte slice.
// If the key does not exist in the database, the method logs a warning and returns nil. If an error occurs while retrieving
// the value, the method logs an error and returns the error.
func (sdk *RedisStorageSDK) GetNoUnmarshal(ctx context.Context, key string) ([]byte, error) {
	val, err := sdk.Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Warn().Msgf("%s has not been added in the cache yet: %v", key, err)
		} else {
			log.Error().Msgf("Error getting key %v from cache: %v", key, err)
			return nil, err
		}
	}
	return val, nil
}

// Put marshals the provided body into a byte slice and stores it in the Redis database with the specified key. If an error
// occurs while marshalling the body or storing the value, the method logs an error and returns the error.
func (sdk *RedisStorageSDK) Put(ctx context.Context, key string, body any) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(&body)
	if err != nil {
		log.Error().Msgf("Error marshalling body for redis: %v", err)
		return err
	}

	_, err = sdk.Client.Set(ctx, key, b, 0).Result()
	if err != nil {
		log.Error().Msgf("Error putting key %v in cache: %v", key, err)
	}
	return err
}
