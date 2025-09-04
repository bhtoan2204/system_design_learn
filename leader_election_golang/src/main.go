package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TaskLeaderPing = "leader:ping"
)

func main() {
	redisAddr := getenv("REDIS_ADDR", "redis:6379")
	redisPass := os.Getenv("REDIS_PASSWORD")
	zkAddrs := getenv("ZK_ADDRS", "zookeeper:2181")
	leaderPath := getenv("ZK_LEADER_PATH", "/asynq/leader")

	r := asynq.RedisClientOpt{Addr: redisAddr, Password: redisPass}

	srv := asynq.NewServer(r, asynq.Config{
		Concurrency: 1,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskLeaderPing, func(ctx context.Context, t *asynq.Task) error {
		log.Println("I'm leader")
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zkList := splitCSV(zkAddrs)
	leader, err := NewLeaderCtrl(ctx, zkList, leaderPath)
	if err != nil {
		log.Fatalf("leader ctrl error: %v", err)
	}
	defer leader.Close()

	becameLeader, lostLeader := leader.RunElection(ctx)

	var (
		sched      *asynq.Scheduler
		schedDoneC chan struct{}
	)

	startScheduler := func() {
		if sched != nil {
			return
		}
		sched = asynq.NewScheduler(r, &asynq.SchedulerOpts{})
		if _, err := sched.Register("@every 10s", asynq.NewTask(TaskLeaderPing, nil)); err != nil {
			log.Fatalf("register schedule failed: %v", err)
		}
		schedDoneC = make(chan struct{})
		go func() {
			log.Println("[scheduler] starting (leader acquired)")
			if err := sched.Run(); err != nil {
				log.Printf("[scheduler] stopped with error: %v", err)
			} else {
				log.Printf("[scheduler] stopped")
			}
			close(schedDoneC)
		}()
	}

	stopScheduler := func() {
		if sched == nil {
			return
		}
		log.Println("[scheduler] shutting down (lost leadership)")
		sched.Shutdown()
		// chờ goroutine Run() kết thúc (tránh rò rỉ)
		select {
		case <-schedDoneC:
		case <-time.After(5 * time.Second):
		}
		sched = nil
		schedDoneC = nil
	}

	errCh := make(chan error, 1)
	go func() {
		log.Println("[server] starting...")
		if err := srv.Run(mux); err != nil {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-becameLeader:
			startScheduler()
		case <-lostLeader:
			stopScheduler()

		case sig := <-sigCh:
			log.Printf("received signal: %v, shutting down...", sig)
			break loop

		case err := <-errCh:
			log.Printf("server runtime error: %v", err)
			break loop
		}
	}

	stopScheduler()
	srv.Shutdown()
	log.Println("bye!")
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func splitCSV(s string) []string {
	out := []string{}
	cur := ""
	for _, c := range s {
		if c == ',' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
		} else {
			cur += string(c)
		}
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
