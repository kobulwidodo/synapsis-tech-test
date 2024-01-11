package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

const (
	Nil = redis.Nil
)

type Interface interface {
	Get(ctx context.Context, key string) (string, error)
	SetEX(ctx context.Context, key string, val string, expTime time.Duration) error
}

type TLSConfig struct {
	Enabled            bool
	InsecureSkipVerify bool
}

type Config struct {
	Protocol string
	Host     string
	Port     string
	Username string
	Password string
	TLS      TLSConfig
}

type cache struct {
	conf  Config
	rdb   *redis.Client
	rlock *redislock.Client
}

func Init(cfg Config) Interface {
	c := &cache{
		conf: cfg,
	}
	c.connect(context.Background())
	return c
}

func (c *cache) connect(ctx context.Context) {
	redisOpts := redis.Options{
		Network:  c.conf.Protocol,
		Addr:     fmt.Sprintf("%s:%s", c.conf.Host, c.conf.Port),
		Username: c.conf.Username,
		Password: c.conf.Password,
	}

	if c.conf.TLS.Enabled {
		redisOpts.TLSConfig = &tls.Config{
			InsecureSkipVerify: c.conf.TLS.InsecureSkipVerify,
		}
	}

	client := redis.NewClient(&redisOpts)

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("[FATAL] cannot connect to redis on address @%s:%v, with error: %s", c.conf.Host, c.conf.Port, err)
	}
	c.rdb = client
	log.Printf("REDIS: Address @%s:%v", c.conf.Host, c.conf.Port)

	c.rlock = redislock.New(client)
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	s, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return s, err
	}

	return s, nil
}

func (c *cache) SetEX(ctx context.Context, key string, val string, expTime time.Duration) error {
	err := c.rdb.SetEx(ctx, key, val, expTime).Err()
	if err != nil {
		return err
	}

	return nil
}
