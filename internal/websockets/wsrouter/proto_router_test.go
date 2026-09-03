package wsrouter_test

// MockProtoRouter is a test double for IProtoRouter that records
// written messages and registered handlers for inspection in tests.
type MockProtoRouter struct {
	WrittenMessages []MockWrittenMessage
	Handlers        map[int]func(ProtoMessage)
	// WriteErr, if set, is returned by WriteMessage instead of nil.
	WriteErr error
}

// MockWrittenMessage captures a single WriteMessage invocation.
type MockWrittenMessage struct {
	TypeID int
	Data   []byte
}

var _ IProtoRouter = &MockProtoRouter{}

func NewMockProtoRouter() *MockProtoRouter {
	return &MockProtoRouter{
		Handlers: make(map[int]func(ProtoMessage)),
	}
}

func (m *MockProtoRouter) WriteMessage(mt int, data []byte) error {
	m.WrittenMessages = append(m.WrittenMessages, MockWrittenMessage{
		TypeID: mt,
		Data:   append([]byte(nil), data...),
	})
	return m.WriteErr
}

func (m *MockProtoRouter) Handle(mt int, handler func(ProtoMessage)) {
	if m.Handlers == nil {
		m.Handlers = make(map[int]func(ProtoMessage))
	}
	m.Handlers[mt] = handler
}

// Dispatch invokes the registered handler for the given message type,
// if any, simulating an inbound message arriving on the router.
func (m *MockProtoRouter) Dispatch(msg ProtoMessage) {
	if h, ok := m.Handlers[msg.TypeID]; ok {
		h(msg)
	}
}
