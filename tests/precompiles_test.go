package keeper_test

// import (
// 	"fmt"
// 	"math/big"
// 	"testing"

// 	"github.com/stretchr/testify/suite"

// 	"github.com/cosmos/cosmos-sdk/client/tx"
// 	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/types/tx/signing"
// 	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
// 	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
// 	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

// 	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"

// 	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
// 	"github.com/tharsis/ethermint/encoding"
// 	"github.com/tharsis/ethermint/tests"
// 	evmtypes "github.com/tharsis/ethermint/x/evm/types"
// 	intertx "github.com/tharsis/ethermint/x/inter-tx/types"
// 	"github.com/tharsis/evmos/v3/app"

// 	"github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"

// 	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
// 	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
// 	abci "github.com/tendermint/tendermint/abci/types"
// 	ibctesting "github.com/tharsis/evmos/v3/ibc/testing"
// )

// type IBCTestingSuite struct {
// 	suite.Suite
// 	coordinator *ibcgotesting.Coordinator

// 	// testing chains used for convenience and readability
// 	chainA *ibcgotesting.TestChain // Ethermint chain A
// 	chainB *ibcgotesting.TestChain // Ethermint chain B
// 	// chainCosmos *ibcgotesting.TestChain // Cosmos chain

// 	// pathEVM *ibcgotesting.Path // chainA (Ethermint) <-->  chainB (Ethermint)
// 	pathICA *ibcgotesting.Path // chainA (Ethermint) <-->  chainB (Ethermint)

// 	ethSigner1  ethtypes.Signer
// 	ethSigner2  ethtypes.Signer
// 	priv        *ethsecp256k1.PrivKey
// 	address     sdk.AccAddress
// 	TestPortID  string
// 	TestVersion string
// }

// func (suite *IBCTestingSuite) SetupTest() {
// 	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2, 0) // initializes 2 Ethermint test chains and 1 Cosmos Chain
// 	suite.chainA = suite.coordinator.GetChain(ibcgotesting.GetChainID(1))
// 	suite.chainB = suite.coordinator.GetChain(ibcgotesting.GetChainID(2))

// 	suite.coordinator.CommitNBlocks(suite.chainA, 2)
// 	suite.coordinator.CommitNBlocks(suite.chainB, 2)

// 	chainA := suite.chainA.App.(*app.Evmos)
// 	chainB := suite.chainB.App.(*app.Evmos)
// 	suite.ethSigner1 = ethtypes.LatestSignerForChainID(chainA.EvmKeeper.ChainID())
// 	suite.ethSigner2 = ethtypes.LatestSignerForChainID(chainB.EvmKeeper.ChainID())

// 	coins := sdk.NewCoins(sdk.NewCoin("aevmos", sdk.NewInt(10000)))
// 	err := FundModuleAccount(suite.chainB.App.(*app.Evmos).BankKeeper, suite.chainB.GetContext(), evmtypes.ModuleName, coins)
// 	suite.Require().NoError(err)

// 	err = FundModuleAccount(suite.chainA.App.(*app.Evmos).BankKeeper, suite.chainA.GetContext(), evmtypes.ModuleName, coins)
// 	suite.Require().NoError(err)

// 	priv, err := ethsecp256k1.GenerateKey()
// 	suite.Require().NoError(err)
// 	suite.priv = priv
// 	suite.address = sdk.AccAddress(priv.PubKey().Address())
// 	fmt.Println("---address", suite.address, suite.address.String())

// 	portid, err := icatypes.NewControllerPortID(suite.address.String())
// 	fmt.Println("---portid", portid)
// 	suite.Require().NoError(err)
// 	suite.TestPortID = portid

// 	suite.TestVersion = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
// 		Version:                icatypes.Version,
// 		ControllerConnectionId: ibcgotesting.FirstConnectionID,
// 		HostConnectionId:       ibcgotesting.FirstConnectionID,
// 	}))

// 	// suite.chainA.App.OnChanOpenInit()
// 	// suite.chainA.App.ibc

// 	// appA, ok := suite.chainA.App.(*simapp.SimApp)
// 	// suite.Require().True(ok)
// 	// appA := suite.chainA.App

