package main

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
)

// --- Activities ---

func BookHotelActivity(ctx context.Context, input HotelWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info(fmt.Sprintf("[BookHotelActivity] Booking hotel with ID: %s", input.BookingID))
	time.Sleep(500 * time.Millisecond) // simulate processing
	logger.Info(fmt.Sprintf("[BookHotelActivity] Hotel booked successfully: %s", input.BookingID))
	return nil
}

func CancelHotelActivity(ctx context.Context, input HotelWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info(fmt.Sprintf("[CancelHotelActivity] Cancelling hotel booking: %s", input.BookingID))
	time.Sleep(500 * time.Millisecond)
	logger.Info(fmt.Sprintf("[CancelHotelActivity] Hotel cancelled: %s", input.BookingID))
	return nil
}

func BookFlightActivity(ctx context.Context, input FlightWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info(fmt.Sprintf("[BookFlightActivity] Booking flight with ID: %s", input.BookingID))
	time.Sleep(500 * time.Millisecond)
	logger.Info(fmt.Sprintf("[BookFlightActivity] Flight booked successfully: %s", input.BookingID))
	return nil
}

func CancelFlightActivity(ctx context.Context, input FlightWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info(fmt.Sprintf("[CancelFlightActivity] Cancelling flight booking: %s", input.BookingID))
	time.Sleep(500 * time.Millisecond)
	logger.Info(fmt.Sprintf("[CancelFlightActivity] Flight cancelled: %s", input.BookingID))
	return nil
}

func BookCarActivity(ctx context.Context, input CarWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info(fmt.Sprintf("[BookCarActivity] Booking car with ID: %s", input.BookingID))
	time.Sleep(500 * time.Millisecond)
	logger.Info(fmt.Sprintf("[BookCarActivity] Car booked successfully: %s", input.BookingID))
	return nil
}

func CancelCarActivity(ctx context.Context, input CarWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info(fmt.Sprintf("[CancelCarActivity] Cancelling car booking: %s", input.BookingID))
	time.Sleep(500 * time.Millisecond)
	logger.Info(fmt.Sprintf("[CancelCarActivity] Car cancelled: %s", input.BookingID))
	return nil
}
