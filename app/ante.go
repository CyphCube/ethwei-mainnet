package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ValidatorProposalDecorator enforces two Ethwei governance rules:
//  1. Only active (bonded) validators may submit proposals.
//  2. The NoWithVeto vote option is disabled.
type ValidatorProposalDecorator struct {
	sk *stakingkeeper.Keeper
}

func NewValidatorProposalDecorator(sk *stakingkeeper.Keeper) ValidatorProposalDecorator {
	return ValidatorProposalDecorator{sk: sk}
}

func (d ValidatorProposalDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *govv1.MsgSubmitProposal:
			proposer, err := sdk.AccAddressFromBech32(m.Proposer)
			if err != nil {
				return ctx, sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", err)
			}
			val, err := d.sk.GetValidator(ctx, sdk.ValAddress(proposer))
			if err != nil || val.Status != stakingtypes.Bonded {
				return ctx, sdkerrors.ErrUnauthorized.Wrap("only active validators can submit proposals on Ethwei")
			}

		case *govv1.MsgVote:
			if m.Option == govv1.VoteOption_VOTE_OPTION_NO_WITH_VETO {
				return ctx, sdkerrors.ErrInvalidRequest.Wrap("NoWithVeto is not supported on Ethwei — use No instead")
			}
		}
	}
	return next(ctx, tx, simulate)
}
