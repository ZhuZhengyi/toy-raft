package raft

import "sync"

var (
	msgPool = &sync.Pool{
		New: func() interface{} {
			return &Message{}
		},
	}

	bytePool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 4096)
		},
	}
)

func takeMessage() *Message {
	msg := msgPool.Get().(*Message)

	return msg
}

func putMessage(msg *Message) {
	if msg != nil {
		msgPool.Put(msg)
	}
}

func takeBytes() []byte {
	return bytePool.Get().([]byte)
}

func putBytes(b []byte) {
	bytePool.Put(b)
}
