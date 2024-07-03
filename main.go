package main

import (
	"fmt"

	"net/url"
)

func main() {
	fmt.Println("This is my implementation of solana actions ink")
	urlString := "solana-action:https://actions.alice.com/donate?alice=wonderland&bob=the_builder"
	// "https://example.domain?action=solana-action%3Ahttps%3A%2F%2Factions.alice.com%2Fdonate"

	// Parse the URL string into a *url.URL struct
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Get the protocol (scheme) of the URL
	val, err := ParseURL(parsedURL)
	if err != nil {
		fmt.Println(err)
	}

	action, ok := val.(*ActionRequestURLFields)
	if ok {
		//do action url stuff

		url, err := EncodeUrl(action, SOLANA_ACTIONS_PROTOCOL)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(url.String())
	}

	blink, ok := val.(*BlinkURLFields)
	if ok {
		//do blinks url stuff
		fmt.Println(blink.Action.Link)
		_, err := EncodeUrl(blink, SOLANA_ACTIONS_PROTOCOL)
		if err != nil {
			fmt.Println(err)
		}
	}

}
