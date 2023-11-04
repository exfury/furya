package app

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

const (
	// DefaultFuryaInstanceCost is initially set the same as in wasmd
	DefaultFuryaInstanceCost uint64 = 60_000
	// DefaultFuryaCompileCost set to a large number for testing
	DefaultFuryaCompileCost uint64 = 3
)

// FuryaGasRegisterConfig is defaults plus a custom compile amount
func FuryaGasRegisterConfig() wasmtypes.WasmGasRegisterConfig {
	gasConfig := wasmtypes.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultFuryaInstanceCost
	gasConfig.CompileCost = DefaultFuryaCompileCost

	return gasConfig
}

func NewFuryaWasmGasRegister() wasmtypes.WasmGasRegister {
	return wasmtypes.NewWasmGasRegister(FuryaGasRegisterConfig())
}
