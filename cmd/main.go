package main

import (
	"context"

	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rome314/idkb-events/cmd/config"
	"github.com/rome314/idkb-events/internal/events"
	eventsRepository "github.com/rome314/idkb-events/internal/events/repository/postgres"
	eventsRedisRepository "github.com/rome314/idkb-events/internal/events/repository/redis"
	eventsWeb "github.com/rome314/idkb-events/internal/events/web"
	"github.com/rome314/idkb-events/pkg/connections"
	"github.com/rome314/idkb-events/pkg/logging"
)

func main() {

	logger := logging.GetLogger("main", "")

	logger.Info("Preparing config...")
	cfg := config.GetConfig()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	logger.Info("Preparing DB connections...")
	logger.Info("	Preparing Postgres connection...")
	pgConn, err := connections.GetPostgresDatabase(ctx, cfg.PgConnString)
	emperror.Panic(err)

	logger.Info("	Preparing Redis connection...")
	redisConn, err := connections.GetRedisConnection(ctx, connections.RedisConfig{
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
		Db:       cfg.Redis.Db,
	})
	emperror.Panic(err)

	logger.Info("Configuring internal modules...")
	redisPubSub, err := connections.GetRedisPubSub(ctx, redisConn.Connection, cfg.PubSub)
	emperror.Panic(err)

	channelPubSub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NewStdLogger(false, false))

	eventsRepo := eventsRepository.NewPostgres(logging.GetLogger("events", "repository"), pgConn)
	eventsBufferRepo := eventsRedisRepository.New(redisConn.Connection.(*redis.Client))

	pgInfoProvider := eventsRepository.NewIpInfoManager(logging.GetLogger("events", "ip_info_pg"), pgConn)
	ipInfoProvider := events.NewAggregatedIpProvider(pgInfoProvider)

	ucInput := events.CreateUseCaseInput{
		Repo:           eventsRepo,
		BufferRepo:     eventsBufferRepo,
		Sub:            redisPubSub.Sub,
		FallbackSub:    channelPubSub,
		IpInfoProvider: ipInfoProvider,
		Config:         cfg.App,
	}

	eventsUC := events.NewUseCase(logging.GetLogger("events", "use_case"), ucInput)

	logger.Info("Running main listener...")
	err = eventsUC.Run(ctx)
	emperror.Panic(err)

	eventsGin := eventsWeb.NewGinDelivery(logging.GetLogger("events", "gin"), redisPubSub.Pub, channelPubSub, cfg.App.EventsTopic)

	logger.Info("Configuring web server...")
	router := gin.Default()

	apiGroup := router.Group("/api")

	eventsGin.SetEndpoints(apiGroup)

	logger.Info("Running...")
	err = router.Run("0.0.0.0:" + cfg.ServerPort)
	emperror.Panic(err)

}
