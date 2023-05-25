package v2

import (
	"github.com/CosmosContracts/juno/v15/x/feeshare/exported"
	"github.com/CosmosContracts/juno/v15/x/feeshare/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "feeshare"
)

var (
	// feeshare/types/keys.go -> prefixParams
	ParamsKey = []byte{0x04}
)

// Migrate migrates the x/feeshare module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/feeshare
// module state.
func Migrate(
	ctx sdk.Context,
	store sdk.KVStore,
	legacySubspace exported.Subspace,
	cdc codec.BinaryCodec,
) error {
	var currParams types.Params
	legacySubspace.GetParamSet(ctx, &currParams)

	if err := currParams.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&currParams)
	store.Set(ParamsKey, bz)

	return nil
}
