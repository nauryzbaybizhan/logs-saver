package eventsWeb

import (
	"encoding/json"
	"github.com/didip/tollbooth"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gin-gonic/gin"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
	"github.com/rome314/idkb-events/pkg/logging"
)

type delivery struct {
	publisher         message.Publisher
	fallBackPublisher message.Publisher
	topic             string
	logger            *logging.Entry
}

func NewGinDelivery(logger *logging.Entry, publisher, fallBackPublisher message.Publisher, topic string) *delivery {
	return &delivery{logger: logger, publisher: publisher, fallBackPublisher: fallBackPublisher, topic: topic}
}

func (d *delivery) SetEndpoints(group *gin.RouterGroup) {
	lmt := tollbooth.NewLimiter(50, time.Second)
	group.POST("/event", LimitHandler(lmt), d.handleEvent)
}

func (d *delivery) handleEvent(ctx *gin.Context) {

	logger := d.logger.WithMethod("handleEvent")

	input := eventEntities.RawEvent{}

	if err := ctx.ShouldBind(&input); err != nil {

		logger.WithPlace("read_request").Error(input)
		ctx.Status(http.StatusBadRequest)
		return
	}

	bts, _ := json.Marshal(input)
	msg := message.NewMessage(watermill.NewUUID(), bts)

	if err := d.publisher.Publish(d.topic, msg); err != nil {
		logger.WithPlace("publish_message").Error(err)
		d.fallBackPublisher.Publish(d.topic, msg)
	}
	ctx.Status(http.StatusOK)
	return

}
