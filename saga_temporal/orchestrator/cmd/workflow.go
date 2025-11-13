package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type TravelBookingInput struct {
	BookingID string
}

type TravelBookingResult struct {
	Success      bool
	Message      string
	HotelBooked  bool
	FlightBooked bool
	CarBooked    bool
}

type HotelWorkflowInput struct {
	BookingID string
}

type FlightWorkflowInput struct {
	BookingID string
}

type CarWorkflowInput struct {
	BookingID string
}

func TravelBookingOrchestrator(ctx workflow.Context, input TravelBookingInput) (TravelBookingResult, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	result := TravelBookingResult{}

	hotelInput := HotelWorkflowInput(input)
	err := workflow.ExecuteActivity(ctx, BookHotelActivity, hotelInput).Get(ctx, nil)
	if err != nil {
		return TravelBookingResult{
			Success: false,
			Message: "Failed to book hotel: " + err.Error(),
		}, err
	}
	result.HotelBooked = true

	flightInput := FlightWorkflowInput(input)
	err = workflow.ExecuteActivity(ctx, BookFlightActivity, flightInput).Get(ctx, nil)
	if err != nil {
		workflow.ExecuteActivity(ctx, CancelHotelActivity, hotelInput).Get(ctx, nil)
		return TravelBookingResult{
			Success:     false,
			Message:     "Failed to book flight: " + err.Error() + ". Hotel booking cancelled.",
			HotelBooked: false,
		}, err
	}
	result.FlightBooked = true

	carInput := CarWorkflowInput(input)
	err = workflow.ExecuteActivity(ctx, BookCarActivity, carInput).Get(ctx, nil)
	if err != nil {
		workflow.ExecuteActivity(ctx, CancelFlightActivity, flightInput).Get(ctx, nil)
		workflow.ExecuteActivity(ctx, CancelHotelActivity, hotelInput).Get(ctx, nil)
		return TravelBookingResult{
			Success:      false,
			Message:      "Failed to book car: " + err.Error() + ". Flight and hotel bookings cancelled.",
			HotelBooked:  false,
			FlightBooked: false,
		}, err
	}
	result.CarBooked = true

	result.Success = true
	result.Message = "All bookings completed successfully: hotel, flight, and car"
	return result, nil
}
