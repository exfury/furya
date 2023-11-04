package v12

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/exfury/furya/v18/app/upgrades"
)

const UpgradeName = "v12"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV12UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
