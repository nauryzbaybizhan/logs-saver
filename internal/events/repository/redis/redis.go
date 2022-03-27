package eventsRedisRepository

import (
	"context"
	"encoding/json"
	"github.com/rome314/idkb-events/internal/events/repository"
	"sync"

	"emperror.dev/errors"
	"github.com/go-redis/redis/v8"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
	log "github.com/sirupsen/logrus"
)

type repo struct {
	client *redis.Client
}

func New(client *redis.Client) repository.BufferRepo {
	return &repo{
		client: client,
	}
}

func (r *repo) StoreToErrorStorage(events []*eventEntities.Event) (err error) {

	values := map[string]interface{}{}

	wg := &sync.WaitGroup{}
	wg.Add(len(events))

	mx := &sync.Mutex{}

	for _, event := range events {
		go func(e *eventEntities.Event) {
			defer wg.Done()
			key, value := getKeyValue(e)
			mx.Lock()
			values[key] = value
			mx.Unlock()

		}(event)
	}
	wg.Wait()

	if err = r.client.HSet(context.TODO(), "uninserteds_buffer", values).Err(); err != nil {
		return errors.WithMessage(err, "inserting")
	}
	return nil

}

func (r *repo) Count() (count uint64, err error) {
	count, err = r.client.HLen(context.TODO(), "insert_buffer").Uint64()
	if err != nil {
		err = errors.WithMessage(err, "getting count")
		return
	}
	return count, nil
}

func (r *repo) PopAll() (events []*eventEntities.Event, err error) {
	// count, err := r.Count()
	// if err != nil {
	// 	err = errors.WithMessage(err, "getting count")
	// 	return
	// }

	kvs, err := r.client.HGetAll(context.TODO(), "insert_buffer").Result()
	if err != nil {
		err = errors.WithMessage(err, "getting")
	}

	events = make([]*eventEntities.Event, len(kvs))
	keys := make([]string, len(kvs))
	wg := &sync.WaitGroup{}
	wg.Add(len(kvs))

	index := -1
	for key, value := range kvs {
		index++
		go func(i int, key, encoded string) {
			defer wg.Done()
			tmp := &eventEntities.Event{}
			if e := json.Unmarshal([]byte(encoded), tmp); e != nil {
				log.Error(e)
				return
			}
			events[i] = tmp
			keys[i] = key

		}(index, key, value)
	}
	wg.Wait()

	if err = r.client.HDel(context.TODO(), "insert_buffer", keys...).Err(); err != nil {
		err = errors.WithMessage(err, "clearing hash")
		return
	}

	return events, nil

}
func (r *repo) Store(event *eventEntities.Event) (bufferSize uint64, err error) {
	// bts, err := json.Marshal(event)
	// if err != nil {
	// 	err = errors.WithMessage(err, "marshalling")
	// 	return
	// }

	k, v := getKeyValue(event)

	if err = r.client.HSetNX(context.TODO(), "insert_buffer", k, v).Err(); err != nil {
		err = errors.WithMessage(err, "inserting")
		return
	}

	return r.Count()
}

func (r *repo) Status() error {
	return r.client.Ping(context.TODO()).Err()
}
