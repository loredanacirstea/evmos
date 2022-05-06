package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/evmos/v4/x/cronjobs/types"
)

// GetAllCronjobs - get all registered DevFeeInfo instances
func (k Keeper) GetAllCronjobs(ctx sdk.Context) []types.Cronjob {
	values := []types.Cronjob{}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixCronjob)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var value types.Cronjob
		k.cdc.MustUnmarshal(iterator.Value(), &value)
		values = append(values, value)
	}

	return values
}

// IterateCronjobs iterates over all registered contracts and performs a
// callback with the corresponding DevFeeInfo.
func (k Keeper) IterateCronjobs(
	ctx sdk.Context,
	handlerFn func(fee types.Cronjob) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixCronjob)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var value types.Cronjob
		k.cdc.MustUnmarshal(iterator.Value(), &value)
		if handlerFn(value) {
			break
		}
	}
}

// GetFeeInfo returns DevFeeInfo for a registered contract
func (k Keeper) GetCronjob(ctx sdk.Context, id string) (types.Cronjob, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		return types.Cronjob{}, false
	}

	var value types.Cronjob
	k.cdc.MustUnmarshal(bz, &value)
	return value, true
}

// SetFee stores the developer fee information for a registered contract
func (k Keeper) SetCronjob(ctx sdk.Context, id string, cronjob types.Cronjob) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)
	bz := k.cdc.MustMarshal(&cronjob)
	store.Set([]byte(id), bz)
}

// DeleteFee removes a registered contract
func (k Keeper) DeleteCronjob(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)
	store.Delete([]byte(id))
}

// IsFeeRegistered checks if a contract was registered for receiving fees
func (k Keeper) IsCronjobRegistered(
	ctx sdk.Context,
	id string,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)
	return store.Has([]byte(id))
}
