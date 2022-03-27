package eventsRedisRepository

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

func hash(values ...string) string {
	h := md5.New()
	for _, v := range values {
		h.Write([]byte(v))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func getKeyValue(event *eventEntities.Event) (key string, value string) {
	key = hash(
		event.Url,
		event.UserId,
		event.Ip.String(),
		event.ApiKey,
		event.UserAgent,
		fmt.Sprintf("%d", event.RequestTime.Unix()),
	)
	bts, _ := json.Marshal(event)
	value = string(bts)
	return
}
