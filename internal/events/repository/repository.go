package repository

import "github.com/rome314/idkb-events/internal/events/entities"

type IpInfoManager interface {
	IpInfoProvider
	SetIpInfo(ip string, info *eventEntities.IpInfo) (id int32, err error)
}

type IpInfoProvider interface {
	GetIpInfo(ip string) (info *eventEntities.IpInfo, err error)
}

type BufferRepo interface {
	Store(event *eventEntities.Event) (bufferSize uint64, err error)
	StoreToErrorStorage(events []*eventEntities.Event) (err error)
	PopAll() (events []*eventEntities.Event, err error)
	Status() error
}

type Repository interface {
	Store(event *eventEntities.Event) (err error)
	StoreMany(events ...*eventEntities.Event) (insertedCount int64, err error)
	Status() error
}
