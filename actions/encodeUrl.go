package actions

import (
	"fmt"
	"net/url"
	"strings"
)

//QUESTION: Can we make union types with generics in golang??

type EncodedUrlError struct {
	Message string
}

func (e *EncodedUrlError) Error() string {
	return fmt.Sprintf("EncodeURLError: %s", e.Message)
}

func EncodeUrl(fields any, protocol SupportedProtocol) (*url.URL, error) {
	//check the types
	actionUrlFields, isActionUrlField := fields.(*ActionRequestURLFields)
	if isActionUrlField {
		return encodeActionRequestUrl(actionUrlFields, protocol)
	}
	blinkUrlFields, isBlinkField := fields.(*BlinkURLFields)
	if isBlinkField {
		return encodeBlinkUrl(blinkUrlFields, protocol)
	}

	return nil, &EncodedUrlError{"invalid field type, must be of type *ActionRequestURLFields or *BlinkUrlFields"}

}

func encodeActionRequestUrl(fields *ActionRequestURLFields, protocol SupportedProtocol) (*url.URL, error) {
	if !protocol.IsValidProtocol() {
		return nil, &EncodedUrlError{"Invalid protocol"}
	}
	pathname := fields.Link.RawQuery
	if pathname != "" {
		modifiedURL := strings.Replace(fields.Link.String(), `\?`, "?", 1)
		pathname = url.QueryEscape(modifiedURL)
	} else {
		pathname = strings.Replace(fields.Link.String(), `/$`, "", 1)
	}
	link := fmt.Sprintf("%s:%s", protocol, pathname)
	URL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	if fields.Label != nil {
		queryParams := URL.Query()
		queryParams.Set("label", *fields.Label)
		URL.RawQuery = queryParams.Encode()
	}
	if fields.Message != nil {
		queryParams := URL.Query()
		queryParams.Set("message", *fields.Message)
		URL.RawQuery = queryParams.Encode()
	}
	return URL, nil
}

func encodeBlinkUrl(fields *BlinkURLFields, protocol SupportedProtocol) (*url.URL, error) {
	if !protocol.IsValidProtocol() {
		return nil, &EncodedUrlError{"Invalid protocol"}
	}
	URL, err := url.Parse(fields.Blink.String())
	if err != nil {
		return nil, err
	}
	actionUrl, err := encodeActionRequestUrl(&fields.Action, protocol)

	encodedUri := url.QueryEscape(actionUrl.String())
	queryParams := URL.Query()
	queryParams.Set(BLINKS_QUERY_PARAM, encodedUri)
	URL.RawQuery = queryParams.Encode()
	return URL, nil
}
