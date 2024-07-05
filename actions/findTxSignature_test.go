package actions_test

import (
	"solana-actions/actions"
	"testing"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
)

var account types.Account
var reference actions.Reference
var c *client.Client

func init() {
	account, _ = types.AccountFromBase58("28WJTTqMuurAfz6yqeTrFMXeFd91uzi9i1AW6F5KyHQDS9siXb8TquAuatvLuCEYdggyeiNKLAUr3w7Czmmf2Rav")
	reference = actions.Reference(account.PublicKey)
	c = client.NewClient(rpc.DevnetRPCEndpoint)
}

func TestFindTxSig(t *testing.T) {

	t.Run("findTransactionSignature", func(t *testing.T) {
		t.Run("should return last signature", func(t *testing.T) {
			sig, err := actions.FindReference(c, reference, nil)
			// t.Log(val)
			if err != nil {
				t.Errorf("err should not be nil: %s", err.Error())
			}
			expected := "3DMBLgCQSLdurScFUCPaHy8anCpYkFVaCnAENGvpyeidfUpcDoDjvK5pSVHPkUg8L1qLcnKj6fs7p677ZRiLP37f"
			if sig.Signature != expected {
				t.Logf("expected: %s", expected)
				t.Errorf("got %s want %s", sig.Signature, expected)
				t.Fail()
			}
		})

		t.Run("throws an error on signature not found", func(t *testing.T) {
			account = types.NewAccount()
			reference = actions.Reference(account.PublicKey)
			_, err := actions.FindReference(c, reference, nil)
			if err != nil {
				if err.Error() != "FindReferenceError: not found" {
					t.Logf("expected: not found")
					t.Errorf("got %s want %s", err.Error(), "FindReferenceError: not found")
					t.Fail()
				}
			}
		})
	})
}
