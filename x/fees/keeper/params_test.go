package keeper_test

import "github.com/tharsis/evmos/v3/x/fees/types"

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.FeesKeeper.GetParams(suite.ctx)
	params.EnableFees = false
	suite.Require().Equal(types.DefaultParams(), params)
	params.EnableFees = true
	suite.app.FeesKeeper.SetParams(suite.ctx, params)
	newParams := suite.app.FeesKeeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}

func (suite *KeeperTestSuite) TestParamsIsEnabled() {
	params := suite.app.FeesKeeper.GetParams(suite.ctx)
	params.EnableFees = false
	suite.app.FeesKeeper.SetParams(suite.ctx, params)
	suite.Require().False(suite.app.FeesKeeper.IsEnabled(suite.ctx), "fail - not enabled")

	params.EnableFees = true
	suite.app.FeesKeeper.SetParams(suite.ctx, params)
	suite.Require().True(suite.app.FeesKeeper.IsEnabled(suite.ctx), "fail - enabled")
}
