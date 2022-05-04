package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
	ibctesting "github.com/tharsis/evmos/v3/ibc/testing"

	app "github.com/tharsis/evmos/v3/app"
)

var (
	// TestAccAddress defines a resuable bech32 address for testing purposes
	// TODO: update crypto.AddressHash() when sdk uses address.Module()
	TestAccAddress = icatypes.GenerateAddress(sdk.AccAddress(crypto.AddressHash([]byte(icatypes.ModuleName))), ibcgotesting.FirstConnectionID, TestPortID)
	// TestOwnerAddress defines a reusable bech32 address for testing purposes
	TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"
	// TestPortID defines a resuable port identifier for testing purposes
	TestPortID, _ = icatypes.NewControllerPortID(TestOwnerAddress)
	// TestVersion defines a resuable interchainaccounts version string for testing purposes
	TestVersion = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: ibcgotesting.FirstConnectionID,
		HostConnectionId:       ibcgotesting.FirstConnectionID,
	}))
)

// func init() {
// 	ibcgotesting.DefaultTestingAppInit = SetupICATestingApp
// }

// func SetupICATestingApp() (ibcgotesting.TestingApp, map[string]json.RawMessage) {
// 	return app.Setup(false, nil), app.NewDefaultGenesisState()
// }

// KeeperTestSuite is a testing suite to test keeper functions
type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibcgotesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibcgotesting.TestChain
	chainB *ibcgotesting.TestChain
}

func (suite *KeeperTestSuite) GetApp(chain *ibcgotesting.TestChain) *app.Evmos {
	app, ok := chain.App.(*app.Evmos)
	if !ok {
		panic("not evmos app")
	}

	return app
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2, 0)
	suite.chainA = suite.coordinator.GetChain(ibcgotesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibcgotesting.GetChainID(2))
}

func (suite *KeeperTestSuite) NewICAPath(chainA, chainB *ibcgotesting.TestChain) *ibcgotesting.Path {
	path := ibcgotesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = icatypes.PortID
	path.EndpointB.ChannelConfig.PortID = icatypes.PortID
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointA.ChannelConfig.Version = TestVersion
	path.EndpointB.ChannelConfig.Version = TestVersion

	return path
}

func (suite *KeeperTestSuite) RegisterInterchainAccount(endpoint *ibcgotesting.Endpoint, owner string) error {
	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return err
	}

	channelSequence := endpoint.Chain.App.GetIBCKeeper().ChannelKeeper.GetNextChannelSequence(endpoint.Chain.GetContext())

	if err := suite.GetApp(endpoint.Chain).ICAControllerKeeper.RegisterInterchainAccount(endpoint.Chain.GetContext(), endpoint.ConnectionID, owner); err != nil {
		return err
	}

	// commit state changes for proof verification
	endpoint.Chain.NextBlock()

	// update port/channel ids
	endpoint.ChannelID = channeltypes.FormatChannelIdentifier(channelSequence)
	endpoint.ChannelConfig.PortID = portID
	endpoint.ChannelConfig.Version = TestVersion

	return nil
}

// SetupICAPath invokes the InterchainAccounts entrypoint and subsequent channel handshake handlers
func (suite *KeeperTestSuite) SetupICAPath(path *ibcgotesting.Path, owner string) error {
	if err := suite.RegisterInterchainAccount(path.EndpointA, owner); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenTry(); err != nil {
		return err
	}

	if err := path.EndpointA.ChanOpenAck(); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenConfirm(); err != nil {
		return err
	}

	return nil
}

func (suite *KeeperTestSuite) TestOnChanCloseInit() {
	path := suite.NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)

	err := suite.SetupICAPath(path, TestOwnerAddress)
	suite.Require().NoError(err)

	module, _, err := suite.chainA.App.GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID)
	suite.Require().NoError(err)

	cbs, ok := suite.chainA.App.GetIBCKeeper().Router.GetRoute(module)
	suite.Require().True(ok)

	err = cbs.OnChanCloseInit(
		suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID,
	)

	suite.Require().Error(err)
}

// package keeper_test

// import (
// 	"testing"
// 	"time"

// 	"github.com/ethereum/go-ethereum/common"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"

// 	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
// 	"github.com/tharsis/ethermint/tests"
// 	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

// 	"github.com/cosmos/cosmos-sdk/baseapp"
// 	"github.com/cosmos/cosmos-sdk/crypto/keyring"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
// 	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"
// 	abci "github.com/tendermint/tendermint/abci/types"
// 	"github.com/tendermint/tendermint/crypto/tmhash"
// 	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
// 	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
// 	"github.com/tendermint/tendermint/version"
// 	evm "github.com/tharsis/ethermint/x/evm/types"
// 	"github.com/tharsis/evmos/v3/app"
// 	claimtypes "github.com/tharsis/evmos/v3/x/claims/types"
// )

