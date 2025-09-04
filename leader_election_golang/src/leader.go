package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-zookeeper/zk"
)

type LeaderCtrl struct {
	conn       *zk.Conn
	leaderPath string
	data       []byte
	cancelFn   context.CancelFunc
}

func NewLeaderCtrl(ctx context.Context, addrs []string, leaderPath string) (*LeaderCtrl, error) {
	conn, _, err := zk.Connect(addrs, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("zk connect: %w", err)
	}

	parts := splitPath(leaderPath)
	path := ""
	for _, p := range parts[:len(parts)-1] {
		path += "/" + p
		exists, _, err := conn.Exists(path)
		if err != nil {
			conn.Close()
			return nil, err
		}
		if !exists {
			_, err = conn.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
			if err != nil && err != zk.ErrNodeExists {
				conn.Close()
				return nil, err
			}
		}
	}
	hostname, _ := os.Hostname()
	data := []byte(fmt.Sprintf("host=%s pid=%d ts=%d", hostname, os.Getpid(), time.Now().Unix()))

	cctx, cancel := context.WithCancel(ctx)
	go func() {
		<-cctx.Done()
		conn.Close()
	}()

	return &LeaderCtrl{conn: conn, leaderPath: leaderPath, data: data, cancelFn: cancel}, nil
}

func (l *LeaderCtrl) RunElection(ctx context.Context) (becameLeader <-chan struct{}, lostLeader <-chan struct{}) {
	become := make(chan struct{}, 1)
	lost := make(chan struct{}, 1)

	go func() {
		defer close(become)
		defer close(lost)

		for {
			if l.runElectionOnce(ctx, become, lost) {
				return
			}
		}
	}()

	return become, lost
}

func (l *LeaderCtrl) tryCreateLeaderNodeWithErr() (bool, error) {
	_, err := l.conn.Create(l.leaderPath, l.data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *LeaderCtrl) handleTransientCreateError(ctx context.Context, err error) {
	if err == nil || err == zk.ErrNodeExists {
		return
	}
	log.Printf("[leader] create error: %v", err)
	sleepCtx(ctx, 2*time.Second)
}

func (l *LeaderCtrl) waitForLeaderLoss(ctx context.Context, lost chan<- struct{}) bool {
	_, _, ch, err := l.conn.ExistsW(l.leaderPath)
	if err != nil {
		log.Printf("[leader] ExistsW error: %v", err)
		sleepCtx(ctx, 2*time.Second)
		return false
	}
	select {
	case ev := <-ch:
		if ev.Type == zk.EventNodeDeleted || ev.State == zk.StateDisconnected || ev.State == zk.StateExpired {
			notifyNonBlocking(lost)
		}
		return false
	case <-ctx.Done():
		return true
	}
}

func (l *LeaderCtrl) waitForAnyLeaderChange(ctx context.Context) bool {
	_, _, ch, err := l.conn.ExistsW(l.leaderPath)
	if err != nil {
		log.Printf("[leader] ExistsW error: %v", err)
		sleepCtx(ctx, 2*time.Second)
		return false
	}
	select {
	case <-ch:
		return false
	case <-ctx.Done():
		return true
	}
}

func notifyNonBlocking(ch chan<- struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

func (l *LeaderCtrl) runElectionOnce(ctx context.Context, become chan<- struct{}, lost chan<- struct{}) bool {
	acquired, err := l.tryCreateLeaderNodeWithErr()
	if acquired {
		notifyNonBlocking(become)
		return l.waitForLeaderLoss(ctx, lost)
	}

	if err == zk.ErrNodeExists {
		return l.waitForAnyLeaderChange(ctx)
	}

	l.handleTransientCreateError(ctx, err)
	return false
}

func (l *LeaderCtrl) Close() {
	l.cancelFn()
}

func splitPath(p string) []string {
	out := []string{}
	cur := ""
	for _, c := range p {
		if c == '/' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(c)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func sleepCtx(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
	case <-ctx.Done():
	}
}
