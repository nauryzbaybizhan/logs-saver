package config

import (
	"github.com/rome314/idkb-events/internal/events/config"
	"github.com/rome314/idkb-events/pkg/connections"
	"github.com/spf13/viper"
)

func init() {
	viperInit()
	cfg = &Config{
		ServerPort:   viper.GetString("PORT"),
		PgConnString: viper.GetString("PG_CONN_STRING"),
		Redis: connections.RedisConfig{
			Address:  viper.GetString("REDIS_ADDRESS"),
			Password: viper.GetString("REDIS_PASSWORD"),
			Db:       viper.GetInt("REDIS_DB"),
		},
		App: config.Config{
			EventsTopic: viper.GetString("APP_EVENTS_TOPIC"),
			BufferSize:  viper.GetUint64("APP_BUFFER_SIZE"),
		},
		PubSub: connections.RedisPubSubconfig{
			ConsumerGroup: viper.GetString("PUBSUB_CONSUMER_GROUP"),
		},
	}
}

func viperInit() {

	viper.SetDefault("PORT", "80")

	viper.SetDefault("APP_EVENTS_TOPIC", "events")
	viper.SetDefault("APP_BUFFER_SIZE", 500)
	viper.SetDefault("PUBSUB_CONSUMER_GROUP", "raw_events_group")

	viper.AutomaticEnv()

}