// 	// appA.ICAAuthModule.IBCApp.OnChanOpenInit = func(ctx sdk.Context, order channeltypes.Order, connectionHops []string,
// 	// 	portID, channelID string, chanCap *capabilitytypes.Capability,
// 	// 	counterparty channeltypes.Counterparty, version string,
// 	// ) error {
// 	// 	fmt.Println("-----mock OnChanOpenInit")
// 	// 	return fmt.Errorf("mock ica auth fails")
// 	// }

// 	// suite.chainA.GetSimApp().ICAAuthModule.IBCApp.OnChanOpenInit = func(ctx sdk.Context, order channeltypes.Order, connectionHops []string,
// 	// 	portID, channelID string, chanCap *capabilitytypes.Capability,
// 	// 	counterparty channeltypes.Counterparty, version string,
// 	// ) error {
// 	// 	fmt.Println("-----mock OnChanOpenInit")
// 	// 	return fmt.Errorf("mock ica auth fails")
// 	// }

// 	// suite.pathEVM = ibctesting.NewTransferPath(suite.chainA, suite.chainB) // clientID, connectionID, channelID empty
// 	// suite.coordinator.Setup(suite.pathEVM)                                 // clientID, connectionID, channelID filled
// 	// suite.Require().Equal("07-tendermint-0", suite.pathEVM.EndpointA.ClientID)
// 	// suite.Require().Equal("connection-0", suite.pathEVM.EndpointA.ConnectionID)
// 	// suite.Require().Equal("channel-0", suite.pathEVM.EndpointA.ChannelID)

// 	suite.pathICA = NewICAPath(suite.chainA, suite.chainB, suite.TestVersion)
// 	suite.coordinator.SetupConnections(suite.pathICA)
// 	// suite.coordinator.Setup(suite.pathICA)
// 	suite.Require().Equal("07-tendermint-0", suite.pathICA.EndpointA.ClientID)
// 	suite.Require().Equal("connection-0", suite.pathICA.EndpointA.ConnectionID)
// 	// suite.Require().Equal("channel-0", suite.pathICA.EndpointA.ChannelID)

// 	// err = suite.pathICA.EndpointA.ChanOpenInit()
// 	// suite.Require().NoError(err)

// 	// err = suite.pathICA.EndpointB.ChanOpenTry()
// 	// suite.Require().NoError(err)

// 	// err = suite.pathICA.EndpointA.ChanOpenAck()
// 	// suite.Require().NoError(err)

// 	// err = suite.pathICA.EndpointB.ChanOpenConfirm()
// 	// suite.Require().NoError(err)

// 	// // ensure counterparty is up to date
// 	// err = suite.pathICA.EndpointA.UpdateClient()
// 	// suite.Require().NoError(err)

// 	// msgSrv := keeper.NewMsgServerImpl(suite.GetICAApp(suite.chainA).InterTxKeeper)
// 	// msg := types.NewMsgRegisterAccount(owner, path.EndpointA.ConnectionID)

// 	// res, err := msgSrv.RegisterAccount(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

// 	// transfer := transfertypes.NewFungibleTokenPacketData("aevmos", "100", sender, receiver)
// 	// bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
// 	// packet := channeltypes.NewPacket(bz, 1, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)

// 	// // send on endpointA
// 	// err := path.EndpointA.SendPacket(packet)
// 	// suite.Require().NoError(err)

// 	// // receive on endpointB
// 	// err = path.RelayPacket(packet)
// 	// suite.Require().NoError(err)

// 	registerMsg := intertx.NewMsgRegisterAccount(suite.address.String(), suite.pathICA.EndpointA.ConnectionID)
// 	res, err := chainA.InterTxKeeper.RegisterAccount(suite.chainA.GetContext(), registerMsg)
// 	fmt.Println("---register", res)
// 	suite.Require().NoError(err)
// 	suite.coordinator.CommitBlock(suite.chainA)

// 	// openTry := channeltypes.NewMsgRecvPacket()
// 	// bz := channeltypes.SubModuleCdc.MustMarshalJSON(&openTry)
// 	// packetTry := channeltypes.NewPacket(
// 	// 	bz,
// 	// 	1,
// 	// 	suite.pathICA.EndpointA.ChannelConfig.PortID,
// 	// 	suite.pathICA.EndpointA.ChannelID,
// 	// 	suite.pathICA.EndpointB.ChannelConfig.PortID,
// 	// 	suite.pathICA.EndpointB.ChannelID,
// 	// 	timeoutHeight,
// 	// 	0,
// 	// )
// 	// // send on endpointB
// 	// err = suite.pathICA.EndpointB.SendPacket(packetTry)
// 	// suite.Require().NoError(err)

