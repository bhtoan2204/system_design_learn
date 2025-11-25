package patterns

import (
	"fmt"
	"sync"
	"time"
)

func producer(ch chan int) {
	for i := 1; i <= 10; i++ {
		ch <- i
		time.Sleep(time.Millisecond * 100)
	}
	close(ch)
}

func worker(id int, ch <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range ch {
		fmt.Printf("Worker %d processing %d\n", id, job)
		results <- job * 2
	}
}

func CallerFanOut() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	go producer(jobs)
	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println("Result:", res)
	}
}
