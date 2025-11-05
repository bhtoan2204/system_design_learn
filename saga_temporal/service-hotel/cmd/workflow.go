package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type HotelWorkflowInput struct {
	BookingID string
}

type HotelWorkflowResult struct {
	Success bool
	Message string
}

func HotelWorkflow(ctx workflow.Context, input HotelWorkflowInput) (HotelWorkflowResult, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, BookHotelActivity, input).Get(ctx, nil)
	if err != nil {
		return HotelWorkflowResult{
			Success: false,
			Message: "Failed to book hotel",
		}, err
	}

	return HotelWorkflowResult{
		Success: true,
		Message: "Hotel booked successfully",
	}, nil
}