// type KeeperTestSuite struct {
// 	suite.Suite

// 	ctx sdk.Context

// 	app            *app.Evmos
// 	queryClientEvm evm.QueryClient
// 	address        common.Address
// 	signer         keyring.Signer
// 	ethSigner      ethtypes.Signer
// 	consAddress    sdk.ConsAddress
// 	validator      stakingtypes.Validator
// 	denom          string
// }

// var s *KeeperTestSuite

// func TestKeeperTestSuite(t *testing.T) {
// 	s = new(KeeperTestSuite)
// 	suite.Run(t, s)

// 	// Run Ginkgo integration tests
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "Keeper Suite")
// }

// func (suite *KeeperTestSuite) SetupTest() {
// 	t := suite.T()
// 	// account key
// 	priv, err := ethsecp256k1.GenerateKey()
// 	require.NoError(t, err)
// 	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
// 	suite.signer = tests.NewSigner(priv)

// 	suite.denom = claimtypes.DefaultClaimsDenom

// 	// consensus key
// 	privCons, err := ethsecp256k1.GenerateKey()
// 	require.NoError(t, err)
// 	suite.consAddress = sdk.ConsAddress(privCons.PubKey().Address())
// 	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState())
// 	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
// 		Height:          1,
// 		ChainID:         "evmos_9001-1",
// 		Time:            time.Now().UTC(),
// 		ProposerAddress: suite.consAddress.Bytes(),

// 		Version: tmversion.Consensus{
// 			Block: version.BlockProtocol,
// 		},
// 		LastBlockId: tmproto.BlockID{
// 			Hash: tmhash.Sum([]byte("block_id")),
// 			PartSetHeader: tmproto.PartSetHeader{
// 				Total: 11,
// 				Hash:  tmhash.Sum([]byte("partset_header")),
// 			},
// 		},
// 		AppHash:            tmhash.Sum([]byte("app")),
// 		DataHash:           tmhash.Sum([]byte("data")),
// 		EvidenceHash:       tmhash.Sum([]byte("evidence")),
// 		ValidatorsHash:     tmhash.Sum([]byte("validators")),
// 		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
// 		ConsensusHash:      tmhash.Sum([]byte("consensus")),
// 		LastResultsHash:    tmhash.Sum([]byte("last_result")),
// 	})

// 	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
// 	evm.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
// 	suite.queryClientEvm = evm.NewQueryClient(queryHelperEvm)

// 	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
// 	stakingParams.BondDenom = suite.denom
// 	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

// 	evmParams := suite.app.EvmKeeper.GetParams(suite.ctx)
// 	evmParams.EvmDenom = suite.denom
// 	suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)

// 	inflationParams := suite.app.InflationKeeper.GetParams(suite.ctx)
// 	inflationParams.EnableInflation = false
// 	suite.app.InflationKeeper.SetParams(suite.ctx, inflationParams)

// 	// Set Validator
// 	valAddr := sdk.ValAddress(suite.address.Bytes())
// 	validator, err := stakingtypes.NewValidator(valAddr, privCons.PubKey(), stakingtypes.Description{})
// 	require.NoError(t, err)
// 	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
// 	suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
// 	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
// 	require.NoError(t, err)
// 	validators := s.app.StakingKeeper.GetValidators(s.ctx, 1)
// 	suite.validator = validators[0]

// 	suite.ethSigner = ethtypes.LatestSignerForChainID(s.app.EvmKeeper.ChainID())
// }

// // Commit commits and starts a new block with an updated context.
// func (suite *KeeperTestSuite) Commit() {
// 	suite.CommitAfter(time.Second * 0)
// }

// // Commit commits a block at a given time.
// func (suite *KeeperTestSuite) CommitAfter(t time.Duration) {
// 	header := suite.ctx.BlockHeader()
// 	suite.app.EndBlock(abci.RequestEndBlock{Height: header.Height})
// 	_ = suite.app.Commit()

// 	header.Height += 1
// 	header.Time = header.Time.Add(t)
// 	suite.app.BeginBlock(abci.RequestBeginBlock{
// 		Header: header,
// 	})

// 	// update ctx
// 	suite.ctx = suite.app.BaseApp.NewContext(false, header)

// 	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
// 	evm.RegisterQueryServer(queryHelperEvm, suite.app.EvmKeeper)
// 	suite.queryClientEvm = evm.NewQueryClient(queryHelperEvm)
// }
