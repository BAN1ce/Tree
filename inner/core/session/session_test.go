package session

import (
	"context"
	"sync"
	"testing"
)

func TestSession_PutKey(t *testing.T) {
	type fields struct {
		Node *Node
		mux  sync.RWMutex
	}
	type args struct {
		ctx   context.Context
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Node: tt.fields.Node,
				mux:  tt.fields.mux,
			}
			if err := s.PutKey(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("PutKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