// 	// openTryMsg := &channeltypes.MsgChannelOpenTry{
// 	// 	PortId: suite.TestPortID,
// 	// 	PreviousChannelId: "",
// 	// 	Channel: "",
// 	// 	CounterpartyVersion: suite.TestVersion,
// 	// 	ProofInit: ,
// 	// 	ProofHeight: ,
// 	// 	Signer: ,
// 	// }

// 	// // suite.chainA.App.GetIBCKeeper().ChannelKeeper.ChanOpenTry(
// 	// chainB.IBCKeeper.ChannelOpenTry(
// 	// 	sdk.WrapSDKContext(suite.chainB.GetContext()),
// 	// 	openTryMsg,
// 	// )

// 	// suite.pathICA.EndpointA.Chain.Coordinator.UpdateTimeForChain(suite.pathICA.EndpointA.Chain)
// 	// suite.pathICA.EndpointA.Chain.App.GetBaseApp().Commit()
// 	// suite.pathICA.EndpointA.Chain.NextBlock()
// 	// suite.pathICA.EndpointA.Chain.Coordinator.IncrementTime()

// 	suite.coordinator.CommitNBlocks(suite.chainA, 1)
// 	// suite.coordinator.CommitNBlocks(suite.chainB, 1)

// 	fmt.Println("--.EndpointA PortID", suite.pathICA.EndpointA.ChannelConfig.PortID)
// 	fmt.Println("--.EndpointA Version", suite.pathICA.EndpointA.ChannelConfig.Version)
// 	fmt.Println("--.EndpointA Order", suite.pathICA.EndpointA.ChannelConfig.Order)
// 	fmt.Println("--.EndpointA ChannelID", suite.pathICA.EndpointA.ChannelID)

// 	channels := chainA.IBCKeeper.ChannelKeeper.GetAllChannels(suite.chainA.GetContext())
// 	fmt.Println("channels", channels)

// 	// // send ChanOpenTry message to ChainB
// 	// suite.pathICA.EndpointA.ChannelID = "channel-0"
// 	// suite.pathICA.EndpointB.Counterparty.ChannelConfig.PortID = suite.pathICA.EndpointA.ChannelConfig.PortID
// 	// suite.pathICA.EndpointB.Counterparty.ChannelConfig.Version = suite.pathICA.EndpointA.ChannelConfig.Version
// 	// suite.pathICA.EndpointB.Counterparty.ChannelConfig.Order = suite.pathICA.EndpointA.ChannelConfig.Order
// 	// suite.pathICA.EndpointB.Counterparty.ChannelID = suite.pathICA.EndpointA.ChannelID
// 	// err = suite.pathICA.EndpointB.ChanOpenTry()
// 	// suite.Require().NoError(err)

// 	// ica, err := chainB.InterTxKeeper.InterchainAccountFromAddressInner(suite.chainB.GetContext(), &intertx.QueryInterchainAccountFromAddressRequest{Owner: suite.address.String(), ConnectionId: "connection-0"})
// 	// suite.Require().NoError(err)
// 	// suite.Require().Equal(suite.TestPortID, ica.InterchainAccountAddress)

// 	// fmt.Println("---ica", ica)
// }

// func TestIBCTestingSuite(t *testing.T) {
// 	suite.Run(t, new(IBCTestingSuite))
// }

// func (suite *IBCTestingSuite) TestIcaEthereumNoSignTx() {
// 	ica := suite.TestPortID
// 	chainA := suite.chainA.App.(*app.Evmos)
// 	fmt.Println("----ica----", ica)
// 	nonce := uint64(0)
// 	value := big.NewInt(0)
// 	to := common.BytesToAddress(suite.address.Bytes())
// 	gasLimit := uint64(300000)
// 	gasPrice := big.NewInt(20)
// 	gasFeeCap := big.NewInt(20)
// 	gasTipCap := big.NewInt(20)
// 	data := make([]byte, 0)
// 	accesses := &ethtypes.AccessList{}
// 	ethtx := evmtypes.NewIcaTx(chainA.EvmKeeper.ChainID(), nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)

