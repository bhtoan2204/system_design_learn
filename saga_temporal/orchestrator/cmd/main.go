package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "travel-booking-orchestrator",
		TaskQueue: "travel-booking-queue",
	}
	workflowInput := TravelBookingInput{
		BookingID: "travel-booking-123",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, TravelBookingOrchestrator, workflowInput)
	if err != nil {
		log.Fatalln("unable to execute Workflow", err)
	}

	var result TravelBookingResult
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("unable to get workflow result", err)
	}

	printResults(workflowInput, result, we.GetID(), we.GetRunID())
}

func printResults(input TravelBookingInput, result TravelBookingResult, workflowID, runID string) {
	log.Printf(
		"\nWorkflowInput: %+v\n",
		input,
	)
	log.Printf(
		"\nWorkflowResult: %+v\n",
		result,
	)
	log.Printf(
		"\nWorkflowID: %s RunID: %s\n",
		workflowID,
		runID,
	)
}
