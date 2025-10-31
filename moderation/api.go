package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIClient struct {
	BaseURL   string
	AuthToken string
	client    *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *APIClient) Login(username, password string) error {
	data := map[string]string{"username": username, "password": password}
	body, _ := json.Marshal(data)

	resp, err := c.client.Post(c.BaseURL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed with status %s", resp.Status)
	}

	c.AuthToken = resp.Header.Get("Set-Cookie")
	return nil
}

func (c *APIClient) GetThreads() error {
	req, _ := http.NewRequest("GET", c.BaseURL+"/api/threads", nil)
	if c.AuthToken != "" {
		req.Header.Set("Cookie", c.AuthToken)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch threads: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

func (c *APIClient) DeleteThread(id string) error {
	req, _ := http.NewRequest("DELETE", c.BaseURL+"/api/threads/"+id, nil)
	if c.AuthToken != "" {
		req.Header.Set("Cookie", c.AuthToken)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// TODO MORE ROBUST ERROR HANDLING
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete thread: %s", resp.Status)
	}

	fmt.Println("Thread deleted successfully.")
	return nil
}
