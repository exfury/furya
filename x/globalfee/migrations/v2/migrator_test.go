package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/exfury/furya/v18/x/globalfee"
	v2 "github.com/exfury/furya/v18/x/globalfee/migrations/v2"
	"github.com/exfury/furya/v18/x/globalfee/types"
)

func TestMigrateMainnet(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(globalfee.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := sdk.NewKVStoreKey(v2.ModuleName)
	tKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	store := ctx.KVStore(storeKey)

	params := types.Params{
		MinimumGasPrices: sdk.DecCoins{
			sdk.NewDecCoinFromDec("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", sdk.NewDecWithPrec(3, 3)),
			sdk.NewDecCoinFromDec("ufury", sdk.NewDecWithPrec(75, 3)),
		},
	}

	require.NoError(t, v2.Migrate(ctx, store, cdc, "ufury"))

	var res types.Params
	bz := store.Get(v2.ParamsKey)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, params, res)
}

func TestMigrateTestnet(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(globalfee.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := sdk.NewKVStoreKey(v2.ModuleName)
	tKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	store := ctx.KVStore(storeKey)

	params := types.Params{
		MinimumGasPrices: sdk.DecCoins{
			sdk.NewDecCoinFromDec("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", sdk.NewDecWithPrec(1, 3)),
			sdk.NewDecCoinFromDec("ufuryx", sdk.NewDecWithPrec(25, 4)),
		},
	}

	require.NoError(t, v2.Migrate(ctx, store, cdc, "ufuryx"))

	var res types.Params
	bz := store.Get(v2.ParamsKey)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, params, res)
}
