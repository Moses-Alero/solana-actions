package main

import (
	"fmt"
	"net/url"
	"solana-actions/types"
)

type URL url.URL

// Thrown when a URL cant be parsed as a Solana Action URL
type ParseUrlError struct {
	Message string
}

func (err ParseUrlError) Error() string {
	return fmt.Sprintf("ParseURLError: %s", err.Message)
}

func ParseURL(url *url.URL) {
	fmt.Println(url.Scheme)

	parseActionRequestURL(url)

}

func parseActionRequestURL(url *url.URL) (*types.ActionRequestURLFields, error) {
	fmt.Println(url.Path)
	path := url.Path
	queryParams := url.Query()

	link, err := url.Parse(path)
	if err != nil {
		return nil, &ParseUrlError{"invalid url"}
	}
	if link.Scheme != types.HTTPS_PROTOCOL {
		return nil, &ParseUrlError{"invalid link"}
	}
	label := queryParams.Get("label")
	message := queryParams.Get("message")

	actionUrlFields := &types.ActionRequestURLFields{
		Link:    *link,
		Label:   &label,
		Message: &message,
	}

	return actionUrlFields, nil
}

func parseBlinksURL(blink url.URL) (*types.BlinkURLFields, error) {
	blinkQuery := blink.Query()
	link := blinkQuery.Get(types.BLINKS_QUERY_PARAM)
	if link == "" {
		return nil, &ParseUrlError{"invalid blink url"}
	}
	_, err := url.Parse(link)
	if err != nil {
		return nil, &ParseUrlError{"error parsing url"}
	}
	blinkUrlFields := &types.BlinkURLFields{
		Blink: blink,
	}
	//todo: complete this
	return blinkUrlFields, nil
}
