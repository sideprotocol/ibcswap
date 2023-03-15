package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	// ProposalTypeMarketFeeUpdate defines the type for a MarketFeeUpdateProposal
	ProposalTypeMarketFeeUpdate = "MarketFeeUpdate"
)

// Assert MarketFeeUpdateProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &MarketFeeUpdateProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeMarketFeeUpdate)
}

// NewMarketFeeUpdateProposal creates a new community pool spend proposal.
//
//nolint:interfacer
func NewMarketFeeUpdateProposal(title, description string, poolId string, fee uint32) *MarketFeeUpdateProposal {
	return &MarketFeeUpdateProposal{title, description, poolId, fee}
}

// GetTitle returns the title of a community pool spend proposal.
func (csp *MarketFeeUpdateProposal) GetTitle() string { return csp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (csp *MarketFeeUpdateProposal) GetDescription() string { return csp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (csp *MarketFeeUpdateProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (csp *MarketFeeUpdateProposal) ProposalType() string { return ProposalTypeMarketFeeUpdate }

// ValidateBasic runs basic stateless validity checks
func (csp *MarketFeeUpdateProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(csp)
	if err != nil {
		return err
	}
	// if !csp.Amount.IsValid() {
	// 	return ErrInvalidProposalAmount
	// }
	// if csp.Recipient == "" {
	// 	return ErrEmptyProposalRecipient
	// }

	return nil
}
