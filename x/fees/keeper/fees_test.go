package keeper_test

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tharsis/ethermint/tests"
	"github.com/tharsis/evmos/v3/x/fees/types"
)

type FeeEntry struct {
	contract common.Address
	deployer sdk.AccAddress
	withdraw sdk.AccAddress
}

func (suite *KeeperTestSuite) TestDeployer() {
	contractAddress := tests.GenerateAddress()
	testCases := []struct {
		name     string
		malleate func() (sdk.AccAddress, bool)
	}{
		{
			"ok - deployer is found",
			func() (sdk.AccAddress, bool) {
				deployer := sdk.AccAddress(tests.GenerateAddress().Bytes())
				suite.app.FeesKeeper.SetDeployer(suite.ctx, contractAddress, deployer)
				return deployer, true
			},
		},
		{
			"ok - no deployer",
			func() (sdk.AccAddress, bool) {
				return nil, false
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			expectedDeployer, expectedFound := tc.malleate()
			deployer, found := suite.app.FeesKeeper.GetDeployer(suite.ctx, contractAddress)

			suite.Require().Equal(expectedFound, found)
			suite.Require().Equal(expectedDeployer, deployer)
		})
	}
}

func (suite *KeeperTestSuite) TestDevFee() {
	contractAddress := tests.GenerateAddress()
	deployerAddress := sdk.AccAddress(tests.GenerateAddress().Bytes())
	withdrawalAddress := sdk.AccAddress(tests.GenerateAddress().Bytes())
	testCases := []struct {
		name          string
		contract      common.Address
		deployer      sdk.AccAddress
		withdraw      sdk.AccAddress
		foundDeployer bool
		foundWithdraw bool
		malleate      func(common.Address, sdk.AccAddress, sdk.AccAddress)
	}{
		{
			"with withdraw address",
			contractAddress,
			deployerAddress,
			withdrawalAddress,
			true,
			true,
			func(contract common.Address, deployer sdk.AccAddress, withdraw sdk.AccAddress) {
				suite.app.FeesKeeper.SetFee(suite.ctx, contract, deployer, withdraw)
			},
		},
		{
			"without withdraw address",
			contractAddress,
			deployerAddress,
			nil,
			true,
			false,
			func(contract common.Address, deployer sdk.AccAddress, withdraw sdk.AccAddress) {
				suite.app.FeesKeeper.SetFee(suite.ctx, contract, deployer, withdraw)
			},
		},
		{
			"deployer same as withdraw address",
			contractAddress,
			deployerAddress,
			deployerAddress,
			true,
			false,
			func(contract common.Address, deployer sdk.AccAddress, withdraw sdk.AccAddress) {
				suite.app.FeesKeeper.SetFee(suite.ctx, contract, deployer, withdraw)
			},
		},
		{
			"not registered",
			common.Address{},
			nil,
			nil,
			false,
			false,
			func(contract common.Address, deployer sdk.AccAddress, withdraw sdk.AccAddress) {},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate(tc.contract, tc.deployer, tc.withdraw)

			deployer, found := suite.app.FeesKeeper.GetDeployer(suite.ctx, contractAddress)
			suite.Require().Equal(tc.foundDeployer, found, "wrong deployer found")
			suite.Require().Equal(tc.deployer, deployer, "wrong deployer")

			withdraw, found := suite.app.FeesKeeper.GetWithdrawal(suite.ctx, contractAddress)
			suite.Require().Equal(tc.foundWithdraw, found, "wrong withdraw found")
			if tc.foundWithdraw {
				suite.Require().Equal(tc.withdraw, withdraw, "wrong withdraw")
			}

			feeInfo, found := suite.app.FeesKeeper.GetFeeInfo(suite.ctx, contractAddress)
			suite.Require().Equal(tc.foundDeployer, found, "wrong fee found")
			if found {
				suite.Require().Equal(tc.contract.String(), feeInfo.ContractAddress, "wrong fee contract")
				suite.Require().Equal(tc.deployer.String(), feeInfo.DeployerAddress, "wrong fee deployer")
				if tc.foundWithdraw {
					suite.Require().Equal(tc.withdraw.String(), feeInfo.WithdrawAddress, "wrong fee withdraw")
				} else {
					suite.Require().Equal("", feeInfo.WithdrawAddress, "wrong fee withdraw")
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestIsFeeRegistered() {
	contract := tests.GenerateAddress()
	deployer := sdk.AccAddress(tests.GenerateAddress().Bytes())
	withdraw := sdk.AccAddress(tests.GenerateAddress().Bytes())
	suite.app.FeesKeeper.SetFee(suite.ctx, contract, deployer, withdraw)
	found := suite.app.FeesKeeper.IsFeeRegistered(suite.ctx, contract)
	suite.Require().True(found)

	contract = tests.GenerateAddress()
	found = suite.app.FeesKeeper.IsFeeRegistered(suite.ctx, contract)
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestDeleteDevFee() {
	contract := tests.GenerateAddress()
	deployer := sdk.AccAddress(tests.GenerateAddress().Bytes())
	withdraw := sdk.AccAddress(tests.GenerateAddress().Bytes())
	suite.app.FeesKeeper.SetFee(suite.ctx, contract, deployer, withdraw)
	found := suite.app.FeesKeeper.IsFeeRegistered(suite.ctx, contract)
	suite.Require().True(found)

	suite.app.FeesKeeper.DeleteFee(suite.ctx, contract)
	found = suite.app.FeesKeeper.IsFeeRegistered(suite.ctx, contract)
	suite.Require().False(found, "wrong fee found")
	_, found = suite.app.FeesKeeper.GetDeployer(suite.ctx, contract)
	suite.Require().False(found, "wrong deployer found")
	_, found = suite.app.FeesKeeper.GetWithdrawal(suite.ctx, contract)
	suite.Require().False(found, "wrong withdraw found")
}

func (suite *KeeperTestSuite) TestAllFees() {
	fees := make([]FeeEntry, 5)
	for i := 0; i < 4; i++ {
		fee := FeeEntry{
			contract: tests.GenerateAddress(),
			deployer: sdk.AccAddress(tests.GenerateAddress().Bytes()),
			withdraw: sdk.AccAddress(tests.GenerateAddress().Bytes()),
		}
		fees[i] = fee
		suite.app.FeesKeeper.SetFee(suite.ctx, fee.contract, fee.deployer, fee.withdraw)
	}
	fee := FeeEntry{
		contract: tests.GenerateAddress(),
		deployer: sdk.AccAddress(tests.GenerateAddress().Bytes()),
	}
	fees[4] = fee
	suite.app.FeesKeeper.SetFee(suite.ctx, fee.contract, fee.deployer, nil)

	sort.Slice(fees, func(i int, j int) bool {
		return sort.StringsAreSorted([]string{
			fees[i].contract.String(),
			fees[j].contract.String(),
		})
	})

	// Test GetAllFees
	all := suite.app.FeesKeeper.GetAllFees(suite.ctx)
	suite.Require().Equal(len(fees), len(all))
	for i := 0; i < 5; i++ {
		suite.Require().Equal(
			fees[i].contract.String(),
			all[i].ContractAddress,
			fmt.Sprintf("wrong contract %d", i),
		)
		suite.Require().Equal(
			fees[i].deployer.String(),
			all[i].DeployerAddress,
			fmt.Sprintf("wrong deployer %d", i),
		)
		suite.Require().Equal(
			fees[i].withdraw.String(),
			all[i].WithdrawAddress,
			fmt.Sprintf("wrong withdraw %d", i),
		)
	}

	// Test IterateFees
	i := 0
	suite.app.FeesKeeper.IterateFees(suite.ctx, func(fee types.DevFeeInfo) bool {
		suite.Require().Equal(
			fees[i].contract.String(),
			all[i].ContractAddress,
			fmt.Sprintf("iterate - wrong contract %d", i),
		)
		suite.Require().Equal(
			fees[i].deployer.String(),
			all[i].DeployerAddress,
			fmt.Sprintf("iterate - wrong deployer %d", i),
		)
		suite.Require().Equal(
			fees[i].withdraw.String(),
			all[i].WithdrawAddress,
			fmt.Sprintf("iterate - wrong withdraw %d", i),
		)
		return false
	})
}

func (suite *KeeperTestSuite) TestDeleteFeeInverse() {
	contract := tests.GenerateAddress()
	deployer := sdk.AccAddress(tests.GenerateAddress().Bytes())

	found := suite.app.FeesKeeper.HasFeeInverse(suite.ctx, deployer)
	suite.Require().False(found)

	suite.app.FeesKeeper.SetFeeInverse(suite.ctx, deployer, contract)
	found = suite.app.FeesKeeper.HasFeeInverse(suite.ctx, deployer)
	suite.Require().True(found)
	contracts := suite.app.FeesKeeper.GetFeesInverse(suite.ctx, deployer)
	suite.Require().Equal(contract, contracts[0])

	suite.app.FeesKeeper.DeleteFeeInverse(suite.ctx, deployer, contract)
	found = suite.app.FeesKeeper.HasFeeInverse(suite.ctx, deployer)
	suite.Require().True(found)
	contracts = suite.app.FeesKeeper.GetFeesInverse(suite.ctx, deployer)
	suite.Require().Equal(len(contracts), 0)
}

func (suite *KeeperTestSuite) TestFeesInverse() {
	contract1 := tests.GenerateAddress()
	contract2 := tests.GenerateAddress()
	deployer := sdk.AccAddress(tests.GenerateAddress().Bytes())

	found := suite.app.FeesKeeper.HasFeeInverse(suite.ctx, deployer)
	suite.Require().False(found)

	suite.app.FeesKeeper.SetFeeInverse(suite.ctx, deployer, contract1)
	suite.app.FeesKeeper.SetFeeInverse(suite.ctx, deployer, contract2)

	contracts := suite.app.FeesKeeper.GetFeesInverse(suite.ctx, deployer)
	suite.Require().Equal(len(contracts), 2)
	suite.Require().Equal(contract1, contracts[0])
	suite.Require().Equal(contract2, contracts[1])

	found = suite.app.FeesKeeper.HasFeeInverse(suite.ctx, deployer)
	suite.Require().True(found)

	suite.app.FeesKeeper.DeleteFeeInverse(suite.ctx, deployer, contract1)
	contracts = suite.app.FeesKeeper.GetFeesInverse(suite.ctx, deployer)
	suite.Require().Equal(len(contracts), 1)
	suite.Require().Equal(contract2, contracts[0])
}
