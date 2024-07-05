package actions

import (
	"net/url"

	"github.com/blocto/solana-go-sdk/common"
)

// @internal
type SupportedProtocol string

const (
	SOLANA_PAY_PROTOCOL SupportedProtocol = "solana"

	SOLANA_ACTIONS_PROTOCOL SupportedProtocol = "solana-action"

	SOLANA_ACTIONS_PROTOCOL_PLURAL SupportedProtocol = "solana-actions"

	HTTPS_PROTOCOL = "https"

	//Program Id for the SPL Memo program
	MEMO_PROGRAM_ID = "MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr"

	BLINKS_QUERY_PARAM = "action"
)

var SupportedProtocols = []SupportedProtocol{
	SOLANA_ACTIONS_PROTOCOL,
	SOLANA_PAY_PROTOCOL,
	SOLANA_ACTIONS_PROTOCOL_PLURAL,
}

func (sp SupportedProtocol) IsValidProtocol() bool {
	for _, protocol := range SupportedProtocols {
		if sp == protocol {
			return true
		}
	}
	return false
}

/*
Standard headers
*/
var ACTIONS_CORS_HEADERS = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "GET, POST, PUT, OPTIONS",
	"Access-Control-Allow-Headers": "Content-Type, Authorization, Content-Encoding, Accept-Encoding",
	"Content-Type":                 "application/json",
}

// `reference` in the [Solana Actions spec](https://github.com/solana-labs/solana-pay/blob/master/SPEC.md#reference)
type Reference common.PublicKey

func (ref Reference) String() string {
	return common.PublicKey(ref).String()

}

// `memo` in the [Solana Actions spec](https://github.com/solana-labs/solana-pay/blob/master/SPEC.md#memo)
type Memo string

type ActionsJson struct {
	Rules []ActionRuleObject
}

type ActionRuleObject struct {
	//relative (preferred) or absolute path to perform the rule mapping from
	PathPattern string

	//relative (preferred) or absolute path that supports Action requests
	ApiPath string
}

/*
Fields of a Solana Action transaction request URL.
*/
type ActionRequestURLFields struct {
	//`link` in the Solana Action spec
	Link *url.URL

	//`label` in the Solana Action spec
	Label *string `json:"label,omitempty"`

	//`message` in the Solana Action spec
	Message *string `json:"message,omitempty"`
}

/**
 * Fields of a blink URL to support a Solana Action.
 */
type BlinkURLFields struct {
	//base URL for the `blink` in the Solana Action spec
	Blink *url.URL

	//`action` passed via the blink `action` query param
	Action ActionRequestURLFields
}

/**
 * # Reserved for future use
 *
 * Response body payload sent via the Action GET Request
 */
type ActionGetRequest struct{}

/**
 * Response body payload returned from the Action GET Request
 */
type ActionGetResponse struct {
	//image url that represents the source of the action request
	Icon string

	//describes the source of the action request
	Title string

	//brief summary of the action to be performed
	Description string

	//button text rendered to the user
	Label string

	//UI state for the button being rendered to the user
	Disabled *bool

	Links *struct {
		//list of related Actions a user could perform
		Actions []LinkedAction
	}

	//non-fatal error message to be displayed to the user
	Error *ActionError
}

/**
 * Related action on a single endpoint
 */
type LinkedAction struct {
	//URL endpoint for an action
	Href string

	//button text rendered to the user
	Label string

	//parameters used to accept user input within an action
	Parameters *[]ActionParameter `json:"parameters,omitempty"`
}

// Parameter to accept user input within an action
type ActionParameter struct {
	// Parameter name in URL
	Name string `json:"name"`
	// Placeholder text for the user input field (optional)
	Label *string `json:"label,omitempty"`
	// Declare if this field is required (defaults to false) (optional)
	Required *bool `json:"required,omitempty"`
}

// Response body payload sent via the Action POST Request
type ActionPostRequest struct {
	// Base58-encoded public key of an account that may sign the transaction
	Account string `json:"account"`
}

// Response body payload returned from the Action POST Request
type ActionPostResponse struct {
	// Base64 encoded serialized transaction
	Transaction string `json:"transaction"`
	// Describes the nature of the transaction (optional)
	Message *string `json:"message,omitempty"`
}

// Non-fatal error message to be displayed to the user
type ActionError struct {
	// Non-fatal error message to be displayed to the user
	Message string `json:"message"`
}
