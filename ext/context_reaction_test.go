package ext

import (
	"context"
	"errors"
	"testing"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

type reactionInvoker struct {
	updates tg.UpdatesClass
	err     error
}

func (i reactionInvoker) Invoke(_ context.Context, _ bin.Encoder, output bin.Decoder) error {
	output.(*tg.UpdatesBox).Updates = i.updates
	return i.err
}

func TestSendReaction(t *testing.T) {
	rpcErr := errors.New("send failed")
	tests := []struct {
		name    string
		updates tg.UpdatesClass
		err     error
		wantID  int
	}{
		{"reaction update", &tg.Updates{Updates: []tg.UpdateClass{&tg.UpdateMessageReactions{}}}, nil, 0},
		{"new message", &tg.Updates{Updates: []tg.UpdateClass{&tg.UpdateNewMessage{Message: &tg.Message{ID: 42}}}}, nil, 42},
		{"RPC error", nil, rpcErr, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{Context: context.Background(), Raw: tg.NewClient(reactionInvoker{tt.updates, tt.err}), setReply: true}
			msg, err := ctx.SendReaction(0, &tg.MessagesSendReactionRequest{Peer: &tg.InputPeerSelf{}})
			if !errors.Is(err, tt.err) {
				t.Fatalf("error = %v, want %v", err, tt.err)
			}
			if tt.wantID == 0 && msg != nil {
				t.Fatalf("message = %v, want nil", msg)
			}
			if tt.wantID != 0 && (msg == nil || msg.ID != tt.wantID) {
				t.Fatalf("message = %v, want ID %d", msg, tt.wantID)
			}
		})
	}
}
