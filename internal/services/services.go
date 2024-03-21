package services

import (
	"context"
	"github.com/moonen-home-automation/hass-ws-client/internal"

	ws "github.com/moonen-home-automation/hass-ws-client/internal/websocket"
)

func BuildService[T Event](conn *ws.WebSocketWriter, ctx context.Context) *T {
	return &T{conn: conn, ctx: ctx}
}

type BaseServiceRequest struct {
	Id          int64          `json:"id"`
	RequestType string         `json:"type"`
	Domain      string         `json:"domain"`
	Service     string         `json:"service"`
	ServiceData map[string]any `json:"service_data,omitempty"`
	Target      struct {
		EntityId string `json:"entity_id,omitempty"`
	} `json:"target,omitempty"`
}

func NewBaseServiceRequest(entityId string) BaseServiceRequest {
	id := internal.GetId()
	bsr := BaseServiceRequest{
		Id:          id,
		RequestType: "call_service",
	}
	if entityId != "" {
		bsr.Target.EntityId = entityId
	}
	return bsr
}
