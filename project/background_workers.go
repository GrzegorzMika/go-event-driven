package main

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

type Message struct {
	Task     Task
	TicketID string
}

type Callback func(context.Context, string) error

type Worker struct {
	queue     chan Message
	callbacks map[Task]Callback
}

func NewWorker() *Worker {
	return &Worker{
		queue:     make(chan Message, 100),
		callbacks: make(map[Task]Callback),
	}
}

func (w *Worker) AddCallback(task Task, callback Callback) {
	w.callbacks[task] = callback
}

func (w *Worker) Send(msg ...Message) {
	for _, m := range msg {
		w.queue <- m
	}
}

func (w *Worker) Run(ctx context.Context) {

	for {
		select {
		case msgp := <-w.queue:
			err := w.callbacks[msgp.Task](ctx, msgp.TicketID)
			if err != nil {
				logrus.WithError(err).Error("requeuing task")
				w.Send(Message{
					Task:     msgp.Task,
					TicketID: msgp.TicketID,
				})
			}
		case <-ctx.Done():
			return
		}
	}
}
