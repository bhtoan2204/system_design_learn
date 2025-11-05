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

func BookCarActivity(ctx context.Context, input CarWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Booking car", "input", input)

	carServiceURL := getEnv("CAR_SERVICE_URL", "http://localhost:8083")

	bookingReq := BookingRequest{
		BookingID: input.BookingID,
	}

	reqBody, err := json.Marshal(bookingReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", carServiceURL+"/book",
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
		return fmt.Errorf("failed to call car service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("car service returned error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var bookingResp BookingResponse
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !bookingResp.Success {
		return fmt.Errorf("car booking failed: %s", bookingResp.Message)
	}

	logger.Info("Car booked successfully")
	return nil
}

func CancelCarActivity(ctx context.Context, input CarWorkflowInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Cancelling car", "input", input)

	carServiceURL := getEnv("CAR_SERVICE_URL", "http://localhost:8083")

	bookingReq := BookingRequest{
		BookingID: input.BookingID,
	}

	reqBody, err := json.Marshal(bookingReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", carServiceURL+"/cancel",
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
		return fmt.Errorf("failed to call car service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("car service returned error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var bookingResp BookingResponse
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !bookingResp.Success {
		return fmt.Errorf("car cancellation failed: %s", bookingResp.Message)
	}

	logger.Info("Car cancelled successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
