package patterns

import (
	"fmt"
	"time"
)

func workerPool(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		time.Sleep(time.Second) // Simulate work
		results <- j * 2        // Return result
	}
}

func CallerWorkerPools() {
	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go workerPool(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs)
	for a := 1; a <= numJobs; a++ {
		fmt.Printf("Result: %d\n", <-results)
	}
}
