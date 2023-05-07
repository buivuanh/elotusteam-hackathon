package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestRegister(t *testing.T) {
	// remove to run this test
	t.Skip()
	userName := "vuanh"
	password := "password123"

	// Create a payload with the username and password
	payload := url.Values{}
	payload.Set("user_name", userName)
	payload.Set("password", password)

	// Encode the payload into a URL-encoded string
	data := []byte(payload.Encode())

	// Create a request to the register API endpoint
	req, err := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	// Set the content type header for URL-encoded form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Register failed. Status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("User registered successfully.")
}

func TestHello(t *testing.T) {
	// remove to run this test
	t.Skip()
	// Create a request to the register API endpoint
	req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	// Set the content type header for URL-encoded form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Register failed. Status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("User registered successfully.")
}

func testLogin(userName, pwd string) string {
	// Create a payload with the username and password
	payload := url.Values{}
	payload.Set("user_name", userName)
	payload.Set("password", pwd)

	// Encode the payload into a URL-encoded string
	data := []byte(payload.Encode())

	// Create a request to the login API endpoint
	req, err := http.NewRequest("POST", "http://localhost:8080/login", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return ""
	}

	// Set the content type header for URL-encoded form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Login failed. Status code: %d\n", resp.StatusCode)
		return ""
	}

	// Parse the response JSON
	var response struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Printf("Failed to parse response JSON: %v\n", err)
		return ""
	}

	return response.Token
}

func TestLogin(t *testing.T) {
	// remove to run this test
	t.Skip()
	userName := "john"
	password := "password123"

	token := testLogin(userName, password)

	fmt.Printf("Login successful. Token: %s\n", token)
}

func TestUploadFile(t *testing.T) {
	// remove to run this test
	t.Skip()
	userName := "vuanh"
	password := "password123"

	token := testLogin(userName, password)

	// Open the file you want to upload
	fileName := "pexels-nout-gons-378570.jpg"
	file, err := os.Open("./test-img/" + fileName)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form field for the file
	fileField, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("Failed to create form field:", err)
		return
	}

	// Copy the file content to the form field
	_, err = io.Copy(fileField, file)
	if err != nil {
		fmt.Println("Failed to copy file content:", err)
		return
	}

	// Close the multipart writer to finalize the body
	writer.Close()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/upload", body)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	// Set the content type header
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Upload failed with status:", resp.Status)
		return
	}

	fmt.Println("File uploaded successfully")
}