// 	// ethicatx := (*types.MsgEthereumIcaTx)(ethtx)
// 	// ethicatx := (*types.MsgEthereumIcaTx)(*ethtx.(*types.MsgEthereumTx))

// 	icaAddr, err := sdk.AccAddressFromBech32(ica)
// 	fmt.Println("ethtx.From", icaAddr, err)
// 	ethtx.From = common.BytesToAddress(icaAddr.Bytes()).Hex()
// 	fmt.Println("ethtx.From", ethtx.From)

// 	res, err := chainA.EvmKeeper.EthereumIcaTx(sdk.WrapSDKContext(suite.chainA.GetContext()), ethtx)
// 	fmt.Println("res", res)
// }

// // func (suite *IBCTestingSuite) TestIcaSubmitTx() {
// // 	ica := suite.TestPortID
// // 	chainA := suite.chainA.App.(*app.Evmos)
// // 	fmt.Println("----ica----", ica)
// // 	nonce := uint64(0)
// // 	value := big.NewInt(0)
// // 	to := common.BytesToAddress(suite.address.Bytes())
// // 	gasLimit := uint64(300000)
// // 	gasPrice := big.NewInt(20)
// // 	gasFeeCap := big.NewInt(20)
// // 	gasTipCap := big.NewInt(20)
// // 	data := make([]byte, 0)
// // 	accesses := &ethtypes.AccessList{}
// // 	ethtx := evmtypes.NewIcaTx(chainA.EvmKeeper.ChainID(), nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)

// // 	// ethicatx := (*types.MsgEthereumIcaTx)(ethtx)
// // 	// ethicatx := (*types.MsgEthereumIcaTx)(*ethtx.(*types.MsgEthereumTx))

// // 	icaAddr, err := sdk.AccAddressFromBech32(ica)
// // 	fmt.Println("ethtx.From", icaAddr, err)
// // 	ethtx.From = common.BytesToAddress(icaAddr.Bytes()).Hex()
// // 	fmt.Println("ethtx.From", ethtx.From)

// // 	msg, err := intertx.NewMsgSubmitTx(ethtx, "connection-0", suite.address.String())
// // 	suite.Require().NoError(err)

// // 	res, err := chainA.InterTxKeeper.SubmitTx(suite.chainA.GetContext(), msg)
// // 	fmt.Println("----SubmitTx res----", res, err)
// // 	suite.Require().NoError(err)
// // }

// // func (suite *IBCTestingSuite) TestIcaPrecompile() {
// // 	suite.SetupTest()
// // 	chainA := suite.chainA.App.(*app.Evmos)
// // 	// chainB := suite.chainB.App.(*app.Evmos)

// // 	chainID := chainA.EvmKeeper.ChainID()
// // 	from := suite.address
// // 	nonce := getNonce(chainA, suite.chainA.GetContext(), from.Bytes())
// // 	data := make([]byte, 0)
// // 	gasLimit := uint64(100000)
// // 	gasPrice := big.NewInt(2000000000)
// // 	gasFeeCap := big.NewInt(2000000000)
// // 	gasTipCap := big.NewInt(2000000000)
// // 	precompileAddress := common.HexToAddress("0x0000000000000000000000000000000000000019")
// // 	msgEthereumTx := evmtypes.NewTx(
// // 		chainID,
// // 		nonce,
// // 		&precompileAddress,
// // 		nil,
// // 		gasLimit,
// // 		gasPrice,
// // 		gasFeeCap,
// // 		gasTipCap,
// // 		data,
// // 		nil,
// // 	)
// // 	msgEthereumTx.From = from.String()

// // 	res := suite.performEthTx(chainA, suite.chainA.GetContext(), suite.ethSigner1, suite.priv, msgEthereumTx)
// // 	suite.coordinator.CommitNBlocks(suite.chainA, 1)
// // 	fmt.Println(res.GetLog())
// // 	suite.Require().True(false)
// // }

// func NewICAPath(chainA, chainB *ibcgotesting.TestChain, TestVersion string) *ibcgotesting.Path {
// 	path := ibcgotesting.NewPath(chainA, chainB)
// 	path.EndpointA.ChannelConfig.PortID = icatypes.PortID
// 	path.EndpointB.ChannelConfig.PortID = icatypes.PortID
// 	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
// 	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
// 	path.EndpointA.ChannelConfig.Version = TestVersion
// 	path.EndpointB.ChannelConfig.Version = TestVersion
// 	return path
// }

