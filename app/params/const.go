package params

const (
	// Name defines the application name of the Furya network.
	Name = "furya"

	// BondDenom defines the native staking token denomination.
	BondDenom = "ufury"

	// DisplayDenom defines the name, symbol, and display value of the Furya token.
	DisplayDenom = "FURY"

	// DefaultGasLimit - set to the same value as cosmos-sdk flags.DefaultGasLimit
	// this value is currently only used in tests.
	DefaultGasLimit = 200000
)
