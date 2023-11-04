package simulation_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/exfury/furya/v18/x/mint/simulation"
	"github.com/exfury/furya/v18/x/mint/types"
)

// TestRandomizedGenState tests the normal scenario of applying RandomizedGenState.
// Abonormal scenarios are not tested here.
func TestRandomizedGenState(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	s := rand.NewSource(1)
	r := rand.New(s)

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: math.NewInt(1_000),
		GenState:     make(map[string]json.RawMessage),
	}

	simulation.RandomizedGenState(&simState)

	var mintGenesis types.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &mintGenesis)

	require.Equal(t, uint64(6311520), mintGenesis.Params.BlocksPerYear)
	require.Equal(t, "stake", mintGenesis.Params.MintDenom)
	require.Equal(t, "0stake", mintGenesis.Minter.BlockProvision(mintGenesis.Params, sdk.NewInt(0)).String())
	require.Equal(t, "0.170000000000000000", mintGenesis.Minter.NextAnnualProvisions(mintGenesis.Params, sdk.OneInt()).String())
	require.Equal(t, "0.400000000000000000", mintGenesis.Minter.PhaseInflationRate(1).String())
	require.Equal(t, "0.170000000000000000", mintGenesis.Minter.Inflation.String())
	require.Equal(t, uint64(1), mintGenesis.Minter.NextPhase(mintGenesis.Params, sdk.NewInt(1)))
	require.Equal(t, uint64(0), mintGenesis.Minter.Phase)
	require.Equal(t, "0.000000000000000000", mintGenesis.Minter.AnnualProvisions.String())
}

// TestRandomizedGenState tests abnormal scenarios of applying RandomizedGenState.
func TestRandomizedGenState1(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	s := rand.NewSource(1)
	r := rand.New(s)
	// all these tests will panic
	tests := []struct {
		simState module.SimulationState
		panicMsg string
	}{
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{}, "invalid memory address or nil pointer dereference"},
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{
				AppParams: make(simtypes.AppParams),
				Cdc:       cdc,
				Rand:      r,
			}, "assignment to entry in nil map"},
	}

	for _, tt := range tests {
		tt := tt
		require.Panicsf(t, func() { simulation.RandomizedGenState(&tt.simState) }, tt.panicMsg)
	}
}
