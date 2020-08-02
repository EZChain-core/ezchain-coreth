package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/ava-labs/coreth"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/eth"
	"github.com/ava-labs/coreth/params"
	"github.com/ava-labs/go-ethereum/common"
	"github.com/ava-labs/go-ethereum/common/hexutil"
	"github.com/ava-labs/go-ethereum/rlp"
	"math/big"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// configure the chain
	config := eth.DefaultConfig
	config.ManualCanonical = true
	chainConfig := &params.ChainConfig{
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       nil,
		Ethash:              nil,
	}

	// configure the genesis block
	genBalance := big.NewInt(100000000000000000)
	genKey, _ := coreth.NewKey(rand.Reader)

	config.Genesis = &core.Genesis{
		Config:     chainConfig,
		Nonce:      0,
		Number:     0,
		ExtraData:  hexutil.MustDecode("0x00"),
		GasLimit:   100000000,
		Difficulty: big.NewInt(0),
		Alloc:      core.GenesisAlloc{genKey.Address: {Balance: genBalance}},
	}

	// grab the control of block generation and disable auto uncle
	config.Miner.ManualMining = true
	config.Miner.DisableUncle = true

	chain := coreth.NewETHChain(&config, nil, nil, nil)
	buff := new(bytes.Buffer)
	blk := chain.GetGenesisBlock()
	err := blk.EncodeRLPEth(buff)
	buff.WriteString("somesuffix")
	checkError(err)
	var blk2 *types.Block
	blk2 = new(types.Block)
	fmt.Println(buff.Len())
	fmt.Println(common.ToHex(buff.Bytes()))

	err = rlp.Decode(buff, blk2)
	fmt.Println(buff.Len())
	checkError(err)
	buff.Reset()
	err = blk2.EncodeRLPEth(buff)
	checkError(err)
	fmt.Println(buff.Len())
	fmt.Println(common.ToHex(buff.Bytes()))

	err = rlp.Decode(buff, blk2)
	fmt.Println(buff.Len())
	checkError(err)
	buff.Reset()
	err = blk2.EncodeRLP(buff)
	checkError(err)
	buff.WriteString("somesuffix")
	fmt.Println(buff.Len())
	fmt.Println(common.ToHex(buff.Bytes()))

	err = rlp.Decode(buff, blk2)
	fmt.Println(buff.Len())
	checkError(err)
	buff.Reset()
	err = blk2.EncodeRLP(buff)
	checkError(err)
	fmt.Println(buff.Len())
	fmt.Println(common.ToHex(buff.Bytes()))

	err = rlp.Decode(buff, blk2)
	fmt.Println(buff.Len())
	checkError(err)
	buff.Reset()
	extra, err := rlp.EncodeToBytes("test extra data")
	err = blk2.EncodeRLPTest(buff, 0xffffffff, extra)
	checkError(err)
	buff.WriteString("somesuffix")
	fmt.Println(buff.Len())
	fmt.Println(common.ToHex(buff.Bytes()))

	err = rlp.Decode(buff, blk2)
	fmt.Println(buff.Len(), (string)(blk2.RawExtraData()))
	checkError(err)
	buff.Reset()
	type NestedData struct {
		A uint16
		B uint16
		S string
	}
	type MyData struct {
		X     uint32
		Y     uint32
		Msg   string
		Inner NestedData
	}
	extra, err = rlp.EncodeToBytes(MyData{
		X: 4200, Y: 4300, Msg: "hello", Inner: NestedData{A: 1, B: 2, S: "world"},
	})
	checkError(err)
	err = blk2.EncodeRLPTest(buff, 0xfffffffe, extra)
	checkError(err)
	decoded1 := new(string)
	err = rlp.DecodeBytes(blk2.RawExtraData(), decoded1)
	checkError(err)
	fmt.Println(buff.Len(), decoded1)
	fmt.Println(common.ToHex(buff.Bytes()))

	err = rlp.Decode(buff, blk2)
	checkError(err)
	decoded2 := new(MyData)
	err = rlp.DecodeBytes(blk2.RawExtraData(), decoded2)
	checkError(err)
	fmt.Println(buff.Len(), decoded2)
	buff.Reset()
}
