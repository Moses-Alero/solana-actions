package main

import (
	"fmt"
	"net/url"
)

func main() {
	fmt.Println("This is my implementation of solana actions ink")
	urlString := "https://example.com/path/to/resource"

	// Parse the URL string into a *url.URL struct
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Get the protocol (scheme) of the URL
	ParseURL(parsedURL)
}
