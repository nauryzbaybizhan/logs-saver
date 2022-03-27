package eventEntities

type UseCase interface {
	HandleEvent(input Event) (err error)
}
