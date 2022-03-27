package connections

import (
	"context"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-redis/redis/v8"
	redisPubSub "github.com/minghsu0107/watermill-redistream/pkg/redis"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

type RedisPubSub struct {
	Pub message.Publisher
	Sub message.Subscriber
}

type RedisPubSubconfig struct {
	ConsumerGroup string
}

func GetRedisPubSub(ctx context.Context, client redis.UniversalClient, cfg RedisPubSubconfig) (pubSub *RedisPubSub, err error) {
	pubSubMarshaler := eventEntities.RedisMarshaller{}
	sub, err := redisPubSub.NewSubscriber(
		ctx,
		redisPubSub.SubscriberConfig{Consumer: watermill.NewShortUUID(), ConsumerGroup: cfg.ConsumerGroup},
		client,
		pubSubMarshaler,
		nil,
	)
	if err != nil {
		err = errors.WithMessage(err, "creating sub")
		return
	}

	pub, err := redisPubSub.NewPublisher(
		ctx,
		redisPubSub.PublisherConfig{
			Maxlens: map[string]int64{
				"events": 50000,
			},
		},
		client,
		pubSubMarshaler,
		nil,
	)
	if err != nil {
		err = errors.WithMessage(err, "creating pub")
		return
	}
	return &RedisPubSub{
		Pub: pub,
		Sub: sub,
	}, nil

}
