package mxc_test

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/mxc"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
)

var (
	testL2RollupAddress = common.HexToAddress("0x79fcdef22feed20eddacbb2587640e45491b757f")
	testKey, _          = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr            = crypto.PubkeyToAddress(testKey.PublicKey)

	genesis    *core.Genesis
	txs        []*types.Transaction
	testEngine *mxc.Mxc
)

func init() {
	config := params.TestChainConfig
	config.GrayGlacierBlock = nil
	config.ArrowGlacierBlock = nil
	config.Ethash = nil
	config.Mxc = true
	testEngine = mxc.New()

	genesis = &core.Genesis{
		Config:     config,
		Alloc:      core.GenesisAlloc{testAddr: {Balance: big.NewInt(2e15)}},
		ExtraData:  []byte("test genesis"),
		Timestamp:  9000,
		Difficulty: common.Big0,
		BaseFee:    big.NewInt(params.InitialBaseFee),
	}

	txs = []*types.Transaction{
		types.MustSignNewTx(testKey, types.LatestSigner(genesis.Config), &types.LegacyTx{
			Nonce:    0,
			Value:    big.NewInt(12),
			GasPrice: big.NewInt(params.InitialBaseFee),
			Gas:      params.TxGas,
			To:       &common.Address{2},
		}),
		types.MustSignNewTx(testKey, types.LatestSigner(genesis.Config), &types.LegacyTx{
			Nonce:    1,
			Value:    big.NewInt(8),
			GasPrice: big.NewInt(params.InitialBaseFee),
			Gas:      params.TxGas,
			To:       &common.Address{2},
		}),
		// prepareBlockTx
		types.MustSignNewTx(testKey, types.LatestSigner(genesis.Config), &types.LegacyTx{
			Nonce:    2,
			Value:    big.NewInt(8),
			GasPrice: big.NewInt(params.InitialBaseFee),
			Gas:      params.TxGas,
			To:       &testL2RollupAddress,
		}),
	}
}

func newTestBackend(t *testing.T) (*eth.Ethereum, []*types.Block) {
	// Generate test chain.
	blocks := generateTestChain()

	// Create node
	n, err := node.New(&node.Config{})
	if err != nil {
		t.Fatalf("can't create new node: %v", err)
	}

	// Create Ethereum Service
	config := &ethconfig.Config{
		Genesis: genesis,
	}

	ethservice, err := eth.New(n, config)
	if err != nil {
		t.Fatalf("can't create new ethereum service: %v", err)
	}

	// Import the test chain.
	if err := n.Start(); err != nil {
		t.Fatalf("can't start test node: %v", err)
	}

	if _, err := ethservice.BlockChain().InsertChain(blocks[1:]); err != nil {
		t.Fatalf("can't import test blocks: %v", err)
	}

	if _, ok := ethservice.Engine().(*mxc.Mxc); !ok {
		t.Fatalf("not use mxc engine")
	}

	return ethservice, blocks
}

func generateTestChain() []*types.Block {
	db := rawdb.NewMemoryDatabase()
	generate := func(i int, g *core.BlockGen) {
		g.OffsetTime(5)

		g.SetExtra([]byte("test_mxc"))
		g.SetDifficulty(common.Big0)

		for _, tx := range txs {
			g.AddTx(tx)
		}
	}

	gblock := genesis.MustCommit(db)

	blocks, _ := core.GenerateChain(genesis.Config, gblock, testEngine, db, 1, generate)

	blocks = append([]*types.Block{gblock}, blocks...)
	return blocks
}

