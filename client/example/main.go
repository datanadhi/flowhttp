package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/datanadhi/flowhttp/client"
)

/*
FlowHTTP Client Example
------------------------

Demonstrates most FlowHTTP Client features:

1. Simple GET and POST requests
2. Query parameters and headers
3. JSON response parsing
4. String and byte helpers
5. Timeout support
6. Status helpers (.IsSuccess, .StatusText)

--------------------------------------
Run the demo:
  go run ./client/examples/basic

--------------------------------------
Features:
--------------------------------------
1. GET with query params and headers
2. POST with JSON payload
3. Automatic Content-Type handling
4. Response helper functions (.Json, .String, .Bytes)
5. Reusable client with timeout control
6. Error and status code handling

--------------------------------------
CURL equivalents:
--------------------------------------
# GET request
curl -v "https://httpbin.org/get?q=flowhttp" -H "X-App: FlowHTTP"

# POST request
curl -v -X POST "https://httpbin.org/post" \
     -H "Content-Type: application/json" \
     -H "X-App: FlowHTTP" \
     -d '{"hello":"world"}'
*/

func main() {
	// Create a reusable client with 5-second timeout
	c := client.NewClient(5 * time.Second)

	// 1. GET request with params and headers
	resp, err := c.Get(
		"https://httpbin.org/get",
		map[string]string{"q": "flowhttp"},
		map[string]string{"X-App": "FlowHTTP"},
	)
	if err != nil {
		panic(fmt.Sprintf("GET failed: %v", err))
	}

	// Parse as JSON
	jsonData, err := resp.Json()
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
	} else {
		fmt.Println("GET JSON Parsed:")
		fmt.Printf("  - URL: %v\n", jsonData["url"])
		fmt.Printf("  - Headers: %v\n", jsonData["headers"])
	}

	fmt.Printf("Status: %d %s\n\n", resp.StatusCode, resp.StatusText())

	// 2. POST request with JSON payload
	payload := strings.NewReader(`{"hello": "world", "from": "FlowHTTP"}`)
	resp2, err := c.Post(
		"https://httpbin.org/post",
		nil,
		map[string]string{"X-App": "FlowHTTP"},
		payload,
		"application/json",
	)
	if err != nil {
		panic(fmt.Sprintf("POST failed: %v", err))
	}

	// Check success
	if resp2.IsSuccess() {
		fmt.Println("POST Success")
	} else {
		fmt.Printf("POST Failed (%d: %s)\n", resp2.StatusCode, resp2.StatusText())
	}

	// Read body as string
	bodyStr, _ := resp2.String()
	fmt.Println("\nPOST Raw Body:\n", bodyStr)

	// Parse JSON again
	postData, _ := resp2.Json()
	fmt.Println("\nPOST JSON (parsed):")
	fmt.Printf("  - Data: %v\n", postData["json"])
	fmt.Printf("  - Headers: %v\n", postData["headers"])

	fmt.Println("\nDone")
}
