package main

import (
	"fmt"
	"net/url"
	"solana-actions/actions"
)

func main() {
	fmt.Println("This is my implementation of solana actions ink")
	urlString := "solana-action:https://actions.alice.com/donate?label=wonderland&message=the_builder"
	// "https://example.domain?action=solana-action%3Ahttps%3A%2F%2Factions.alice.com%2Fdonate"

	// Parse the URL string into a *url.URL struct
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Get the protocol (scheme) of the URL
	val, err := actions.ParseURL(parsedURL)
	if err != nil {
		fmt.Println(err)
	}

	action, ok := val.(*actions.ActionRequestURLFields)
	if ok {
		//do action url stuff

		url, err := actions.EncodeUrl(action, actions.SOLANA_ACTIONS_PROTOCOL)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(url.String())
	}

	blink, ok := val.(*actions.BlinkURLFields)
	if ok {
		//do blinks url stuff
		fmt.Println(blink.Action.Link)
		_, err := actions.EncodeUrl(blink, actions.SOLANA_ACTIONS_PROTOCOL)
		if err != nil {
			fmt.Println(err)
		}
	}
}
