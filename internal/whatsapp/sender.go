package whatsapp

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendText sends a plain text message to a JID.
func SendText(wa *whatsmeow.Client, jid types.JID, text string) error {
	msg := &waE2E.Message{
		Conversation: proto.String(text),
	}
	_, err := wa.SendMessage(context.Background(), jid, msg)
	return err
}

// SendReply sends a text message as a reply to a specific message.
func SendReply(wa *whatsmeow.Client, jid types.JID, text string, quoted *waE2E.Message, quotedID string, quotedSender string) error {
	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String(quotedID),
				Participant:   proto.String(quotedSender),
				QuotedMessage: quoted,
			},
		},
	}
	_, err := wa.SendMessage(context.Background(), jid, msg)
	return err
}

// SendTyping sends a composing presence indicator.
func SendTyping(wa *whatsmeow.Client, jid types.JID, delayMs int) {
	ctx := context.Background()
	if err := wa.SendPresence(ctx, types.PresenceAvailable); err != nil {
		log.Debug().Err(err).Msg("failed to send presence")
	}
	if err := wa.SendChatPresence(ctx, jid, types.ChatPresenceComposing, types.ChatPresenceMediaText); err != nil {
		log.Debug().Err(err).Msg("failed to send chat presence")
	}
	if delayMs > 0 {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}

// StopTyping clears the composing presence indicator.
func StopTyping(wa *whatsmeow.Client, jid types.JID) {
	_ = wa.SendChatPresence(context.Background(), jid, types.ChatPresencePaused, types.ChatPresenceMediaText)
}