// func getNonce(chainApp *app.Evmos, ctx sdk.Context, addressBytes []byte) uint64 {
// 	return chainApp.EvmKeeper.GetNonce(
// 		ctx,
// 		common.BytesToAddress(addressBytes),
// 	)
// }

// // FundModuleAccount is a utility function that funds a module account by
// // minting and sending the coins to the address. This should be used for testing
// // purposes only!
// func FundModuleAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, recipientMod string, amounts sdk.Coins) error {
// 	if err := bankKeeper.MintCoins(ctx, evmtypes.ModuleName, amounts); err != nil {
// 		return err
// 	}

// 	return bankKeeper.SendCoinsFromModuleToModule(ctx, evmtypes.ModuleName, recipientMod, amounts)
// }

// // func TestIBCTestingSuite(t *testing.T) {
// // 	suite.Run(t, new(IBCTestingSuite))
// // }

// // func getNonce(app, ctx, addressBytes []byte) uint64 {
// // 	return s.app.EvmKeeper.GetNonce(
// // 		s.ctx,
// // 		common.BytesToAddress(addressBytes),
// // 	)
// // }

// // func (s *KeeperTestSuite) deployContract(priv *ethsecp256k1.PrivKey, contractCode string) common.Address {
// // 	chainID := s.app.EvmKeeper.ChainID()
// // 	from := common.BytesToAddress(priv.PubKey().Address().Bytes())
// // 	nonce := s.getNonce(from.Bytes())

// // 	data := common.Hex2Bytes(contractCode)
// // 	gasLimit := uint64(100000)
// // 	msgEthereumTx := evmtypes.NewTxContract(
// // 		chainID,
// // 		nonce,
// // 		nil,
// // 		gasLimit,
// // 		nil,
// // 		s.app.FeeMarketKeeper.GetBaseFee(s.ctx),
// // 		big.NewInt(1),
// // 		data,
// // 		&ethtypes.AccessList{},
// // 	)
// // 	msgEthereumTx.From = from.String()

// // 	res := s.performEthTx(priv, msgEthereumTx)
// // 	s.Commit()

// // 	ethereumTx := res.GetEvents()[10]
// // 	s.Require().Equal(ethereumTx.Type, "ethereum_tx")
// // 	s.Require().Equal(string(ethereumTx.Attributes[1].Key), "ethereumTxHash")

// // 	contractAddress := crypto.CreateAddress(from, nonce)
// // 	acc := s.app.EvmKeeper.GetAccountWithoutBalance(s.ctx, contractAddress)
// // 	s.Require().NotEmpty(acc)
// // 	s.Require().True(acc.IsContract())
// // 	return contractAddress
// // }

// // func (s *KeeperTestSuite) contractInteract(
// // 	priv *ethsecp256k1.PrivKey,
// // 	contractAddr *common.Address,
// // 	gasPrice *big.Int,
// // 	gasFeeCap *big.Int,
// // 	gasTipCap *big.Int,
// // 	accesses *ethtypes.AccessList,
// // ) abci.ResponseDeliverTx {
// // 	chainID := s.app.EvmKeeper.ChainID()
// // 	from := common.BytesToAddress(priv.PubKey().Address().Bytes())
// // 	nonce := s.getNonce(from.Bytes())
// // 	data := make([]byte, 0)
// // 	gasLimit := uint64(100000)
// // 	msgEthereumTx := evmtypes.NewTx(
// // 		chainID,
// // 		nonce,
// // 		contractAddr,
// // 		nil,
// // 		gasLimit,
// // 		gasPrice,
// // 		gasFeeCap,
// // 		gasTipCap,
// // 		data,
// // 		accesses,
// // 	)
// // 	msgEthereumTx.From = from.String()

// // 	return s.performEthTx(priv, msgEthereumTx)
// // }

// func (s *IBCTestingSuite) performEthTx(chainApp *app.Evmos, ctx sdk.Context, ethSigner ethtypes.Signer, priv *ethsecp256k1.PrivKey, msgEthereumTx *evmtypes.MsgEthereumTx) abci.ResponseDeliverTx {
// 	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
// 	option, err := codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
// 	s.Require().NoError(err)

