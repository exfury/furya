package v11

import (
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"

	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/exfury/furya/v18/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the Furya v11 upgrade.
const UpgradeName = "v11" // maybe multiverse?

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV11UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{icacontrollertypes.StoreKey, icahosttypes.StoreKey},
	},
}
