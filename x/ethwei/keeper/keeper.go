package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"ethwei/x/ethwei/types"
)

type Keeper struct {
	storeService   corestore.KVStoreService
	cdc            codec.Codec
	addressCodec   address.Codec
	authority      []byte
	mintBankKeeper types.MintBankKeeper

	Schema      collections.Schema
	Params      collections.Item[types.Params]
	TotalMinted collections.Item[uint64]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	mintBankKeeper types.MintBankKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:   storeService,
		cdc:            cdc,
		addressCodec:   addressCodec,
		authority:      authority,
		mintBankKeeper: mintBankKeeper,

		Params:      collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		TotalMinted: collections.NewItem(sb, types.TotalMintedKey, "total_minted", collections.Uint64Value),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
