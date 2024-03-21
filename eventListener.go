package hasswsclient

type EventListener struct {
	eventType string
}

type EventListenerCallback func(EventData)

type EventData struct {
	Type         string
	RawEventJSON []byte
}

type BaseEventMsg struct {
	Event struct {
		EventType string `json:"event_type"`
	} `json:"event"`
}
