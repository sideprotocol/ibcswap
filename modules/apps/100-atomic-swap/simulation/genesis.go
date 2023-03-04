package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

// Simulation parameter constants
const port = "port_id"

// RadomEnabled randomized send or receive enabled param with 75% prob of being true.
func RadomEnabled(r *rand.Rand) bool {
	return r.Int63n(101) <= 75
}

func RadomInt(r *rand.Rand) uint32 {
	return uint32(r.Int63n(101))
}

// RandomizedGenState generates a random GenesisState for transfer.
func RandomizedGenState(simState *module.SimulationState) {
	var portID string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, port, &portID, simState.Rand,
		func(r *rand.Rand) { portID = strings.ToLower(simtypes.RandStringOfLength(r, 20)) },
	)

	var swapEnabled bool
	simState.AppParams.GetOrGenerate(
		simState.Cdc, string(types.KeySwapEnabled), &swapEnabled, simState.Rand,
		func(r *rand.Rand) { swapEnabled = RadomEnabled(r) },
	)

	var swapMaxFeeRate uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, string(types.KeySwapMaxFeeRate), &swapMaxFeeRate, simState.Rand,
		func(r *rand.Rand) { swapMaxFeeRate = RadomInt(r) },
	)

	transferGenesis := types.GenesisState{
		PortId: portID,
		Params: types.NewParams(swapEnabled, swapMaxFeeRate),
	}

	bz, err := json.MarshalIndent(&transferGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&transferGenesis)
}
