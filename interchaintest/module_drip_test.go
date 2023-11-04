package interchaintest

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"

	helpers "github.com/CosmosContracts/furya/tests/interchaintest/helpers"
)

// TestFuryaDrip ensures the drip module properly distributes tokens from whitelisted accounts.
func TestFuryaDrip(t *testing.T) {
	t.Parallel()

	// Setup new pre determined user (from test_node.sh)
	mnemonic := "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry"
	addr := "furya1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl"

	// Base setup
	newCfg := furyaConfig
	newCfg.ModifyGenesis = cosmos.ModifyGenesis(append(defaultGenesisKV, []cosmos.GenesisKV{
		{
			Key:   "app_state.drip.params.allowed_addresses",
			Value: []string{addr},
		},
	}...))

	chains := CreateChainWithCustomConfig(t, 1, 0, newCfg)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// Chains
	furya := chains[0].(*cosmos.CosmosChain)
	nativeDenom := furya.Config().Denom

	// User
	user, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, "default", mnemonic, int64(1000000_000_000), furya)
	if err != nil {
		t.Fatal(err)
	}

	// New TF token to distributes
	tfDenom := helpers.CreateTokenFactoryDenom(t, ctx, furya, user, "dripme", fmt.Sprintf("0%s", Denom))
	distributeAmt := math.NewInt(1_000_000)
	helpers.MintTokenFactoryDenom(t, ctx, furya, user, distributeAmt.Uint64(), tfDenom)

	// Stake some tokens
	vals := helpers.GetValidators(t, ctx, furya)
	valoper := vals.Validators[0].OperatorAddress

	stakeAmt := 100000_000_000
	helpers.StakeTokens(t, ctx, furya, user, valoper, fmt.Sprintf("%d%s", stakeAmt, nativeDenom))

	// Drip the TF Tokens to all stakers
	distribute := int64(1_000_000)
	helpers.DripTokens(t, ctx, furya, user, fmt.Sprintf("%d%s", distribute, tfDenom))

	// Claim staking rewards to capture the drip
	helpers.ClaimStakingRewards(t, ctx, furya, user, valoper)

	// Check balances has the TF Denom from the claim
	bals, _ := furya.AllBalances(ctx, user.FormattedAddress())
	fmt.Println("balances", bals)

	found := false
	for _, bal := range bals {
		if bal.Denom == tfDenom {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("did not find drip token")
	}

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
