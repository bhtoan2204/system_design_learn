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

// Input types for activities
type HotelWorkflowInput struct {
	BookingID string
}

type FlightWorkflowInput struct {
	BookingID string
}

type CarWorkflowInput struct {
	BookingID string
}

// TravelBookingOrchestrator orchestrates the saga pattern:
// 1. Book hotel
// 2. Book flight (if hotel succeeds)
// 3. Book car (if flight succeeds)
// Compensation:
// - If car fails -> cancel flight
// - If flight fails -> cancel hotel
func TravelBookingOrchestrator(ctx workflow.Context, input TravelBookingInput) (TravelBookingResult, error) {
	// Set workflow options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	result := TravelBookingResult{}

	// Step 1: Book hotel
	hotelInput := HotelWorkflowInput{
		BookingID: input.BookingID,
	}
	err := workflow.ExecuteActivity(ctx, BookHotelActivity, hotelInput).Get(ctx, nil)
	if err != nil {
		return TravelBookingResult{
			Success: false,
			Message: "Failed to book hotel: " + err.Error(),
		}, err
	}
	result.HotelBooked = true

	// Step 2: Book flight
	flightInput := FlightWorkflowInput{
		BookingID: input.BookingID,
	}
	err = workflow.ExecuteActivity(ctx, BookFlightActivity, flightInput).Get(ctx, nil)
	if err != nil {
		// Compensation: Cancel hotel because flight failed
		workflow.ExecuteActivity(ctx, CancelHotelActivity, hotelInput).Get(ctx, nil)
		return TravelBookingResult{
			Success:     false,
			Message:     "Failed to book flight: " + err.Error() + ". Hotel booking cancelled.",
			HotelBooked: false,
		}, err
	}
	result.FlightBooked = true

	// Step 3: Book car
	carInput := CarWorkflowInput{
		BookingID: input.BookingID,
	}
	err = workflow.ExecuteActivity(ctx, BookCarActivity, carInput).Get(ctx, nil)
	if err != nil {
		// Compensation: Cancel flight because car failed
		workflow.ExecuteActivity(ctx, CancelFlightActivity, flightInput).Get(ctx, nil)
		// Compensation: Cancel hotel because car failed
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