// 	txBuilder := encodingConfig.TxConfig.NewTxBuilder()
// 	builder, ok := txBuilder.(authtx.ExtensionOptionsTxBuilder)
// 	s.Require().True(ok)
// 	builder.SetExtensionOptions(option)

// 	err = msgEthereumTx.Sign(ethSigner, tests.NewSigner(priv))
// 	s.Require().NoError(err)

// 	err = txBuilder.SetMsgs(msgEthereumTx)
// 	s.Require().NoError(err)

// 	txData, err := evmtypes.UnpackTxData(msgEthereumTx.Data)
// 	s.Require().NoError(err)

// 	evmDenom := chainApp.EvmKeeper.GetParams(ctx).EvmDenom
// 	fees := sdk.Coins{{Denom: evmDenom, Amount: sdk.NewIntFromBigInt(txData.Fee())}}
// 	builder.SetFeeAmount(fees)
// 	builder.SetGasLimit(msgEthereumTx.GetGas())

// 	// bz are bytes to be broadcasted over the network
// 	bz, err := encodingConfig.TxConfig.TxEncoder()(txBuilder.GetTx())
// 	s.Require().NoError(err)

// 	req := abci.RequestDeliverTx{Tx: bz}
// 	res := chainApp.BaseApp.DeliverTx(req)
// 	s.Require().Equal(res.IsOK(), true, res.GetLog())
// 	return res
// }

// func (s *IBCTestingSuite) deliverTx(chainApp *app.Evmos, ctx sdk.Context, ethSigner ethtypes.Signer, priv *ethsecp256k1.PrivKey, msgs ...sdk.Msg) abci.ResponseDeliverTx {
// 	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
// 	accountAddress := sdk.AccAddress(priv.PubKey().Address().Bytes())

// 	txBuilder := encodingConfig.TxConfig.NewTxBuilder()
// 	evmDenom := chainApp.EvmKeeper.GetParams(ctx).EvmDenom

// 	txBuilder.SetGasLimit(1000000)
// 	txBuilder.SetFeeAmount(sdk.Coins{{Denom: evmDenom, Amount: sdk.NewInt(1)}})
// 	err := txBuilder.SetMsgs(msgs...)
// 	s.Require().NoError(err)

// 	seq, err := chainApp.AccountKeeper.GetSequence(ctx, accountAddress)
// 	s.Require().NoError(err)

// 	// First round: we gather all the signer infos. We use the "set empty
// 	// signature" hack to do that.
// 	sigV2 := signing.SignatureV2{
// 		PubKey: priv.PubKey(),
// 		Data: &signing.SingleSignatureData{
// 			SignMode:  encodingConfig.TxConfig.SignModeHandler().DefaultMode(),
// 			Signature: nil,
// 		},
// 		Sequence: seq,
// 	}

// 	sigsV2 := []signing.SignatureV2{sigV2}

// 	err = txBuilder.SetSignatures(sigsV2...)
// 	s.Require().NoError(err)

// 	// Second round: all signer infos are set, so each signer can sign.
// 	accNumber := chainApp.AccountKeeper.GetAccount(ctx, accountAddress).GetAccountNumber()
// 	signerData := authsigning.SignerData{
// 		ChainID:       ctx.ChainID(),
// 		AccountNumber: accNumber,
// 		Sequence:      seq,
// 	}
// 	sigV2, err = tx.SignWithPrivKey(
// 		encodingConfig.TxConfig.SignModeHandler().DefaultMode(), signerData,
// 		txBuilder, priv, encodingConfig.TxConfig,
// 		seq,
// 	)
// 	s.Require().NoError(err)

// 	sigsV2 = []signing.SignatureV2{sigV2}
// 	err = txBuilder.SetSignatures(sigsV2...)
// 	s.Require().NoError(err)

// 	// bz are bytes to be broadcasted over the network
// 	bz, err := encodingConfig.TxConfig.TxEncoder()(txBuilder.GetTx())
// 	s.Require().NoError(err)

// 	req := abci.RequestDeliverTx{Tx: bz}
// 	res := chainApp.BaseApp.DeliverTx(req)
// 	return res
// }
