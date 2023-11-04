package types

import (
	"testing"

	"github.com/stretchr/testify/suite"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodecTestSuite struct {
	suite.Suite
}

func TestCodecSuite(t *testing.T) {
	suite.Run(t, new(CodecTestSuite))
}

func (suite *CodecTestSuite) TestRegisterInterfaces() {
	registry := codectypes.NewInterfaceRegistry()
	registry.RegisterInterface(sdk.MsgInterfaceProtoName, (*sdk.Msg)(nil))
	RegisterInterfaces(registry)

	impls := registry.ListImplementations(sdk.MsgInterfaceProtoName)
	suite.Require().Equal(4, len(impls))
	suite.Require().ElementsMatch([]string{
		"/furya.feeshare.v1.MsgRegisterFeeShare",
		"/furya.feeshare.v1.MsgCancelFeeShare",
		"/furya.feeshare.v1.MsgUpdateFeeShare",
		"/furya.feeshare.v1.MsgUpdateParams",
	}, impls)
}
