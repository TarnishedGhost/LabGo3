package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver receives the texture that was prepared as a result of executing commands in the event loop.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop implements the event loop to generate a texture obtained by executing operations received from an internal queue.
type Loop struct {
	Receiver Receiver

	next screen.Texture // The texture currently being formed
	prev screen.Texture // The texture that was last sent to the Receiver

	Mq MessageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

// Start launches the event loop. This method must be called before any other methods are called on it.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.stop = make(chan struct{})
	go func() {
		for !l.stopReq || !l.Mq.empty() {
			op := l.Mq.pull()
			update := op.Do(l.next)
			if update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
		close(l.stop)
	}()
}

// Post adds a new operation to the internal queue.
func (l *Loop) Post(op Operation) {
	if update := op.Do(l.next); update {
		l.Receiver.Update(l.next)
		l.next, l.prev = l.prev, l.next
	}
}

// StopAndWait signals the need to stop the loop and blocks until it's completely stopped.
func (l *Loop) StopAndWait() {
}

type MessageQueue struct {
	Ops     []Operation
	mu      sync.Mutex
	blocked chan struct{}
}

func (mq *MessageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.Ops = append(mq.Ops, op)
	if mq.blocked != nil {
		close(mq.blocked)
		mq.blocked = nil
	}
}

func (mq *MessageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	for len(mq.Ops) == 0 {
		mq.blocked = make(chan struct{})
		mq.mu.Unlock()
		<-mq.blocked
		mq.mu.Lock()
	}
	op := mq.Ops[0]
	mq.Ops[0] = nil
	mq.Ops = mq.Ops[1:]
	return op
}

func (mq *MessageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.Ops) == 0
}
