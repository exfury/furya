package interchaintest

import (
	"fmt"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/stretchr/testify/require"

	helpers "github.com/exfury/furya/tests/interchaintest/helpers"
)

// TestFuryaCwHooks
// from x/cw-hooks/keeper/msg_server_test.go -> TestContractExecution
func TestFuryaCwHooks(t *testing.T) {
	t.Parallel()

	cfg := furyaConfig

	// Base setup
	chains := CreateChainWithCustomConfig(t, 1, 0, cfg)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// Chains
	furya := chains[0].(*cosmos.CosmosChain)

	// Users
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", int64(10_000_000_000), furya, furya)
	user := users[0]

	// Upload & init contract payment to another address
	_, contractAddr := helpers.SetupContract(t, ctx, furya, user.KeyName(), "contracts/furya_staking_hooks_example.wasm", `{}`)

	// register staking contract (to be tested)
	helpers.RegisterCwHooksStaking(t, ctx, furya, user, contractAddr)
	sc := helpers.GetCwHooksStakingContracts(t, ctx, furya)
	require.Equal(t, 1, len(sc))
	require.Equal(t, contractAddr, sc[0])

	// Validate that governance contract is added
	helpers.RegisterCwHooksGovernance(t, ctx, furya, user, contractAddr)
	gc := helpers.GetCwHooksGovernanceContracts(t, ctx, furya)
	require.Equal(t, 1, len(gc))
	require.Equal(t, contractAddr, gc[0])

	// Perform a Staking Action
	vals := helpers.GetValidators(t, ctx, furya)
	valoper := vals.Validators[0].OperatorAddress

	stakeAmt := 1_000_000
	helpers.StakeTokens(t, ctx, furya, user, valoper, fmt.Sprintf("%d%s", stakeAmt, furya.Config().Denom))

	// Query the smart contract to validate it saw the fire-and-forget update
	res := helpers.GetCwStakingHookLastDelegationChange(t, ctx, furya, contractAddr, user.FormattedAddress())
	require.Equal(t, valoper, res.Data.ValidatorAddress)
	require.Equal(t, user.FormattedAddress(), res.Data.DelegatorAddress)
	require.Equal(t, fmt.Sprintf("%d.000000000000000000", stakeAmt), res.Data.Shares)

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
