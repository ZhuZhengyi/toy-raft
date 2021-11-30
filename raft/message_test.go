// message_test.go

package raft

import (
	"bytes"
	"testing"
)

func assertEqMsgEvent(event1, event2 MsgEvent) bool {
	return event1.Type() == event2.Type()
}

func (msg1 *Message) equalTo(msg2 *Message) bool {
	return msg1.from == msg2.from &&
		msg1.to == msg2.to &&
		msg1.term == msg2.term &&
		msg1.MsgType() == msg2.MsgType()
}

func TestMessage(t *testing.T) {
	tests := []*Message{
		NewMessage(&AddrLocal{}, &AddrLocal{}, 1, &EventHeartbeatReq{}),
		NewMessage(&AddrLocal{}, &AddrPeer{}, 2, &EventHeartbeatResp{}),
		NewMessage(&AddrPeer{}, &AddrLocal{}, 3, &EventClientReq{}),
		NewMessage(&AddrLocal{}, &AddrPeers{}, 4, &EventClientResp{}),
	}

	logger.SetLogLevel("Debug")

	for _, msg1 := range tests {
		logger.Info("msg size: %v", msg1.Size())
		data := make([]byte, msg1.Size())
		msg1.Marshal(data)
		msg11 := new(Message)
		msg11.Unmarshal(data)
		if msg1.equalTo(msg11) {
			t.Errorf("msg (%v) marshal (%v) ", msg1, msg11)
		}

		buffer := bytes.NewBuffer(data)
		msg1.SendTo(buffer)

		msg12 := new(Message)
		msg12.RecvFrom(buffer)
		if msg1.equalTo(msg12) {
			t.Errorf("msg (%v) marshal (%v) ", msg1, msg11)
		}
	}

}
