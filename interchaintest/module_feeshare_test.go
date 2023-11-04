package interchaintest

import (
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"

	helpers "github.com/CosmosContracts/furya/tests/interchaintest/helpers"
)

// TestFuryaFeeShare ensures the feeshare module register and execute sharing functions work properly on smart contracts.
func TestFuryaFeeShare(t *testing.T) {
	t.Parallel()

	// Base setup
	chains := CreateThisBranchChain(t, 1, 0)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// Chains
	furya := chains[0].(*cosmos.CosmosChain)

	nativeDenom := furya.Config().Denom

	// Users
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", int64(10_000_000), furya, furya)
	user := users[0]
	feeRcvAddr := "furya1v75wlkccpv7le3560zw32v2zjes5n0e7csr4qh"

	// Upload & init contract payment to another address
	_, contractAddr := helpers.SetupContract(t, ctx, furya, user.KeyName(), "contracts/cw_template.wasm", `{"count":0}`)

	// register contract to a random address (since we are the creator, though not the admin)
	helpers.RegisterFeeShare(t, ctx, furya, user, contractAddr, feeRcvAddr)
	if balance, err := furya.GetBalance(ctx, feeRcvAddr, nativeDenom); err != nil {
		t.Fatal(err)
	} else if balance.Int64() != 0 {
		t.Fatal("balance not 0")
	}

	// execute with a 10000 fee (so 5000 denom should be in the contract now with 50% feeshare default)
	helpers.ExecuteMsgWithFee(t, ctx, furya, user, contractAddr, "", "10000"+nativeDenom, `{"increment":{}}`)

	// check balance of nativeDenom now
	if balance, err := furya.GetBalance(ctx, feeRcvAddr, nativeDenom); err != nil {
		t.Fatal(err)
	} else if balance.Int64() != 5000 {
		t.Fatal("balance not 5,000. it is ", balance, nativeDenom)
	}

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
