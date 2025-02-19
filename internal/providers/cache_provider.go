package providers

import (
	"context"
	"time"

	"github.com/rbcervilla/redisstore/v9"
	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/modules"
)

// CacheProvider is a port for the applications underlying storage/persistence layer
type CacheProvider struct {
	Session CacheSession
}

// Compile time proof of implementation
var _ ICacheProvider = (*CacheProvider)(nil)

// CacheSession represents a session with the redis storage provider
type CacheSession struct {
	SDK modules.RedisStorageSDK
}

// InitCacheProvider will initialize the storage provider implementation.
func (p *ModuleProviders) InitCacheProvider() error {
	redisSDK := modules.NewRedisStorageSDK(p.Config.RedisAddress)
	p.CacheProvider = &CacheProvider{
		Session: CacheSession{
			SDK: redisSDK,
		},
	}
	return nil
}

func (p *CacheProvider) Get(key string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.Session.SDK.Get(ctx, key)
}

func (p *CacheProvider) GetNoUnmarshal(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.Session.SDK.GetNoUnmarshal(ctx, key)
}

func (p *CacheProvider) Put(key string, value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.Session.SDK.Put(ctx, key, value)
}

func (p *CacheProvider) InitAuthSessionStore() *redisstore.RedisStore {
	sess, err := redisstore.NewRedisStore(context.Background(), p.Session.SDK.Client)
	if err != nil {
		log.Fatal().Msgf("unable to create redis auth session store: %s", err.Error())
	}

	return sess
}
