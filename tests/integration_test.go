package keeper_test

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	// . "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

func (suite *KeeperTestSuite) TestIcaEthereumNoSignTx() {
	ica := "evmos1m9d75kmeagp4pl48ldy5m4akrjsvt83qunurzpp5gvwmhqhmgmzqdx9sgz"
	nonce := uint64(0)
	value := big.NewInt(0)
	to := common.BytesToAddress(s.address.Bytes())
	gasLimit := uint64(300000)
	gasPrice := big.NewInt(20)
	gasFeeCap := big.NewInt(20)
	gasTipCap := big.NewInt(20)
	data := make([]byte, 0)
	accesses := &ethtypes.AccessList{}
	ethtx := evmtypes.NewIcaTx(s.app.EvmKeeper.ChainID(), nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)

	icaAddr, err := sdk.AccAddressFromBech32(ica)
	fmt.Println("ethtx.From", icaAddr, err)
	ethtx.From = common.BytesToAddress(icaAddr.Bytes()).Hex()
	fmt.Println("ethtx.From", ethtx.From)

	res, err := s.app.EvmKeeper.EthereumIcaTx(sdk.WrapSDKContext(s.ctx), ethtx)
	fmt.Println("res", res)
	s.Require().True(false)
}

// var _ = Describe("Fee distribution:", Ordered, func() {
// 	Context("with fees param disabled", func() {
// 		It("should", func() {
// 			ica := "evmos1m9d75kmeagp4pl48ldy5m4akrjsvt83qunurzpp5gvwmhqhmgmzqdx9sgz"
// 			nonce := uint64(0)
// 			value := big.NewInt(0)
// 			to := common.BytesToAddress(s.address.Bytes())
// 			gasLimit := uint64(300000)
// 			gasPrice := big.NewInt(20)
// 			gasFeeCap := big.NewInt(20)
// 			gasTipCap := big.NewInt(20)
// 			data := make([]byte, 0)
// 			accesses := &ethtypes.AccessList{}
// 			ethtx := evmtypes.NewIcaTx(s.app.EvmKeeper.ChainID(), nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)

// 			icaAddr, err := sdk.AccAddressFromBech32(ica)
// 			fmt.Println("ethtx.From", icaAddr, err)
// 			ethtx.From = common.BytesToAddress(icaAddr.Bytes()).Hex()
// 			fmt.Println("ethtx.From", ethtx.From)

// 			res, err := s.app.EvmKeeper.EthereumIcaTx(sdk.WrapSDKContext(s.ctx), ethtx)
// 			fmt.Println("res", res)
// 			s.Require().True(false)
// 		})
// 	})
// })
