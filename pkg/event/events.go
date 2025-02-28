package event

const (
	EventLinkVisited = "link.visited"
)

type EventLinkVisitedPayload struct {
	LinkId uint
}
