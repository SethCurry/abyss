package wsmessage

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/coder/acp-go-sdk"
)

type ACPMessage struct {
	TypeID    MessageType
	Type      reflect.Type
	Unmarshal func([]byte) (any, error)
}

func GetMessageTypeByID(msgTypeID int32) (ACPMessage, error) {
	for _, v := range AllMessages {
		if int32(v.TypeID) == msgTypeID {
			return v, nil
		}
	}

	return ACPMessage{}, fmt.Errorf("no ACP message type with ID %d", msgTypeID)
}

func GetMessageTypeByType(msg any) (ACPMessage, error) {
	msgType := reflect.TypeOf(msg)
	for _, v := range AllMessages {
		if msgType == v.Type {
			return v, nil
		}
	}

	return ACPMessage{}, fmt.Errorf("unknown router message type %T", msg)
}

var AllMessages = []ACPMessage{
	{
		TypeID: RequestPermissionRequestType,
		Type:   reflect.TypeOf(acp.RequestPermissionRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.RequestPermissionRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: RequestPermissionResponseType,
		Type:   reflect.TypeOf(acp.RequestPermissionResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.RequestPermissionResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: WriteTextFileRequestType,
		Type:   reflect.TypeOf(acp.WriteTextFileRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.WriteTextFileRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: WriteTextFileResponseType,
		Type:   reflect.TypeOf(acp.WriteTextFileResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.WriteTextFileResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: ReadTextFileRequestType,
		Type:   reflect.TypeOf(acp.ReadTextFileRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.ReadTextFileRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: ReadTextFileResponseType,
		Type:   reflect.TypeOf(acp.ReadTextFileResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.ReadTextFileResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: CreateTerminalRequestType,
		Type:   reflect.TypeOf(acp.CreateTerminalRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.CreateTerminalRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: CreateTerminalResponseType,
		Type:   reflect.TypeOf(acp.CreateTerminalResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.CreateTerminalResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: TerminalOutputRequestType,
		Type:   reflect.TypeOf(acp.TerminalOutputRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.TerminalOutputRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: TerminalOutputResponseType,
		Type:   reflect.TypeOf(acp.TerminalOutputResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.TerminalOutputResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: ReleaseTerminalRequestType,
		Type:   reflect.TypeOf(acp.ReleaseTerminalRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.ReleaseTerminalRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: ReleaseTerminalResponseType,
		Type:   reflect.TypeOf(acp.ReleaseTerminalResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.ReleaseTerminalResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: WaitForTerminalExitRequestType,
		Type:   reflect.TypeOf(acp.WaitForTerminalExitRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.WaitForTerminalExitRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: WaitForTerminalExitResponseType,
		Type:   reflect.TypeOf(acp.WaitForTerminalExitResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.WaitForTerminalExitResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: KillTerminalRequestType,
		Type:   reflect.TypeOf(acp.KillTerminalRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.KillTerminalRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: KillTerminalResponseType,
		Type:   reflect.TypeOf(acp.KillTerminalResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.KillTerminalResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: SessionNotificationType,
		Type:   reflect.TypeOf(acp.SessionNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.SessionNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: SetSessionModeRequestType,
		Type:   reflect.TypeOf(acp.SetSessionModeRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.SetSessionModeRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: SetSessionModeResponseType,
		Type:   reflect.TypeOf(acp.SetSessionModeResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.SetSessionModeResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableForkSessionRequestType,
		Type:   reflect.TypeOf(acp.UnstableForkSessionRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableForkSessionRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableForkSessionResponseType,
		Type:   reflect.TypeOf(acp.UnstableForkSessionResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableForkSessionResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: ListSessionsRequestType,
		Type:   reflect.TypeOf(acp.ListSessionsRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.ListSessionsRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: ListSessionsResponseType,
		Type:   reflect.TypeOf(acp.ListSessionsResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.ListSessionsResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: ResumeSessionRequestType,
		Type:   reflect.TypeOf(acp.ResumeSessionRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.ResumeSessionRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: ResumeSessionResponseType,
		Type:   reflect.TypeOf(acp.ResumeSessionResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.ResumeSessionResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: SetSessionConfigOptionRequestType,
		Type:   reflect.TypeOf(acp.SetSessionConfigOptionRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.SetSessionConfigOptionRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: SetSessionConfigOptionResponseType,
		Type:   reflect.TypeOf(acp.SetSessionConfigOptionResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.SetSessionConfigOptionResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: LogoutRequestType,
		Type:   reflect.TypeOf(acp.LogoutRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.LogoutRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: LogoutResponseType,
		Type:   reflect.TypeOf(acp.LogoutResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.LogoutResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableCloseNesRequestType,
		Type:   reflect.TypeOf(acp.UnstableCloseNesRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableCloseNesRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableCloseNesResponseType,
		Type:   reflect.TypeOf(acp.UnstableCloseNesResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableCloseNesResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableStartNesRequestType,
		Type:   reflect.TypeOf(acp.UnstableStartNesRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableStartNesRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableStartNesResponseType,
		Type:   reflect.TypeOf(acp.UnstableStartNesResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableStartNesResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableSuggestNesRequestType,
		Type:   reflect.TypeOf(acp.UnstableSuggestNesRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableSuggestNesRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableSuggestNesResponseType,
		Type:   reflect.TypeOf(acp.UnstableSuggestNesResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableSuggestNesResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableAcceptNesNotificationType,
		Type:   reflect.TypeOf(acp.UnstableAcceptNesNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableAcceptNesNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableRejectNesNotificationType,
		Type:   reflect.TypeOf(acp.UnstableRejectNesNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableRejectNesNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDidChangeDocumentNotificationType,
		Type:   reflect.TypeOf(acp.UnstableDidChangeDocumentNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDidChangeDocumentNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDidCloseDocumentNotificationType,
		Type:   reflect.TypeOf(acp.UnstableDidCloseDocumentNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDidCloseDocumentNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDidFocusDocumentNotificationType,
		Type:   reflect.TypeOf(acp.UnstableDidFocusDocumentNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDidFocusDocumentNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDidOpenDocumentNotificationType,
		Type:   reflect.TypeOf(acp.UnstableDidOpenDocumentNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDidOpenDocumentNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDidSaveDocumentNotificationType,
		Type:   reflect.TypeOf(acp.UnstableDidSaveDocumentNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDidSaveDocumentNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDisableProviderRequestType,
		Type:   reflect.TypeOf(acp.UnstableDisableProviderRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDisableProviderRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDisableProviderResponseType,
		Type:   reflect.TypeOf(acp.UnstableDisableProviderResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDisableProviderResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableListProvidersRequestType,
		Type:   reflect.TypeOf(acp.UnstableListProvidersRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableListProvidersRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableListProvidersResponseType,
		Type:   reflect.TypeOf(acp.UnstableListProvidersResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableListProvidersResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableSetProviderRequestType,
		Type:   reflect.TypeOf(acp.UnstableSetProviderRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableSetProviderRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableSetProviderResponseType,
		Type:   reflect.TypeOf(acp.UnstableSetProviderResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableSetProviderResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDeleteSessionRequestType,
		Type:   reflect.TypeOf(acp.UnstableDeleteSessionRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDeleteSessionRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: UnstableDeleteSessionResponseType,
		Type:   reflect.TypeOf(acp.UnstableDeleteSessionResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.UnstableDeleteSessionResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: CloseSessionRequestType,
		Type:   reflect.TypeOf(acp.CloseSessionRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.CloseSessionRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: CloseSessionResponseType,
		Type:   reflect.TypeOf(acp.CloseSessionResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.CloseSessionResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: InitializeRequestType,
		Type:   reflect.TypeOf(acp.InitializeRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.InitializeRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: InitializeResponseType,
		Type:   reflect.TypeOf(acp.InitializeResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.InitializeResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: NewSessionRequestType,
		Type:   reflect.TypeOf(acp.NewSessionRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.NewSessionRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: NewSessionResponseType,
		Type:   reflect.TypeOf(acp.NewSessionResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.NewSessionResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: AuthenticateRequestType,
		Type:   reflect.TypeOf(acp.AuthenticateRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.AuthenticateRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: AuthenticateResponseType,
		Type:   reflect.TypeOf(acp.AuthenticateResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.AuthenticateResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: LoadSessionRequestType,
		Type:   reflect.TypeOf(acp.LoadSessionRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.LoadSessionRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: LoadSessionResponseType,
		Type:   reflect.TypeOf(acp.LoadSessionResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.LoadSessionResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: PromptRequestType,
		Type:   reflect.TypeOf(acp.PromptRequest{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.PromptRequest

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: PromptResponseType,
		Type:   reflect.TypeOf(acp.PromptResponse{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.PromptResponse

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
	{
		TypeID: CancelNotificationType,
		Type:   reflect.TypeOf(acp.CancelNotification{}),
		Unmarshal: func(content []byte) (any, error) {
			var resp acp.CancelNotification

			err := json.Unmarshal(content, &resp)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	},
}

var MessageTypeToMessage = map[MessageType]func() any{
	RequestPermissionRequestType: func() any {
		return &acp.RequestPermissionRequest{}
	},
	RequestPermissionResponseType: func() any {
		return &acp.RequestPermissionResponse{}
	},
	WriteTextFileRequestType: func() any {
		return &acp.WriteTextFileRequest{}
	},
	WriteTextFileResponseType: func() any {
		return &acp.WriteTextFileResponse{}
	},
	ReadTextFileRequestType: func() any {
		return &acp.ReadTextFileRequest{}
	},
	ReadTextFileResponseType: func() any {
		return &acp.ReadTextFileResponse{}
	},
	CreateTerminalRequestType: func() any {
		return &acp.CreateTerminalRequest{}
	},
	CreateTerminalResponseType: func() any {
		return &acp.CreateTerminalResponse{}
	},
	TerminalOutputRequestType: func() any {
		return &acp.TerminalOutputRequest{}
	},
	TerminalOutputResponseType: func() any {
		return &acp.TerminalOutputResponse{}
	},
	ReleaseTerminalRequestType: func() any {
		return &acp.ReleaseTerminalRequest{}
	},
	ReleaseTerminalResponseType: func() any {
		return &acp.ReleaseTerminalResponse{}
	},
	WaitForTerminalExitRequestType: func() any {
		return &acp.WaitForTerminalExitRequest{}
	},
	WaitForTerminalExitResponseType: func() any {
		return &acp.WaitForTerminalExitResponse{}
	},
	KillTerminalRequestType: func() any {
		return &acp.KillTerminalRequest{}
	},
	KillTerminalResponseType: func() any {
		return &acp.KillTerminalResponse{}
	},
	SessionNotificationType: func() any {
		return &acp.SessionNotification{}
	},
	SetSessionModeRequestType: func() any {
		return &acp.SetSessionModeRequest{}
	},
	SetSessionModeResponseType: func() any {
		return &acp.SetSessionModeResponse{}
	},
	UnstableForkSessionRequestType: func() any {
		return &acp.UnstableForkSessionRequest{}
	},
	UnstableForkSessionResponseType: func() any {
		return &acp.UnstableForkSessionResponse{}
	},
	ListSessionsRequestType: func() any {
		return &acp.ListSessionsRequest{}
	},
	ListSessionsResponseType: func() any {
		return &acp.ListSessionsResponse{}
	},
	ResumeSessionRequestType: func() any {
		return &acp.ResumeSessionRequest{}
	},
	ResumeSessionResponseType: func() any {
		return &acp.ResumeSessionResponse{}
	},
	SetSessionConfigOptionRequestType: func() any {
		return &acp.SetSessionConfigOptionRequest{}
	},
	SetSessionConfigOptionResponseType: func() any {
		return &acp.SetSessionConfigOptionResponse{}
	},
	LogoutRequestType: func() any {
		return &acp.LogoutRequest{}
	},
	LogoutResponseType: func() any {
		return &acp.LogoutResponse{}
	},
	UnstableCloseNesRequestType: func() any {
		return &acp.UnstableCloseNesRequest{}
	},
	UnstableCloseNesResponseType: func() any {
		return &acp.UnstableCloseNesResponse{}
	},
	UnstableStartNesRequestType: func() any {
		return &acp.UnstableStartNesRequest{}
	},
	UnstableStartNesResponseType: func() any {
		return &acp.UnstableStartNesResponse{}
	},
	UnstableSuggestNesRequestType: func() any {
		return &acp.UnstableSuggestNesRequest{}
	},
	UnstableSuggestNesResponseType: func() any {
		return &acp.UnstableSuggestNesResponse{}
	},
	UnstableAcceptNesNotificationType: func() any {
		return &acp.UnstableAcceptNesNotification{}
	},
	UnstableRejectNesNotificationType: func() any {
		return &acp.UnstableRejectNesNotification{}
	},
	UnstableDidChangeDocumentNotificationType: func() any {
		return &acp.UnstableDidChangeDocumentNotification{}
	},
	UnstableDidCloseDocumentNotificationType: func() any {
		return &acp.UnstableDidCloseDocumentNotification{}
	},
	UnstableDidFocusDocumentNotificationType: func() any {
		return &acp.UnstableDidFocusDocumentNotification{}
	},
	UnstableDidOpenDocumentNotificationType: func() any {
		return &acp.UnstableDidOpenDocumentNotification{}
	},
	UnstableDidSaveDocumentNotificationType: func() any {
		return &acp.UnstableDidSaveDocumentNotification{}
	},
	UnstableDisableProviderRequestType: func() any {
		return &acp.UnstableDisableProviderRequest{}
	},
	UnstableDisableProviderResponseType: func() any {
		return &acp.UnstableDisableProviderResponse{}
	},
	UnstableListProvidersRequestType: func() any {
		return &acp.UnstableListProvidersRequest{}
	},
	UnstableListProvidersResponseType: func() any {
		return &acp.UnstableListProvidersResponse{}
	},
	UnstableSetProviderRequestType: func() any {
		return &acp.UnstableSetProviderRequest{}
	},
	UnstableSetProviderResponseType: func() any {
		return &acp.UnstableSetProviderResponse{}
	},
	UnstableDeleteSessionRequestType: func() any {
		return &acp.UnstableDeleteSessionRequest{}
	},
	UnstableDeleteSessionResponseType: func() any {
		return &acp.UnstableDeleteSessionResponse{}
	},
	CloseSessionRequestType: func() any {
		return &acp.CloseSessionRequest{}
	},
	CloseSessionResponseType: func() any {
		return &acp.CloseSessionResponse{}
	},
	InitializeRequestType: func() any {
		return &acp.InitializeRequest{}
	},
	InitializeResponseType: func() any {
		return &acp.InitializeResponse{}
	},
	NewSessionRequestType: func() any {
		return &acp.NewSessionRequest{}
	},
	NewSessionResponseType: func() any {
		return &acp.NewSessionResponse{}
	},
	AuthenticateRequestType: func() any {
		return &acp.AuthenticateRequest{}
	},
	AuthenticateResponseType: func() any {
		return &acp.AuthenticateResponse{}
	},
	LoadSessionRequestType: func() any {
		return &acp.LoadSessionRequest{}
	},
	LoadSessionResponseType: func() any {
		return &acp.LoadSessionResponse{}
	},
	PromptRequestType: func() any {
		return &acp.PromptRequest{}
	},
	PromptResponseType: func() any {
		return &acp.PromptResponse{}
	},
	CancelNotificationType: func() any {
		return &acp.CancelNotification{}
	},
}

// unmarshalMessage decodes a raw websocket message into the acp struct that
// corresponds to the given MessageType.
func UnmarshalMessage(mt MessageType, data []byte) (any, error) {
	constructor, ok := MessageTypeToMessage[mt]
	if !ok {
		return nil, fmt.Errorf("unknown websocket message type: %d", mt)
	}

	target := constructor()
	if err := json.Unmarshal(data, target); err != nil {
		return nil, err
	}

	return target, nil
}
