package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.temporal.io/sdk/activity"
)

type BookingRequest struct {
	BookingID string `json:"bookingId"`
}

type BookingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func BookHotelActivity(ctx context.Context, input HotelWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Booking hotel", "input", input)

	hotelServiceURL := getEnv("HOTEL_SERVICE_URL", "http://localhost:8081")

	bookingReq := BookingRequest{
		BookingID: activity.GetInfo(ctx).WorkflowExecution.ID,
	}

	reqBody, err := json.Marshal(bookingReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", hotelServiceURL+"/book",
		bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call hotel service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("hotel service returned error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var bookingResp BookingResponse
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !bookingResp.Success {
		return fmt.Errorf("hotel booking failed: %s", bookingResp.Message)
	}

	logger.Info("Hotel booked successfully")
	return nil
}

func CancelHotelActivity(ctx context.Context, input HotelWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Cancelling hotel", "input", input)

	hotelServiceURL := getEnv("HOTEL_SERVICE_URL", "http://localhost:8081")

	bookingReq := BookingRequest{
		BookingID: activity.GetInfo(ctx).WorkflowExecution.ID,
	}

	reqBody, err := json.Marshal(bookingReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", hotelServiceURL+"/cancel",
		bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call hotel service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("hotel service returned error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var bookingResp BookingResponse
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !bookingResp.Success {
		return fmt.Errorf("hotel cancellation failed: %s", bookingResp.Message)
	}

	logger.Info("Hotel cancelled successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
