package hasswsclient

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	ws "github.com/moonen-home-automation/hass-ws-client/internal/websocket"
	"github.com/moonen-home-automation/hass-ws-client/pkg/services"
	"log/slog"
	"slices"
	"sync"
	"time"
)

var ErrInvalidToken = ws.ErrInvalidToken

var ErrInvalidArgs = errors.New("invalid arguments provided")

var lock = &sync.Mutex{}

var appInstance *App

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	conn      *websocket.Conn

	wsWriter *ws.WebSocketWriter

	ServiceCaller services.ServiceCaller

	eventListenerTypes []string
}

type DurationString string

type TimeString string

type timeRange struct {
	start time.Time
	end   time.Time
}

type InitializeAppRequest struct {
	URL         string
	HAAuthToken string
	Secure      bool
}

func InitializeAppInstance(request InitializeAppRequest) (*App, error) {
	if request.URL == "" || request.HAAuthToken == "" {
		slog.Error("URL and HAAuthToken are required arguments for InitializeAppRequest")
		return nil, ErrInvalidArgs
	}

	var (
		conn      *websocket.Conn
		ctx       context.Context
		ctxCancel context.CancelFunc
		err       error
	)

	if request.Secure {
		conn, ctx, ctxCancel, err = ws.SetupSecureConnection(request.URL, request.HAAuthToken)
	} else {
		conn, ctx, ctxCancel, err = ws.SetupConnection(request.URL, request.HAAuthToken)
	}

	if conn == nil {
		return nil, err
	}

	wsWriter := &ws.WebSocketWriter{Conn: conn}

	serviceCaller := services.ServiceCaller{Conn: wsWriter, Ctx: ctx}

	appInstance = &App{
		conn:          conn,
		wsWriter:      wsWriter,
		ctx:           ctx,
		ctxCancel:     ctxCancel,
		ServiceCaller: serviceCaller,
	}

	return appInstance, nil
}

func GetAppInstance() *App {
	return appInstance
}

func (a *App) Cleanup() {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
}

func (a *App) RegisterEventListener(listener EventListener) {
	if !slices.Contains(a.eventListenerTypes, listener.EventType) {
		ws.SubscribeToEventType(listener.EventType, a.wsWriter, a.ctx)
		a.eventListenerTypes = append(a.eventListenerTypes, listener.EventType)
	}
}

func (a *App) ListenForEvents(listener EventListener, eventChan chan EventData) {
	elChan := make(chan ws.ChanMsg, 25)
	go ws.ListenWebsocket(a.conn, a.ctx, elChan)

	for {
		msg, ok := <-elChan
		if !ok {
			break
		}

		if msg.Type != "event" {
			continue
		}

		baseEventMsg := BaseEventMsg{}
		_ = json.Unmarshal(msg.Raw, &baseEventMsg)
		if baseEventMsg.Event.EventType != listener.EventType {
			continue
		}
		eventData := EventData{
			Type:         baseEventMsg.Event.EventType,
			RawEventJSON: msg.Raw,
		}
		eventChan <- eventData
	}
}
