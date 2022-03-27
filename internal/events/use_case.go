package events

import (
	"context"
	"encoding/json"
	"github.com/rome314/idkb-events/internal/events/config"
	"github.com/rome314/idkb-events/internal/events/repository"
	"github.com/rome314/idkb-events/internal/events/util"
	"sync"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rome314/idkb-events/internal/events/entities"
	"github.com/rome314/idkb-events/pkg/logging"
)

const numWorkers = 10

type uc struct {
	logger     *logging.Entry
	repo       repository.Repository
	bufferRepo repository.BufferRepo

	sub         message.Subscriber
	fallbackSub message.Subscriber

	ipInfoProvider repository.IpInfoProvider

	cfg config.Config

	clearInProgress bool
	mx              *sync.Mutex
}

type CreateUseCaseInput struct {
	Repo       repository.Repository
	BufferRepo repository.BufferRepo

	Sub         message.Subscriber
	FallbackSub message.Subscriber

	IpInfoProvider repository.IpInfoProvider

	Config config.Config
}

func NewUseCase(logger *logging.Entry, input CreateUseCaseInput) *uc {
	return &uc{
		logger:         logger,
		repo:           input.Repo,
		bufferRepo:     input.BufferRepo,
		sub:            input.Sub,
		fallbackSub:    input.FallbackSub,
		ipInfoProvider: input.IpInfoProvider,
		cfg:            input.Config,
		mx:             &sync.Mutex{},
	}
}

func (u *uc) Run(ctx context.Context) error {

	for i := 0; i < numWorkers; i++ {
		messages, err := u.sub.Subscribe(ctx, u.cfg.EventsTopic)
		if err != nil {
			return errors.WithMessage(err, "subscribing to topic")
		}
		fallbackMessage, err := u.fallbackSub.Subscribe(ctx, u.cfg.EventsTopic)
		if err != nil {
			return errors.WithMessage(err, "subscribing to fallback topic")
		}
		go u.listener(ctx, messages, fallbackMessage)
	}
	return nil
}

func (u *uc) listener(ctx context.Context, msgs, fbMsgs <-chan *message.Message) {
	// logger := u.logger.WithMethod("listener")
	for {
		select {
		case _ = <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				u.logger.WithMethod("listener").WithPlace("default chan").Error("channel closed")
				return
			}
			u.handleMessage(msg)
		case msg, ok := <-fbMsgs:
			if !ok {
				u.logger.WithMethod("listener").WithPlace("fallback chan").Error("channel closed")
				return
			}
			u.handleMessage(msg)
		}
	}
}

func (u *uc) handleMessage(msg *message.Message) {
	logger := u.logger.WithMethod("handleMessage")

	rawEvent := eventEntities.RawEvent{}

	if err := json.Unmarshal(msg.Payload, &rawEvent); err != nil {
		logger.WithPlace("read_message").Error(err)
		msg.Ack()
		return
	}

	event, err := util.RawToEvent(rawEvent)
	if err != nil {
		logger.WithPlace("validate_message").Error(err)
		msg.Ack()
		return
	}

	ipInfo, err := u.ipInfoProvider.GetIpInfo(event.Ip.String())
	if err != nil {
		logger.WithPlace("get_ip_info").Error(err)
	}

	event.IpInfo = ipInfo

	if err = u.bufferRepo.Status(); err != nil {
		logger.WithPlace("insert_event").Warn("buffer is off, inserting to pg...")
		e := u.repo.Store(event)
		if e != nil {
			logger.WithPlace("insert_event").Error(e)
			msg.Nack()
			return
		}
	}

	bufferSize, err := u.bufferRepo.Store(event)
	if err != nil {
		logger.WithPlace("insert_event").Warnf("error inserting to buffer (%s), insert to pg... ", err.Error())
		e := u.repo.Store(event)
		if e != nil {
			msg.Nack()
			return
		}
	}
	msg.Ack()
	if bufferSize >= u.cfg.BufferSize && !u.clearInProgress {
		err = u.clearBuffer()
		if err != nil {
			logger.WithPlace("clearBuffer").Error(err)
		}
	}
	return

}

func (u *uc) clearBuffer() error {
	u.mx.Lock()
	u.clearInProgress = true
	logger := u.logger.WithMethod("clearBuffer")
	logger.Info("Starting...")
	defer logger.Info("Finish...")
	defer func() {
		u.clearInProgress = false
		u.mx.Unlock()
	}()

	if u.repo.Status() != nil {
		logger.Warn("Pg unresponding, skipping...")
		return nil
	}

	events, err := u.bufferRepo.PopAll()
	if err != nil {
		return errors.WithMessage(err, "getting events from buffer")
	}
	logger.Infof("Poped %d events from buffer", len(events))

	logger = logger.WithPlace("insert_to_db")
	inserted, err := u.repo.StoreMany(events...)
	if err != nil {
		logger.Warnf("error occured return events to buffer")
		if e := u.bufferRepo.StoreToErrorStorage(events); e != nil {
			logger.Errorf("insert to errors buffer: %s", e.Error())
		}
		return errors.WithMessage(err, "inserting events to main db")
	}
	logger.Infof("inserted %d events", inserted)

	if inserted == 0 {
		logger.Warnf("0 events inserted, return events to buffer")
		if e := u.bufferRepo.StoreToErrorStorage(events); e != nil {
			logger.Errorf("insert to errors buffer: %s", e.Error())

		}

	}
	return nil
}
