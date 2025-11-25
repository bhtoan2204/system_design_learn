package message_queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Producer interface {
	Produce(ctx context.Context, message Message) error
	Close() error
	GetChannels() []chan Message
}

type producer struct {
	channels []chan Message
	mu       sync.RWMutex
	running  bool
}

func NewProducer(numPartitions int, bufferSize int) Producer {
	channels := make([]chan Message, numPartitions)
	for i := 0; i < numPartitions; i++ {
		channels[i] = make(chan Message, bufferSize)
	}

	return &producer{
		channels: channels,
		running:  true,
	}
}

func (p *producer) Produce(ctx context.Context, message Message) error {
	p.mu.RLock()
	if !p.running {
		p.mu.RUnlock()
		return fmt.Errorf("producer is closed")
	}
	p.mu.RUnlock()

	partition := int(message.Partition)
	if partition < 0 || partition >= len(p.channels) {
		return fmt.Errorf("invalid partition %d, must be between 0 and %d", partition, len(p.channels)-1)
	}

	select {
	case p.channels[partition] <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout: failed to produce message to partition %d", partition)
	}
}

func (p *producer) GetChannels() []chan Message {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.channels
}

func (p *producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return nil
	}

	p.running = false

	for _, ch := range p.channels {
		close(ch)
	}

	return nil
}
