package messaging

import (
	"errors"
	"github.com/appellative-ai/common/core"
)

func NewConfigMessage(v any) *Message {
	return NewMessage(ChannelControl, ConfigEvent).SetContent(ContentTypeAny, v)
}

func ConfigContent[T any](m *Message) (t T, ok bool) {
	if m == nil || m.Content == nil || m.ContentType() != ContentTypeAny {
		return
	}
	t, ok = m.Content.Value.(T)
	return
}

func UpdateContent[T any](m *Message, t *T) bool {
	if m == nil || m.Content == nil || m.ContentType() != ContentTypeAny {
		return false
	}
	if t1, ok := m.Content.Value.(T); ok {
		*t = t1
		return true
	}
	return false
}

func NewStatusMessage(statusCode int, relatesTo string) *Message {
	m := NewMessage(ChannelControl, StatusEvent).SetContent(ContentTypeStatus, statusCode)
	if relatesTo != "" {
		m.SetRelatesTo(relatesTo)
	}
	return m
}

func StatusContent(m *Message) (int, string, error) {
	if !ValidContent(m, StatusEvent, ContentTypeStatus) {
		return 0, "", errors.New("invalid content")
	}
	t, status := core.New[int](m.Content)
	if status == nil {
		return t, m.RelatesTo(), status
	}
	return 0, "", status
}
