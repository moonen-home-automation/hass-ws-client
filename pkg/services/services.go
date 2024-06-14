package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/moonen-home-automation/hass-ws-client/internal"

	ws "github.com/moonen-home-automation/hass-ws-client/internal/websocket"
)

type ServiceCall struct {
	Id          int64                  `json:"id"`
	RequestType string                 `json:"type"`
	Domain      string                 `json:"domain"`
	Service     string                 `json:"service"`
	ServiceData map[string]interface{} `json:"service_data,omitempty"`
	Target      ServiceTarget          `json:"target"`
	Returns     bool                   `json:"return_response"`
}

type ServiceResponse struct {
	Result struct {
		Response any `json:"response"`
	} `json:"result"`
}

type ServiceTarget struct {
	AreaID   string `json:"area_id"`
	DeviceID string `json:"device_id"`
	EntityID string `json:"entity_id"`
	LabelID  string `json:"label_id"`
}

type ServiceCaller struct {
	Conn *ws.WebSocketWriter
	Ctx  context.Context
}

func (s *ServiceCaller) Call(service ServiceCall) (ServiceResponse, error) {
	err := s.Conn.WriteMessage(service, s.Ctx)
	if err != nil {
		return ServiceResponse{}, err
	}
	if service.Returns == true {
		resp := listenForServiceResponse(s.Conn.Conn, s.Ctx, service.Id)
		return *resp, nil
	}
	return ServiceResponse{}, nil
}

func listenForServiceResponse(conn *websocket.Conn, ctx context.Context, id int64) *ServiceResponse {
	elChan := make(chan ws.ChanMsg)
	go ws.ListenWebsocket(conn, ctx, elChan)

	var serviceResponse *ServiceResponse

	for {
		msg, ok := <-elChan
		if !ok {
			break
		}

		if msg.Type != "result" {
			continue
		}

		if msg.Id != id {
			continue
		}

		srvResp := ServiceResponse{}
		fmt.Println(string(msg.Raw))
		_ = json.Unmarshal(msg.Raw, &srvResp)
		serviceResponse = &srvResp
		break
	}

	return serviceResponse
}

func NewServiceCall(domain string, service string, data map[string]interface{}, target ServiceTarget, returns bool) ServiceCall {
	id := internal.GetId()
	sc := ServiceCall{
		Id:          id,
		RequestType: "call_service",
		Domain:      domain,
		Service:     service,
		ServiceData: data,
		Target:      target,
		Returns:     returns,
	}
	return sc
}
