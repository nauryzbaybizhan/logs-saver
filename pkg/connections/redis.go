package connections

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Address  string
	Password string
	Db       int
}

type RedisConnection struct {
	Connection redis.UniversalClient
}

func (r *RedisConnection) connectionChecker(ctx context.Context) {
	for {
		time.Sleep(time.Second * 5)
		err := r.Connection.Ping(ctx).Err()
		if err != nil {
			log.Warn("Redis not responding to PING...")
		}
	}
}

func GetRedisConnection(ctx context.Context, config RedisConfig) (conn *RedisConnection, err error) {
	options := redis.UniversalOptions{
		Addrs:     []string{config.Address}, // use default Addr
		Password:  config.Password,          // no password set
		DB:        config.Db,
		TLSConfig: &tls.Config{InsecureSkipVerify: true}, // use default DB
	}
	client := redis.NewUniversalClient(&options)
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	conn = &RedisConnection{Connection: client}
	go conn.connectionChecker(ctx)
	return conn, nil
}
