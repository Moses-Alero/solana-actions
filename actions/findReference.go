package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
)

type FindReferenceError struct {
	Message string
}

func (e *FindReferenceError) Error() string {
	return fmt.Sprintf("FindReferenceError: %s", e.Message)
}

/*
 * Find the oldest transaction signature referencing a given public key.
 *
 * @param connection - A connection to the cluster.
 * @param reference - `reference` in the Solana Action spec.
 * @param options - Options for `getSignaturesForAddress`.
 *
 * @throws {FindReferenceError}
 */

func FindReference(connection *client.Client, reference Reference, options *client.GetSignaturesForAddressConfig) (*rpc.SignatureWithStatus, error) {
	signatures, err := connection.GetSignaturesForAddress(context.Background(), reference.String())
	if err != nil {
		return nil, err
	}
	sigLen := len(signatures)
	if sigLen < 1 {
		return nil, &FindReferenceError{"not found"}
	}
	oldest := &signatures[sigLen-1]
	var limit int
	if options != nil && options.Limit > 0 {
		limit = options.Limit
	} else {
		limit = 1000
	}
	if sigLen < limit {
		return oldest, nil
	}

	// Recursively find the oldest one in the unlikely event that signatures up to the limit are found.
	opts := &client.GetSignaturesForAddressConfig{
		Limit:      options.Limit,
		Before:     oldest.Signature,
		Until:      options.Until,
		Commitment: options.Commitment,
	}
	oldest, err = FindReference(connection, reference, opts)
	if err != nil {
		if errors.Is(err, &FindReferenceError{}) {
			return oldest, nil
		}
		return nil, err
	}
	return oldest, nil
}
