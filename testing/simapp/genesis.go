package simapp

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	defaultGenesis := ModuleBasics.DefaultGenesis(cdc)
	tokenDenom1 := "marscoin"
	tokenDenom2 := "venuscoin"
	defaultGenesis[banktypes.ModuleName] = cdc.MustMarshalJSON(&banktypes.GenesisState{
		Params:        banktypes.Params{},
		Balances:      []banktypes.Balance{{Address: "validator2_address", Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500000000)), sdk.NewCoin(tokenDenom1, sdk.NewInt(250000000)), sdk.NewCoin(tokenDenom2, sdk.NewInt(125000000)))}, {Address: "validator2_address", Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500000000)), sdk.NewCoin("mytoken1", sdk.NewInt(250000000)), sdk.NewCoin("mytoken2", sdk.NewInt(125000000)))}},
		Supply:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000)), sdk.NewCoin(tokenDenom1, sdk.NewInt(500000000)), sdk.NewCoin(tokenDenom2, sdk.NewInt(250000000))),
		DenomMetadata: []banktypes.Metadata{},
	})
	return defaultGenesis
	// return GenesisState{
	// 	banktypes.ModuleName: cdc.MustMarshalJSON(&banktypes.GenesisState{
	// 		Params:        banktypes.Params{},
	// 		Balances:      []banktypes.Balance{{Address: "validator2_address", Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500000000)), sdk.NewCoin(tokenDenom1, sdk.NewInt(250000000)), sdk.NewCoin(tokenDenom2, sdk.NewInt(125000000)))}, {Address: "validator2_address", Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500000000)), sdk.NewCoin("mytoken1", sdk.NewInt(250000000)), sdk.NewCoin("mytoken2", sdk.NewInt(125000000)))}},
	// 		Supply:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000)), sdk.NewCoin(tokenDenom1, sdk.NewInt(500000000)), sdk.NewCoin(tokenDenom2, sdk.NewInt(250000000))),
	// 		DenomMetadata: []banktypes.Metadata{},
	// 	}),
	// }

	//return ModuleBasics.DefaultGenesis(cdc)
}
