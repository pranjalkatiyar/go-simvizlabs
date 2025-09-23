package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func FetchTransaction(jwtToken, transactionID string) ([]byte, error) {
	baseURL := os.Getenv("BASE_URL")
	url := fmt.Sprintf("%s/inApps/v1/transactions/%s", baseURL, transactionID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func FetchTransactionHistory(jwtToken, transactionID string) ([]byte, error) {
	baseURL := os.Getenv("BASE_URL")

	url := fmt.Sprintf("%s/inApps/v2/history/%s", baseURL, transactionID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func FetchAllSubscriptionStatuses(jwtToken, originalTransactionId string) ([]byte, error) {
	baseURL := os.Getenv("BASE_URL")
	url := fmt.Sprintf("%s/inApps/v1/subscriptions/%s", baseURL, originalTransactionId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("apple API error: %s", string(body))
	}
	return body, nil
}
