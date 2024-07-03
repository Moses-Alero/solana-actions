package main

import (
	"fmt"
	"net/url"
	"regexp"
)

// Thrown when a URL can't be parsed as a Solana Action URL
type ParseUrlError struct {
	Message string
}

func (err ParseUrlError) Error() string {
	return fmt.Sprintf("ParseURLError: %s", err.Message)
}

func ParseURL(url *url.URL) (any, error) {
	match, _ := regexp.MatchString(`^https?`, url.Scheme)
	if match {
		return parseBlinksURL(url)
	}
	if url.Scheme != string(SOLANA_PAY_PROTOCOL) &&
		url.Scheme != string(SOLANA_ACTIONS_PROTOCOL) &&
		url.Scheme != string(SOLANA_ACTIONS_PROTOCOL_PLURAL) {
		return nil, &ParseUrlError{"protocol invalid"}
	}
	if url.Opaque == "" {
		return nil, &ParseUrlError{"pathname missing"}
	}
	match, _ = regexp.MatchString(`[:%]`, url.Opaque)
	if !match {
		return nil, &ParseUrlError{"pathname invalid"}
	}
	return parseActionRequestURL(url)
}

func parseActionRequestURL(url *url.URL) (*ActionRequestURLFields, error) {
	opaque := url.Opaque
	queryParams := url.Query()

	link, err := url.Parse(opaque)
	if err != nil {
		return nil, &ParseUrlError{"invalid url"}
	}
	if link.Scheme != HTTPS_PROTOCOL {
		return nil, &ParseUrlError{"invalid link"}
	}
	label := queryParams.Get("label")
	message := queryParams.Get("message")

	actionUrlFields := &ActionRequestURLFields{
		Link:    link,
		Label:   &label,
		Message: &message,
	}

	return actionUrlFields, nil
}

func parseBlinksURL(blink *url.URL) (*BlinkURLFields, error) {
	blinkQuery := blink.Query()
	link := blinkQuery.Get(BLINKS_QUERY_PARAM)
	if link == "" {
		return nil, &ParseUrlError{"invalid blink url"}
	}
	linkUrl, err := url.Parse(link)
	if err != nil {
		return nil, &ParseUrlError{"error parsing url"}
	}
	parsedUrl, err := ParseURL(linkUrl)

	if err != nil {
		return nil, &ParseUrlError{err.Error()}
	}

	action, ok := parsedUrl.(*ActionRequestURLFields)

	if !ok {
		return nil, ParseUrlError{"invalid action type"}
	}

	blinkUrlFields := &BlinkURLFields{
		Blink:  blink,
		Action: *action,
	}
	return blinkUrlFields, nil
}
