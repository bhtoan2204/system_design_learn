package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create worker
	w := worker.New(c, "travel-booking-queue", worker.Options{})

	// Register workflow
	w.RegisterWorkflow(TravelBookingOrchestrator)

	// Register activities
	w.RegisterActivity(BookHotelActivity)
	w.RegisterActivity(CancelHotelActivity)
	w.RegisterActivity(BookFlightActivity)
	w.RegisterActivity(CancelFlightActivity)
	w.RegisterActivity(BookCarActivity)
	w.RegisterActivity(CancelCarActivity)

	// Start worker
	log.Println("Starting worker...")
	log.Println("Worker is listening on task queue: travel-booking-queue")
	log.Println("Press Ctrl+C to stop the worker")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
