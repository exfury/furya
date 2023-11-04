package v14

import (
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"

	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/exfury/furya/v18/app/upgrades"
	"github.com/exfury/furya/v18/x/globalfee"
)

// UpgradeName defines the on-chain upgrade name for the upgrade.
const UpgradeName = "v14"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV14UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			globalfee.ModuleName,
			ibchookstypes.StoreKey,
		},
	},
}
