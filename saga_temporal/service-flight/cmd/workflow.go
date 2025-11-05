package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type FlightWorkflowInput struct {
	BookingID string
}

type FlightWorkflowResult struct {
	Success bool
	Message string
}

func FlightWorkflow(ctx workflow.Context, input FlightWorkflowInput) (FlightWorkflowResult, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, BookFlightActivity, input).Get(ctx, nil)
	if err != nil {
		return FlightWorkflowResult{
			Success: false,
			Message: "Failed to book flight",
		}, err
	}

	return FlightWorkflowResult{
		Success: true,
		Message: "Flight booked successfully",
	}, nil
}
