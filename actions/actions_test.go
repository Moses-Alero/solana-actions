package actions_test

import (
	"fmt"
	"net/url"
	"solana-actions/actions"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	t.Run("ActionRequestURL", func(t *testing.T) {
		t.Run("encodes Url without action params", func(t *testing.T) {
			rawUrl := "https://example.com/api/action"
			link, _ := url.Parse(rawUrl)

			actionURLField := &actions.ActionRequestURLFields{
				Link: link,
			}

			URL, err := actions.EncodeUrl(actionURLField, actions.SOLANA_ACTIONS_PROTOCOL)
			if err != nil {
				t.Log("err should be nil", err)
				t.Fail()
			}
			if URL.String() != "solana-action:"+rawUrl {
				t.Errorf("got %s want %s", URL.String(), "solana-action:"+rawUrl)
				t.Fail()
			}
		})

		t.Run("encodes a URL with additional action params", func(t *testing.T) {
			rawUrl := "https://example.com/api/action"
			link, _ := url.Parse(rawUrl)
			label := "label"
			message := "message"

			actionURLField := &actions.ActionRequestURLFields{
				Link:    link,
				Label:   &label,
				Message: &message,
			}

			URL, err := actions.EncodeUrl(actionURLField, actions.SOLANA_ACTIONS_PROTOCOL)
			if err != nil {
				t.Log("err should be nil", err)
				t.Fail()
			}

			expected := fmt.Sprintf("solana-action:%s?label=%s&message=%s", rawUrl, label, message)
			if URL.String() != expected {
				t.Errorf("got %s want %s", URL.String(), expected)
				t.Fail()
			}

		})

		t.Run("encodes a URL with query parameters", func(t *testing.T) {
			rawUrl := "https://example.com/api/action?query=param"
			link, _ := url.Parse(rawUrl)

			actionURLField := &actions.ActionRequestURLFields{
				Link: link,
			}
			URL, err := actions.EncodeUrl(actionURLField, actions.SOLANA_ACTIONS_PROTOCOL)
			if err != nil {
				t.Log("err should be nil", err)
				t.Fail()
			}

			expected := fmt.Sprintf("solana-action:%s", url.QueryEscape(rawUrl))
			if URL.String() != expected {
				t.Logf("expected: %s", expected)
				t.Errorf("got %s want %s", URL.String(), expected)
				t.Fail()
			}
		})

		t.Run("encodes a URL with query parameters AND action params", func(t *testing.T) {
			rawUrl := "https://example.com/api/action?query=param&amount=1337"
			link, _ := url.Parse(rawUrl)
			label := "label"
			message := "message"

			actionURLField := &actions.ActionRequestURLFields{
				Link:    link,
				Label:   &label,
				Message: &message,
			}

			URL, err := actions.EncodeUrl(actionURLField, actions.SOLANA_ACTIONS_PROTOCOL)
			if err != nil {
				t.Log("err should be nil", err)
				t.Fail()
			}
			expected := fmt.Sprintf("solana-action:%s?label=%s&message=%s", url.QueryEscape(rawUrl), label, message)
			if URL.String() != expected {
				t.Logf("expected: %s", expected)
				t.Errorf("got %s want %s", URL.String(), expected)
				t.Fail()
			}
		})
	})

	t.Run("BlinkUrl", func(t *testing.T) {
		t.Run("encodes a URL without action params", func(t *testing.T) {
			rawBlinkUrl := "https://blink.com/"
			rawUrl := "https://action.com/api/action"

			blink, _ := url.Parse(rawBlinkUrl)
			link, _ := url.Parse(rawUrl)

			blinkUrlField := &actions.BlinkURLFields{
				Blink: blink,
				Action: actions.ActionRequestURLFields{
					Link: link,
				},
			}

			URL, err := actions.EncodeUrl(blinkUrlField, actions.SOLANA_ACTIONS_PROTOCOL)
			if err != nil {
				t.Log("err should be nil", err)
				t.Fail()
			}

			expected := fmt.Sprintf("%s", url.QueryEscape("solana-action:"+rawUrl))
			if URL.Query().Get("action") != expected {
				t.Logf("expected: %s", expected)
				t.Errorf("got %s want %s", URL.Query().Get("action"), expected)
				t.Fail()
			}
		})

		t.Run("encodes a URL with action params", func(t *testing.T) {
			rawBlinkUrl := "https://blink.com/"
			rawUrl := "https://action.com/api/action?query=param"

			blink, _ := url.Parse(rawBlinkUrl)
			link, _ := url.Parse(rawUrl)

			blinkUrlField := &actions.BlinkURLFields{
				Blink: blink,
				Action: actions.ActionRequestURLFields{
					Link: link,
				},
			}
			URL, err := actions.EncodeUrl(blinkUrlField, actions.SOLANA_ACTIONS_PROTOCOL)
			if err != nil {
				t.Log("err should be nil", err)
				t.Fail()
			}

			expected := url.QueryEscape("solana-action:" + fmt.Sprintf("%s", url.QueryEscape(rawUrl)))
			if URL.Query().Get("action") != expected {
				t.Logf("expected: %s", expected)
				t.Errorf("got %s want %s", URL.Query().Get("action"), expected)
				t.Fail()
			}
		})

		t.Run("encodes a URL with query params AND without action params", func(t *testing.T) {
			rawBlinkUrl := "https://blink.com/?other=one"
			rawUrl := "https://action.com/api/action"

			blink, _ := url.Parse(rawBlinkUrl)
			link, _ := url.Parse(rawUrl)

			blinkUrlField := &actions.BlinkURLFields{
				Blink: blink,
				Action: actions.ActionRequestURLFields{
					Link: link,
				},
			}

			URL, err := actions.EncodeUrl(blinkUrlField, actions.SOLANA_ACTIONS_PROTOCOL)
			if err != nil {
				t.Log("err should be nil", err)
				t.Fail()
			}

			expected := fmt.Sprintf("%s", url.QueryEscape("solana-action:"+rawUrl))
			if URL.Query().Get("action") != expected {
				t.Logf("expected: %s", expected)
				t.Errorf("got %s want %s", URL.Query().Get("action"), expected)
				t.Fail()
			}
		})

		t.Run("encodes a URL with query params AND with action params", func(t *testing.T) {
			rawBlinkUrl := "https://blink.com/?other=one"
			rawUrl := "https://action.com/api/action?query=param"

			blink, _ := url.Parse(rawBlinkUrl)
			link, _ := url.Parse(rawUrl)

			blinkUrlField := &actions.BlinkURLFields{
				Blink: blink,
				Action: actions.ActionRequestURLFields{
					Link: link,
				},
			}

			URL, err := actions.EncodeUrl(blinkUrlField, actions.SOLANA_ACTIONS_PROTOCOL)
			if err != nil {
				t.Log("err should be nil", err)
				t.Fail()
			}

			expected := url.QueryEscape("solana-action:" + fmt.Sprintf("%s", url.QueryEscape(rawUrl)))
			if URL.Query().Get("action") != expected {
				t.Logf("expected: %s", expected)
				t.Errorf("got %s want %s", URL.Query().Get("action"), expected)
				t.Fail()
			}
		})
	})
}
