package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func logRequest(method, url string, payload map[string]string) {
	fmt.Printf("[HTTP Request] Method: %s, URL: %s, Payload: %v\n", method, url, payload)
}

func logResponse(body []byte) {
	fmt.Printf("[HTTP Response] Body: %s\n", string(body))
}

func (c *Client) RegisterMember(username, password string) {
	payload := map[string]string{
		"Username": username,
		"Password": password,
	}
	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/register", c.baseURL)
	logRequest("POST", url, payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error registering member: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logResponse(body)
}

func (c *Client) CreateCommunity(name, description string) {
	payload := map[string]string{
		"Name":        name,
		"Description": description,
	}
	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/community", c.baseURL)
	logRequest("POST", url, payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error creating community: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logResponse(body)
}

func (c *Client) CreateThread(title, content, creatorID, communityID string) string {
	payload := map[string]string{
		"Title":       title,
		"Content":     content,
		"CreatorID":   creatorID,
		"CommunityID": communityID,
	}
	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/thread", c.baseURL)
	logRequest("POST", url, payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error creating thread: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logResponse(body)

	// Extract the thread ID from the response
	var threadID string
	fmt.Sscanf(string(body), "Thread created with ID: %s", &threadID)
	return threadID
}

func (c *Client) CreateReply(content, creatorID, threadID, parentID string) {
	payload := map[string]string{
		"Content":   content,
		"CreatorID": creatorID,
		"ThreadID":  threadID,
		"ParentID":  parentID,
	}
	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/reply", c.baseURL)
	logRequest("POST", url, payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error creating reply: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logResponse(body)
}

func mainClient() {
	client := NewClient("http://localhost:8080")

	client.RegisterMember("test_user", "password123")
	client.CreateCommunity("test_community", "A test community description.")

	// Create a thread and dynamically capture its ID
	threadID := client.CreateThread("Welcome Thread", "Welcome to the community!", "test_user", "test_community")
	if threadID == "" {
		fmt.Println("Failed to create thread. Exiting...")
		return
	}

	// Use the captured thread ID to create a reply
	client.CreateReply("Thanks for the welcome!", "test_user", threadID, "")
}
