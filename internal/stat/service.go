package stat

import "url-shortener/pkg/event"

type StatServiceDeps struct {
	*StatRepository
	*event.EventBus
}

type StatService struct {
	*StatRepository
	*event.EventBus
}

func NewStatService(deps StatServiceDeps) *StatService {
	return &StatService{
		StatRepository: deps.StatRepository,
		EventBus:       deps.EventBus,
	}
}

func (svc *StatService) AddClick() {
	svc.EventBus.On(event.EventLinkVisited, func(payload any) {
		eventPayload, ok := payload.(event.EventLinkVisitedPayload)
		if !ok {
			return
		}

		svc.StatRepository.AddClick(eventPayload.LinkId)
	})
}
