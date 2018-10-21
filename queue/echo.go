package queue

import (
	"fmt"
	"os"
)

type EchoQueue struct {
	f *os.File
}

func NewEchoQueue() *EchoQueue {
	return &EchoQueue{
		f: os.Stdout,
	}
}

func (q *EchoQueue) SendPipeline(ev Event) error {
	_, err := fmt.Fprintf(q.f, "%+v", ev)
	return err
}
