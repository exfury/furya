package interchaintest

import (
	"context"
	"fmt"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/docker/docker/client"
	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	ibclocalhost "github.com/cosmos/ibc-go/v7/modules/light-clients/09-localhost"

	clocktypes "github.com/CosmosContracts/furya/v18/x/clock/types"
	feepaytypes "github.com/CosmosContracts/furya/v18/x/feepay/types"
	feesharetypes "github.com/CosmosContracts/furya/v18/x/feeshare/types"
	globalfeetypes "github.com/CosmosContracts/furya/v18/x/globalfee/types"
	tokenfactorytypes "github.com/CosmosContracts/furya/v18/x/tokenfactory/types"
)

var (
	VotingPeriod     = "15s"
	MaxDepositPeriod = "10s"
	Denom            = "ufury"

	FuryaE2ERepo  = "ghcr.io/cosmoscontracts/furya-e2e"
	FuryaMainRepo = "ghcr.io/cosmoscontracts/furya"

	IBCRelayerImage   = "ghcr.io/cosmos/relayer"
	IBCRelayerVersion = "main"

	furyaRepo, furyaVersion = GetDockerImageInfo()

	FuryaImage = ibc.DockerImage{
		Repository: furyaRepo,
		Version:    furyaVersion,
		UidGid:     "1025:1025",
	}

	// SDK v47 Genesis
	defaultGenesisKV = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: VotingPeriod,
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: MaxDepositPeriod,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: Denom,
		},
		{
			Key:   "app_state.feepay.params.enable_feepay",
			Value: false,
		},
	}

	furyaConfig = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "furya",
		ChainID:             "furya-2",
		Images:              []ibc.DockerImage{FuryaImage},
		Bin:                 "furyad",
		Bech32Prefix:        "furya",
		Denom:               Denom,
		CoinType:            "118",
		GasPrices:           fmt.Sprintf("0%s", Denom),
		GasAdjustment:       2.0,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		ConfigFileOverrides: nil,
		EncodingConfig:      furyaEncoding(),
		ModifyGenesis:       cosmos.ModifyGenesis(defaultGenesisKV),
	}

	genesisWalletAmount = int64(10_000_000)
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount("furya", "furya")
	sdk.GetConfig().SetBech32PrefixForValidator("furyavaloper", "furya")
	sdk.GetConfig().SetBech32PrefixForConsensusNode("furyavalcons", "furya")
	sdk.GetConfig().SetCoinType(118)
}

// furyaEncoding registers the Furya specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func furyaEncoding() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	ibclocalhost.RegisterInterfaces(cfg.InterfaceRegistry)
	wasmtypes.RegisterInterfaces(cfg.InterfaceRegistry)
	feesharetypes.RegisterInterfaces(cfg.InterfaceRegistry)
	tokenfactorytypes.RegisterInterfaces(cfg.InterfaceRegistry)
	feepaytypes.RegisterInterfaces(cfg.InterfaceRegistry)
	globalfeetypes.RegisterInterfaces(cfg.InterfaceRegistry)
	clocktypes.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}

// CreateChain generates a new chain with a custom image (useful for upgrades)
func CreateChain(t *testing.T, numVals, numFull int, img ibc.DockerImage) []ibc.Chain {
	cfg := furyaConfig
	cfg.Images = []ibc.DockerImage{img}
	return CreateChainWithCustomConfig(t, numVals, numFull, cfg)
}

// CreateThisBranchChain generates this branch's chain (ex: from the commit)
func CreateThisBranchChain(t *testing.T, numVals, numFull int) []ibc.Chain {
	return CreateChain(t, numVals, numFull, FuryaImage)
}

func CreateChainWithCustomConfig(t *testing.T, numVals, numFull int, config ibc.ChainConfig) []ibc.Chain {
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "furya",
			ChainName:     "furya",
			Version:       config.Images[0].Version,
			ChainConfig:   config,
			NumValidators: &numVals,
			NumFullNodes:  &numFull,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	// chain := chains[0].(*cosmos.CosmosChain)
	return chains
}

func BuildInitialChain(t *testing.T, chains []ibc.Chain) (*interchaintest.Interchain, context.Context, *client.Client, string) {
	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain()

	for _, chain := range chains {
		ic = ic.AddChain(chain)
	}

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()
	client, network := interchaintest.DockerSetup(t)

	err := ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
		// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
	})
	require.NoError(t, err)

	return ic, ctx, client, network
}
