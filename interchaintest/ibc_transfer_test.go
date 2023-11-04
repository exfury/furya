package interchaintest

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

// TestFuryaGaiaIBCTransfer spins up a Furya and Gaia network, initializes an IBC connection between them,
// and sends an ICS20 token transfer from Furya->Gaia and then back from Gaia->Furya.
func TestFuryaGaiaIBCTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Furya and Gaia
	numVals := 1
	numFullNodes := 1

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "furya",
			ChainConfig:   furyaConfig,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "gaia",
			Version:       "v9.1.0",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	const (
		path = "ibc-path"
	)

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	client, network := interchaintest.DockerSetup(t)

	furya, gaia := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	relayerType, relayerName := ibc.CosmosRly, "relay"

	// Get a relayer instance
	rf := interchaintest.NewBuiltinRelayerFactory(
		relayerType,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(IBCRelayerImage, IBCRelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	)

	r := rf.Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(furya).
		AddChain(gaia).
		AddRelayer(r, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  furya,
			Chain2:  gaia,
			Relayer: r,
			Path:    path,
		})

	ctx := context.Background()

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, furya, gaia)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, furya, gaia)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	furyaUser, gaiaUser := users[0], users[1]

	furyaUserAddr := furyaUser.FormattedAddress()
	gaiaUserAddr := gaiaUser.FormattedAddress()

	// Get original account balances
	furyaOrigBal, err := furya.GetBalance(ctx, furyaUserAddr, furya.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, furyaOrigBal.Int64())

	gaiaOrigBal, err := gaia.GetBalance(ctx, gaiaUserAddr, gaia.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, gaiaOrigBal.Int64())

	// Compose an IBC transfer and send from Furya -> Gaia
	var transferAmount = math.NewInt(1_000)
	transfer := ibc.WalletAmount{
		Address: gaiaUserAddr,
		Denom:   furya.Config().Denom,
		Amount:  transferAmount,
	}

	channel, err := ibc.GetTransferChannel(ctx, r, eRep, furya.Config().ChainID, gaia.Config().ChainID)
	require.NoError(t, err)

	furyaHeight, err := furya.Height(ctx)
	require.NoError(t, err)

	transferTx, err := furya.SendIBCTransfer(ctx, channel.ChannelID, furyaUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	err = r.StartRelayer(ctx, eRep, path)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, furya, furyaHeight, furyaHeight+50, transferTx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 10, furya)
	require.NoError(t, err)

	// Get the IBC denom for ufury on Gaia
	furyaTokenDenom := transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, furya.Config().Denom)
	furyaIBCDenom := transfertypes.ParseDenomTrace(furyaTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on Furya and are in the user acc on Gaia
	furyaUpdateBal, err := furya.GetBalance(ctx, furyaUserAddr, furya.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, furyaOrigBal.Sub(transferAmount), furyaUpdateBal)

	gaiaUpdateBal, err := gaia.GetBalance(ctx, gaiaUserAddr, furyaIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, gaiaUpdateBal)

	// Compose an IBC transfer and send from Gaia -> Furya
	transfer = ibc.WalletAmount{
		Address: furyaUserAddr,
		Denom:   furyaIBCDenom,
		Amount:  transferAmount,
	}

	gaiaHeight, err := gaia.Height(ctx)
	require.NoError(t, err)

	transferTx, err = gaia.SendIBCTransfer(ctx, channel.Counterparty.ChannelID, gaiaUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, gaia, gaiaHeight, gaiaHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Assert that the funds are now back on Furya and not on Gaia
	furyaUpdateBal, err = furya.GetBalance(ctx, furyaUserAddr, furya.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, furyaOrigBal, furyaUpdateBal)

	gaiaUpdateBal, err = gaia.GetBalance(ctx, gaiaUserAddr, furyaIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(0), gaiaUpdateBal.Int64())
}
