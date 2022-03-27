package events

import (
	"emperror.dev/errors"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
	"github.com/rome314/idkb-events/internal/events/repository"
	eventsSupport "github.com/rome314/idkb-events/internal/events/support"
	log "github.com/sirupsen/logrus"
)

type ipProvider struct {
	manager     repository.IpInfoManager
	apiProvider repository.IpInfoProvider
}

func (i *ipProvider) GetIpInfo(ip string) (info *eventEntities.IpInfo, err error) {
	if info, err = i.manager.GetIpInfo(ip); err == nil {
		return
	}
	info, err = i.apiProvider.GetIpInfo(ip)
	if err != nil {
		log.Info("api")
		err = errors.WithMessage(err, "api error")
		return
	}
	id, err := i.manager.SetIpInfo(ip, info)
	if err != nil {
		err = errors.WithMessage(err, "inserting to manager")
	}
	info.Id = id
	return
}

func NewAggregatedIpProvider(manager repository.IpInfoManager) repository.IpInfoProvider {
	return &ipProvider{
		manager:     manager,
		apiProvider: eventsSupport.NewApiInfoProvider(),
	}
}
