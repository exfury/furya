package interchaintest

import (
	"fmt"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/stretchr/testify/require"

	helpers "github.com/CosmosContracts/furya/tests/interchaintest/helpers"
)

// TestFuryaFeePay
func TestFuryaFeePay(t *testing.T) {
	t.Parallel()

	cfg := furyaConfig
	cfg.GasPrices = "0.0025ufury"

	// 0.002500000000000000
	coin := sdk.NewDecCoinFromDec(cfg.Denom, sdk.NewDecWithPrec(25, 4))
	cfg.ModifyGenesis = cosmos.ModifyGenesis(append(defaultGenesisKV, []cosmos.GenesisKV{
		{
			Key:   "app_state.globalfee.params.minimum_gas_prices",
			Value: sdk.DecCoins{coin},
		},
		{
			// override default impl.
			Key:   "app_state.feepay.params.enable_feepay",
			Value: true,
		},
	}...))

	// Base setup
	chains := CreateChainWithCustomConfig(t, 1, 0, cfg)
	ic, ctx, _, _ := BuildInitialChain(t, chains)

	// Chains
	furya := chains[0].(*cosmos.CosmosChain)

	nativeDenom := furya.Config().Denom

	// Users
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", int64(10_000_000), furya, furya)
	admin := users[0]
	user := users[1]

	// Upload & init contract payment to another address
	codeId, err := furya.StoreContract(ctx, admin.KeyName(), "contracts/cw_template.wasm", "--fees", "50000ufury")
	if err != nil {
		t.Fatal(err)
	}

	contractAddr, err := furya.InstantiateContract(ctx, admin.KeyName(), codeId, `{"count":0}`, true)
	if err != nil {
		t.Fatal(err)
	}

	// Register contract for 0 fee usage (x amount of times)
	limit := 5
	balance := 1_000_000
	helpers.RegisterFeePay(t, ctx, furya, admin, contractAddr, limit)
	helpers.FundFeePayContract(t, ctx, furya, admin, contractAddr, strconv.Itoa(balance)+nativeDenom)

	beforeContract := helpers.GetFeePayContract(t, ctx, furya, contractAddr)
	t.Log("beforeContract", beforeContract)
	require.Equal(t, beforeContract.FeePayContract.Balance, strconv.Itoa(balance))
	require.Equal(t, beforeContract.FeePayContract.WalletLimit, strconv.Itoa(int(limit)))

	// execute against it from another account with enough fees (standard Tx)
	txHash, err := furya.ExecuteContract(ctx, user.KeyName(), contractAddr, `{"increment":{}}`, "--fees", "500"+nativeDenom)
	require.NoError(t, err)
	fmt.Println("txHash", txHash)

	beforeBal, err := furya.GetBalance(ctx, user.FormattedAddress(), nativeDenom)
	require.NoError(t, err)

	// execute against it from another account and have the dev pay it
	txHash, err = furya.ExecuteContract(ctx, user.KeyName(), contractAddr, `{"increment":{}}`, "--fees", "0"+nativeDenom)
	require.NoError(t, err)
	fmt.Println("txHash", txHash)

	afterBal, err := furya.GetBalance(ctx, user.FormattedAddress(), nativeDenom)
	require.NoError(t, err)

	// validate users balance did not change
	require.Equal(t, beforeBal, afterBal)

	// validate the contract balance went down
	afterContract := helpers.GetFeePayContract(t, ctx, furya, contractAddr)
	t.Log("afterContract", afterContract)
	require.Equal(t, afterContract.FeePayContract.Balance, strconv.Itoa(balance-500))

	uses := helpers.GetFeePayUses(t, ctx, furya, contractAddr, user.FormattedAddress())
	t.Log("uses", uses)
	require.Equal(t, uses.Uses, "1")

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
