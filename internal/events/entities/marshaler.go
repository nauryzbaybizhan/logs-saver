package eventEntities

import (
	"bytes"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
)

type RedisMarshaller struct {
}

func (r RedisMarshaller) Marshal(topic string, msg *message.Message) (resp map[string]interface{}, err error) {
	switch topic {
	case EventsTopic, "events_test":
		resp = map[string]interface{}{
			"key":  msg.UUID,
			"data": string(msg.Payload),
		}
		return
	default:
		err = errors.Errorf("unsupported topic: %s", topic)
		return
	}
}

func (r RedisMarshaller) Unmarshal(values map[string]interface{}) (msg *message.Message, err error) {
	uuid, ok := values["key"]
	if !ok {
		err = errors.New("no key in values")
		return
	}
	data, ok := values["data"]
	if !ok {
		err = errors.New("no data in values")
		return
	}

	payload := bytes.NewBufferString(data.(string)).Bytes()

	msg = message.NewMessage(uuid.(string), payload)
	return msg, nil
}
