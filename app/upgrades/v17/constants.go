package v17

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/CosmosContracts/furya/v18/app/upgrades"
	clocktypes "github.com/CosmosContracts/furya/v18/x/clock/types"
	driptypes "github.com/CosmosContracts/furya/v18/x/drip/types"
)

// UpgradeName defines the on-chain upgrade name for the upgrade.
const UpgradeName = "v17"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV17UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			driptypes.ModuleName,
			clocktypes.ModuleName,
		},
	},
}
