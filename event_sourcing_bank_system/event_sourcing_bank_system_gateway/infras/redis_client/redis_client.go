package redis_client

import (
	"context"
	"event_sourcing_bank_system_gateway/package/settings"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedis(cfg *settings.RedisConfig) (*RedisClient, error) {
	var (
		redisClient *redis.Client
		err         error
	)

	if cfg.UseSentinel {
		redisClient, err = newSentinel(cfg)
	} else {
		redisClient, err = newStandAlone(cfg)
	}
	if err != nil {
		return nil, err
	}

	cmd := redisClient.Ping(context.Background())
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		return nil, fmt.Errorf("instrument tracing redis got err=%w", err)
	}
	if err := redisotel.InstrumentMetrics(redisClient); err != nil {
		return nil, fmt.Errorf("instrument metrics redis got err=%w", err)
	}

	return &RedisClient{
		client: redisClient,
	}, nil
}

func newSentinel(cfg *settings.RedisConfig) (*redis.Client, error) {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       cfg.SentinelMasterName,
		SentinelAddrs:    cfg.SentinelServers,
		Password:         cfg.Password,
		SentinelPassword: cfg.SentinelPassword,
		DB:               cfg.DB,
		PoolSize:         cfg.PoolSize,
		DialTimeout:      time.Duration(cfg.DialTimeoutSeconds) * time.Second,
		ReadTimeout:      time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout:     time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
		// Disable unsupported features for compatibility with older Redis versions
		DisableIdentity: true,
		Protocol:        2, // Force RESP2 protocol to avoid newer features
	}), nil
}

func newStandAlone(cfg *settings.RedisConfig) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.ConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("parseURl failed err=%w", err)
	}

	opts.PoolSize = cfg.PoolSize
	opts.DialTimeout = time.Duration(cfg.DialTimeoutSeconds) * time.Second
	opts.ReadTimeout = time.Duration(cfg.ReadTimeoutSeconds) * time.Second
	opts.WriteTimeout = time.Duration(cfg.WriteTimeoutSeconds) * time.Second
	// Disable unsupported features for compatibility with older Redis versions
	opts.DisableIdentity = true
	opts.Protocol = 2 // Force RESP2 protocol to avoid newer features

	return redis.NewClient(opts), nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}
