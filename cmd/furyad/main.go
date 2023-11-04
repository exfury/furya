package main

import (
	"os"

	"cosmossdk.io/log"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/CosmosContracts/furya/v18/app"
	"github.com/CosmosContracts/furya/v18/cmd/furyad/cmd"
)

func main() {
	app.SetAddressPrefixes()
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "FURYAD", app.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
