package message_queue

import (
	"context"
	"fmt"
	"sync"
)

type PartitionID int

type Message struct {
	ID        string
	Partition PartitionID
	Data      []byte
	Timestamp int64
}

type Consumer interface {
	Start() error
	Stop() error
	Consume(ctx context.Context, partition PartitionID, handler func(Message) error) error
	GetChannels() []chan Message
}

type consumer struct {
	partitions int
	channels   []chan Message
	running    bool
	wg         sync.WaitGroup
	mu         sync.RWMutex
}

func NewConsumer(bufferSize int) Consumer {
	channels := make([]chan Message, 4)
	for i := 0; i < 4; i++ {
		channels[i] = make(chan Message, bufferSize)
	}

	return &consumer{
		partitions: 4,
		channels:   channels,
		running:    false,
	}
}

func (c *consumer) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("consumer is already running")
	}

	c.running = true
	return nil
}

func (c *consumer) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return nil
	}

	c.running = false

	c.wg.Wait()

	for _, ch := range c.channels {
		close(ch)
	}

	return nil
}

func (c *consumer) Consume(ctx context.Context, partition PartitionID, handler func(Message) error) error {
	c.mu.RLock()
	if !c.running {
		c.mu.RUnlock()
		return fmt.Errorf("consumer is not running")
	}
	c.mu.RUnlock()

	partitionIdx := int(partition)
	if partitionIdx < 0 || partitionIdx >= len(c.channels) {
		return fmt.Errorf("invalid partition %d, must be between 0 and %d", partitionIdx, len(c.channels)-1)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case msg, ok := <-c.channels[partitionIdx]:
				if !ok {

					return
				}

				if err := handler(msg); err != nil {

					fmt.Printf("Error processing message %s: %v\n", msg.ID, err)
				}

			case <-ctx.Done():

				return
			}
		}
	}()

	return nil
}

func (c *consumer) GetChannels() []chan Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channels
}
