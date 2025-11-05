package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type CarWorkflowInput struct {
	BookingID string
}

type CarWorkflowResult struct {
	Success bool
	Message string
}

func CarWorkflow(ctx workflow.Context, input CarWorkflowInput) (CarWorkflowResult, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, BookCarActivity, input).Get(ctx, nil)
	if err != nil {
		return CarWorkflowResult{
			Success: false,
			Message: "Failed to book car",
		}, err
	}

	return CarWorkflowResult{
		Success: true,
		Message: "Car booked successfully",
	}, nil
}
