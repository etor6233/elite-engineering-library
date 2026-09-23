package whatsappbridge

// AUTHORED projection guard for one exact human-approved Meta text request.
import (
	"bytes"
	"elite.local/enterprise/internal/channels"
	"encoding/json"
	"io"
	"time"
	"unicode/utf8"
)

type ReplyRequest struct {
	Kind            string `json:"kind"`
	Recipient       string `json:"recipient"`
	Text            string `json:"text"`
	SourceMessageID string `json:"source_message_id"`
	LastInboundAt   int64  `json:"last_inbound_at"`
	WindowExpiresAt int64  `json:"window_expires_at"`
}

func replyRequestFor(message channels.Message) (json.RawMessage, error) {
	var body ReplyRequest
	d := json.NewDecoder(bytes.NewBufferString(message.Text))
	d.DisallowUnknownFields()
	if d.Decode(&body) != nil || d.Decode(new(any)) != io.EOF || body.Kind != "text_reply" || body.Recipient != message.ExternalID || !utf8.ValidString(body.Text) || len(body.Text) == 0 || len(body.Text) > 4096 || utf8.RuneCountInString(body.Text) > 1024 || body.SourceMessageID == "" || len(body.SourceMessageID) > 256 || body.LastInboundAt < 1 || body.WindowExpiresAt != body.LastInboundAt+86400 {
		return nil, ErrBridge
	}
	for _, r := range body.Text {
		if r < 32 && r != '\n' && r != '\t' {
			return nil, ErrBridge
		}
	}
	now := time.Now().Unix()
	if now < body.LastInboundAt || now >= body.WindowExpiresAt {
		return nil, ErrBridge
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, ErrBridge
	}
	return raw, nil
}