func TestVerifyHeader(t *testing.T) {
	ethService, blocks := newTestBackend(t)

	for _, b := range blocks {
		err := testEngine.VerifyHeader(ethService.BlockChain(), b.Header(), true)
		assert.NoErrorf(t, err, "VerifyHeader error: %s", err)
	}

	err := testEngine.VerifyHeader(ethService.BlockChain(), &types.Header{
		Number:          common.Big1,
		Time:            uint64(time.Now().Unix()),
		BaseFee:         big.NewInt(params.InitialBaseFee),
		WithdrawalsHash: &types.EmptyWithdrawalsHash,
		UncleHash:       types.EmptyUncleHash,
	}, true)
	assert.ErrorIs(t, err, consensus.ErrUnknownAncestor, "VerifyHeader should throw ErrUnknownAncestor when parentHash is unknown")

	err = testEngine.VerifyHeader(ethService.BlockChain(), &types.Header{
		ParentHash:      blocks[len(blocks)-1].Hash(),
		Number:          common.Big0,
		Time:            uint64(time.Now().Unix()),
		BaseFee:         big.NewInt(params.InitialBaseFee),
		WithdrawalsHash: &types.EmptyWithdrawalsHash,
		UncleHash:       types.EmptyUncleHash,
	}, true)
	assert.ErrorIs(t, err, consensus.ErrInvalidNumber, "VerifyHeader should throw ErrInvalidNumber when the block number is wrong")

	err = testEngine.VerifyHeader(ethService.BlockChain(), &types.Header{
		ParentHash:      blocks[len(blocks)-1].Hash(),
		Number:          new(big.Int).SetInt64(int64(len(blocks))),
		Time:            uint64(time.Now().Unix()),
		Extra:           bytes.Repeat([]byte{1}, int(params.MaximumExtraDataSize+1)),
		BaseFee:         big.NewInt(params.InitialBaseFee),
		WithdrawalsHash: &types.EmptyWithdrawalsHash,
		UncleHash:       types.EmptyUncleHash,
	}, true)
	assert.ErrorContains(t, err, "extra-data too long", "VerifyHeader should throw ErrExtraDataTooLong when the block has too much extra data")

	err = testEngine.VerifyHeader(ethService.BlockChain(), &types.Header{
		ParentHash:      blocks[len(blocks)-1].Hash(),
		Number:          new(big.Int).SetInt64(int64(len(blocks))),
		Time:            uint64(time.Now().Unix()),
		Difficulty:      common.Big1,
		BaseFee:         big.NewInt(params.InitialBaseFee),
		WithdrawalsHash: &types.EmptyWithdrawalsHash,
		UncleHash:       types.EmptyUncleHash,
	}, true)
	assert.ErrorContains(t, err, "invalid difficulty", "VerifyHeader should throw ErrInvalidDifficulty when difficulty is not 0")

	err = testEngine.VerifyHeader(ethService.BlockChain(), &types.Header{
		ParentHash:      blocks[len(blocks)-1].Hash(),
		Number:          new(big.Int).SetInt64(int64(len(blocks))),
		Time:            uint64(time.Now().Unix()),
		GasLimit:        params.MaxGasLimit + 1,
		BaseFee:         big.NewInt(params.InitialBaseFee),
		WithdrawalsHash: &types.EmptyWithdrawalsHash,
		UncleHash:       types.EmptyUncleHash,
	}, true)
	assert.ErrorContains(t, err, "invalid gasLimit", "VerifyHeader should throw ErrInvalidGasLimit when gasLimit is higher than the limit")

	err = testEngine.VerifyHeader(ethService.BlockChain(), &types.Header{
		ParentHash: blocks[len(blocks)-1].Hash(),
		Number:     new(big.Int).SetInt64(int64(len(blocks))),
		Time:       uint64(time.Now().Unix()),
		GasLimit:   params.MaxGasLimit,
		BaseFee:    big.NewInt(params.InitialBaseFee),
		UncleHash:  types.EmptyUncleHash,
	}, true)
	assert.ErrorContains(t, err, "withdrawals hash missing", "VerifyHeader should throw ErrWithdrawalsHashMissing withdrawalshash is nil")

	err = testEngine.VerifyHeader(ethService.BlockChain(), &types.Header{
		ParentHash:      blocks[len(blocks)-1].Hash(),
		Number:          new(big.Int).SetInt64(int64(len(blocks))),
		Time:            uint64(time.Now().Unix()),
		GasLimit:        params.MaxGasLimit,
		BaseFee:         big.NewInt(params.InitialBaseFee),
		WithdrawalsHash: &types.EmptyWithdrawalsHash,
	}, true)
	assert.ErrorContains(t, err, "uncles not empty", "VerifyHeader should throw ErrUnclesNotEmpty if uncles is not the empty hash")
}
