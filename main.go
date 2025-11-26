package main

import (
	"encoding/json"  // For parsing JSON from the API
	"fmt"            // For printing to console
	"io"             // For reading API response
	"net/http"       // For making HTTP requests
	"sync"           // For waiting on goroutines
)

// Joke struct to hold the API response (setup and punchline)
type Joke struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func main() {
	// Channel to receive jokes from goroutines (buffered for 3 jokes)
	jokeChan := make(chan string, 3)
	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Launch 3 goroutines
	for i := 0; i < 3; i++ {
		wg.Add(1)  // Increment wait counter
		go fetchJoke(&wg, jokeChan)  // Run in background
	}

	// Wait for all goroutines to finish
	wg.Wait()
	// Close channel after all are done
	close(jokeChan)

	// Print the jokes
	fmt.Println("Fetched Jokes:")
	for joke := range jokeChan {
		fmt.Println(joke)
	}
}

// Function to fetch a single joke (runs in goroutine)
func fetchJoke(wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()  // Decrement wait counter when done

	// Make HTTP GET request
	resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
	if err != nil {
		ch <- fmt.Sprintf("Error fetching joke: %v", err)  // Send error to channel
		return
	}
	defer resp.Body.Close()  // Close response body

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("Error reading response: %v", err)
		return
	}

	// Parse JSON into Joke struct
	var joke Joke
	err = json.Unmarshal(body, &joke)
	if err != nil {
		ch <- fmt.Sprintf("Error parsing JSON: %v", err)
		return
	}

	// Send formatted joke to channel
	ch <- fmt.Sprintf("%s\n%s", joke.Setup, joke.Punchline)
}