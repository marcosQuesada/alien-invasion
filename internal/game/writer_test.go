package game

import (
	"testing"
	"time"
)

func TestItWrapsChannelWithLockingBehaviour(t *testing.T) {
	w := NewChannelWriter()
	select {
	case <-w.Chan():
		t.Fatal("expected read locking behaviour")
	default:
	}
	_ = w.Close()
}

func TestItWrapsChannelFullWriterWorkFlow(t *testing.T) {
	w := NewChannelWriter()

	payload := "foo"
	go func() {
		_, _ = w.Write([]byte(payload))
	}()
	select {
	case res, ok := <-w.Chan():
		if !ok {
			t.Fatal("expected open channel")
		}
		if res != payload {
			t.Errorf("result does not match, expected %s got %s", payload, res)
		}
	case <-time.NewTimer(time.Second).C:
		t.Fatal("Read Timeout")
	}
	_ = w.Close()
}
