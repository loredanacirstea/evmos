package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/tharsis/ethermint/types"

	"github.com/tharsis/evmos/v3/x/fees/types"
)

var _ types.QueryServer = Keeper{}

// FeeContracts returns all registered contracts for fee distribution
func (k Keeper) FeeContracts(
	c context.Context,
	req *types.QueryFeeContractsRequest,
) (*types.QueryFeeContractsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var feeContracts []types.FeeContract
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixFee)

	pageRes, err := query.Paginate(
		store,
		req.Pagination,
		func(_, value []byte) error {
			var feeContract types.FeeContract
			if err := k.cdc.Unmarshal(value, &feeContract); err != nil {
				return err
			}
			feeContracts = append(feeContracts, feeContract)
			return nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryFeeContractsResponse{
		Fees:       feeContracts,
		Pagination: pageRes,
	}, nil
}

// FeeContract returns a given registered contract
func (k Keeper) FeeContract(
	c context.Context,
	req *types.QueryFeeContractRequest,
) (*types.QueryFeeContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if strings.TrimSpace(req.ContractAddress) == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"contract address is empty",
		)
	}

	// check if the contract is a hex address
	if err := ethermint.ValidateAddress(req.ContractAddress); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid format for contract %s, should be hex ('0x...')", req.ContractAddress,
		)
	}

	feeContract, found := k.GetFee(ctx, common.HexToAddress(req.ContractAddress))
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"fees registered contract '%s'",
			req.ContractAddress,
		)
	}

	return &types.QueryFeeContractResponse{Fee: feeContract}, nil
}

// Params return hub contract param
func (k Keeper) Params(
	c context.Context,
	_ *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// FeeContractsPerDeployer returns all contracts that a deployer has registered
func (k Keeper) FeeContractsPerDeployer(
	c context.Context,
	req *types.QueryFeeContractsPerDeployerRequest,
) (*types.QueryFeeContractsPerDeployerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if strings.TrimSpace(req.DeployerAddress) == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"deployer address is empty",
		)
	}

	deployer, err := sdk.AccAddressFromBech32(req.DeployerAddress)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid format for deployer %s, should be bech32 ('evmos...')", req.DeployerAddress,
		)
	}

	contractAddresses := k.GetFeesInverse(ctx, deployer)
	var feeContracts []types.FeeContract

	for _, contractAddress := range contractAddresses {
		feeContract, found := k.GetFee(ctx, contractAddress)
		if found {
			feeContracts = append(feeContracts, feeContract)
		}
	}

	return &types.QueryFeeContractsPerDeployerResponse{Fees: feeContracts}, nil
}
