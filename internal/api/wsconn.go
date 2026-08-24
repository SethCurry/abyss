package api

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// messageTypePrefixLen is the number of bytes used to encode the MessageType
// and correlation ID at the start of each websocket message.
const messageTypePrefixLen = 8

// wsResponse is a response routed to a pending request.
type wsResponse struct {
	mt   wsmessage.MessageType
	data []byte
}

// wsConn wraps a websocket connection and provides request/response
// correlation plus a demultiplexing read loop. A single reader goroutine
// (serve) reads every message and routes responses to pending requests while
// forwarding requests and notifications to a handler.
type wsConn struct {
	conn    *websocket.Conn
	logger  zerolog.Logger
	writeMu sync.Mutex
	mu      sync.Mutex
	nextID  uint32
	pending map[uint32]chan wsResponse
}

func newWSConn(conn *websocket.Conn, logger zerolog.Logger) *wsConn {
	return &wsConn{
		conn:    conn,
		logger:  logger.With().Str("from", "wsConn").Logger(),
		pending: make(map[uint32]chan wsResponse),
	}
}

// writeMessage writes a MessageType- and correlation-ID-tagged JSON payload to
// the websocket.
func (w *wsConn) writeMessage(mt wsmessage.MessageType, id uint32, data []byte) error {
	w.logger.Debug().Int("mt", int(mt)).Uint32("id", id).Msg("writeMessage")

	buf := make([]byte, messageTypePrefixLen+len(data))
	binary.BigEndian.PutUint32(buf[0:4], uint32(mt))
	binary.BigEndian.PutUint32(buf[4:8], id)
	copy(buf[messageTypePrefixLen:], data)

	w.writeMu.Lock()
	w.logger.Debug().Int("mt", int(mt)).Uint32("id", id).Msg("writeMessage acquired writeMu")
	defer w.writeMu.Unlock()
	err := w.conn.WriteMessage(websocket.BinaryMessage, buf)
	w.logger.Debug().Int("mt", int(mt)).Uint32("id", id).Err(err).Msg("writeMessage wrote")
	return err
}

// readMessage reads a MessageType- and correlation-ID-tagged JSON payload from
// the websocket.
func (w *wsConn) readMessage() (wsmessage.MessageType, uint32, []byte, error) {
	_, data, err := w.conn.ReadMessage()
	if err != nil {
		return wsmessage.MessageTypeNotExist, 0, nil, err
	}
	if len(data) < messageTypePrefixLen {
		return wsmessage.MessageTypeNotExist, 0, nil, fmt.Errorf("websocket message too short: %d bytes", len(data))
	}
	mt := wsmessage.MessageType(binary.BigEndian.Uint32(data[0:4]))
	id := binary.BigEndian.Uint32(data[4:8])
	w.logger.Debug().Int("mt", int(mt)).Uint32("id", id).Msg("readMessage")
	return mt, id, data[messageTypePrefixLen:], nil
}

// request sends a request and waits for the matching response.
func (w *wsConn) request(ctx context.Context, mt wsmessage.MessageType, req any, respType wsmessage.MessageType) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	w.mu.Lock()
	w.nextID++
	id := w.nextID
	ch := make(chan wsResponse, 1)
	w.pending[id] = ch
	pending := len(w.pending)
	w.mu.Unlock()
	w.logger.Debug().Int("mt", int(mt)).Uint32("id", id).Int("pending", pending).Msg("request registered")
	defer func() {
		w.mu.Lock()
		delete(w.pending, id)
		w.mu.Unlock()
		w.logger.Debug().Uint32("id", id).Msg("request pending removed")
	}()

	if err := w.writeMessage(mt, id, data); err != nil {
		w.logger.Error().Err(err).Uint32("id", id).Msg("request write failed")
		return nil, err
	}

	w.logger.Debug().Uint32("id", id).Msg("request waiting for response")
	select {
	case resp := <-ch:
		w.logger.Debug().Uint32("id", id).Int("resp_mt", int(resp.mt)).Msg("request received response")
		if resp.mt != respType {
			return nil, fmt.Errorf("unexpected response type: got %d, want %d", resp.mt, respType)
		}
		return wsmessage.UnmarshalMessage(wsmessage.MessageType(resp.mt), resp.data)
	case <-ctx.Done():
		w.logger.Debug().Uint32("id", id).Msg("request context done")
		return nil, ctx.Err()
	}
}

// notify sends a notification (no response expected).
func (w *wsConn) notify(mt wsmessage.MessageType, data []byte) error {
	return w.writeMessage(mt, 0, data)
}

// serve runs the demultiplexing read loop. Responses are routed to pending
// requests; requests are handled asynchronously (so a blocking stdio forward
// does not stall the reader); notifications are handled inline to preserve
// ordering.
func (w *wsConn) serve(ctx context.Context, handler func(context.Context, wsmessage.MessageType, uint32, []byte) error) error {
	for {
		mt, id, data, err := w.readMessage()
		if err != nil {
			w.logger.Debug().Err(err).Msg("serve read error, exiting read loop")
			return err
		}

		switch classifyMessage(mt) {
		case kindResponse:
			w.mu.Lock()
			ch := w.pending[id]
			w.mu.Unlock()
			if ch != nil {
				w.logger.Debug().Uint32("id", id).Int("type", int(mt)).Msg("routing response to pending request")
				ch <- wsResponse{mt: mt, data: data}
			} else {
				w.logger.Warn().Uint32("id", id).Int("type", int(mt)).Msg("response with no pending request")
			}
		case kindRequest:
			w.logger.Debug().Uint32("id", id).Int("type", int(mt)).Msg("spawning request handler")
			go func() {
				if err := handler(ctx, mt, id, data); err != nil {
					w.logger.Error().Err(err).Int("type", int(mt)).Msg("failed to handle request")
				}
			}()
		case kindNotification:
			w.logger.Debug().Uint32("id", id).Int("type", int(mt)).Msg("handling notification inline")
			if err := handler(ctx, mt, id, data); err != nil {
				return err
			}
		}
	}
}
