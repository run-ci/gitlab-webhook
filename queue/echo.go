package queue

import (
	"fmt"
	"os"

	"gitlab.com/run-ci/webhooks/pkg.git"
)

type EchoQueue struct {
	f *os.File
}

func NewEchoQueue() *EchoQueue {
	return &EchoQueue{
		f: os.Stdout,
	}
}

func (q *EchoQueue) SendPipeline(p pkg.Pipeline) error {
	_, err := fmt.Fprintf(q.f, "%+v", p)
	return err
}
