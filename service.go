package hasswsclient

import (
	"context"
	"github.com/moonen-home-automation/hass-ws-client/internal/services"
	ws "github.com/moonen-home-automation/hass-ws-client/internal/websocket"
)

type Service struct {
	Event *services.Event
}

func newService(conn *ws.WebSocketWriter, ctx context.Context) *Service {
	return &Service{
		Event: services.BuildService[services.Event](conn, ctx),
	}
}
