package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/twystd/tweetnacl-go/tweetnacl"
)

type ActionPostResponseWithSerializedTransaction struct {
	ActionPostResponse
	Transaction types.Transaction
}

type FetchActionError struct {
	Message string
}

func (e *FetchActionError) Error() string {
	return fmt.Sprintf("SerializeTransactionError: %s", e.Message)
}

/*
Fetch the action payload from a Solana Action request link.

@param connection - A connection to the cluster.

@param link - `link` in the Solana Action spec.

@param fields - Action Post Request Fields

@param options - Options for `getRecentBlockhash`.

@throws {FetchActionError}
*/
func FetchTransaction(conn *client.Client, link *url.URL, fields ActionPostRequest, commitment rpc.Commitment) (*ActionPostResponseWithSerializedTransaction, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(fields)
	if err != nil {
		return nil, &FetchActionError{err.Error()}
	}
	req, err := http.NewRequest(http.MethodPost, link.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &FetchActionError{err.Error()}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, &FetchActionError{err.Error()}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &FetchActionError{err.Error()}
	}

	var actionResp ActionPostResponse
	json.Unmarshal(body, &actionResp)

	if actionResp.Transaction == "" {
		return nil, &FetchActionError{"missing transaction"}
	}
	account := common.PublicKeyFromString(fields.Account)
	tx, err := SerializeTransaction(conn, account, actionResp.Transaction, commitment)
	if err != nil {
		return nil, &FetchActionError{err.Error()}
	}

	actionPostResp := new(ActionPostResponseWithSerializedTransaction)
	actionPostResp.ActionPostResponse = actionResp
	actionPostResp.Transaction = *tx

	return actionPostResp, nil
}

/*
Thrown when the base64 encoded action `transaction` cannot be serialized
*/
type SerializeTransactionError struct {
	Message string
}

func (e *SerializeTransactionError) Error() string {
	return fmt.Sprintf("SerializeTransactionError: %s", e.Message)
}

/*
Serialize a base64 encoded transaction into a web3.js `Transaction`.

	@param connection - A connection to the cluster.

	@param account - Account that may sign the transaction.

	@param base64Transaction - `transaction` in the Solana Action spec.

	@param commitment - transaction Commitment.

@throws {SerializeTransactionError}
*/
func SerializeTransaction(conn *client.Client, account common.PublicKey, base64Tx string, commitment rpc.Commitment) (*types.Transaction, error) {

	tx, err := types.TransactionDeserialize([]byte(base64Tx))
	if err != nil {
		return nil, &SerializeTransactionError{err.Error()}
	}
	sigs := tx.Signatures
	feePayer := tx.Message.Accounts[0]
	recentBlockHash := tx.Message.RecentBlockHash

	if len(sigs) > 0 {
		if feePayer.String() == "" {
			return nil, &SerializeTransactionError{"missing fee payer"}
		}
		if feePayer.String() != string(sigs[0]) {
			return nil, &SerializeTransactionError{"invalid fee payer"}
		}
		if recentBlockHash == "" {
			return nil, &SerializeTransactionError{"recent block hash missing"}
		}

		// A valid signature for everything except `account` must be provided.
		msg, err := tx.Serialize()
		if err != nil {
			return nil, &SerializeTransactionError{err.Error()}
		}

		for _, s := range sigs {
			if s != nil {
				ok, err := tweetnacl.CryptoOneTimeAuthVerify(s, msg, feePayer.Bytes())
				if err != nil {
					return nil, &SerializeTransactionError{err.Error()}
				}
				if !ok {
					return nil, &SerializeTransactionError{"invalid signature"}
				} else if feePayer == account {
					// If the only signature expected is for `account`, ignore the recent blockhash in the transaction.
					if len(sigs) == 1 {
						recentBlkHash, err := conn.GetLatestBlockhashWithConfig(context.Background(), client.GetLatestBlockhashConfig{
							Commitment: commitment,
						})
						if err != nil {
							return nil, &SerializeTransactionError{err.Error()}
						}
						tx.Message.RecentBlockHash = recentBlkHash.Blockhash
					}
				}
			} else {
				return nil, &SerializeTransactionError{"missing signature"}
			}
		}
	} else {
		tx.Message.Accounts[0] = account
		recentBlkHash, err := conn.GetLatestBlockhashWithConfig(context.Background(), client.GetLatestBlockhashConfig{
			Commitment: commitment,
		})
		if err != nil {
			return nil, &SerializeTransactionError{err.Error()}
		}
		tx.Message.RecentBlockHash = recentBlkHash.Blockhash
	}
	return &tx, nil
}
