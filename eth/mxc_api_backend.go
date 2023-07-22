package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// MxcAPIBackend handles l2 node related RPC calls.
type MxcAPIBackend struct {
	eth *Ethereum
}

// NewMxcAPIBackend creates a new MxcAPIBackend instance.
func NewMxcAPIBackend(eth *Ethereum) *MxcAPIBackend {
	return &MxcAPIBackend{
		eth: eth,
	}
}

// HeadL1Origin returns the latest L2 block's corresponding L1 origin.
func (s *MxcAPIBackend) HeadL1Origin() (*rawdb.L1Origin, error) {
	blockID, err := rawdb.ReadHeadL1Origin(s.eth.ChainDb())
	if err != nil {
		return nil, err
	}

	if blockID == nil {
		return nil, ethereum.NotFound
	}

	l1Origin, err := rawdb.ReadL1Origin(s.eth.ChainDb(), blockID)
	if err != nil {
		return nil, err
	}

	if l1Origin == nil {
		return nil, ethereum.NotFound
	}

	return l1Origin, nil
}

// L1OriginByID returns the L2 block's corresponding L1 origin.
func (s *MxcAPIBackend) L1OriginByID(blockID *math.HexOrDecimal256) (*rawdb.L1Origin, error) {
	l1Origin, err := rawdb.ReadL1Origin(s.eth.ChainDb(), (*big.Int)(blockID))
	if err != nil {
		return nil, err
	}

	if l1Origin == nil {
		return nil, ethereum.NotFound
	}

	return l1Origin, nil
}

// TxPoolContent retrieves the transaction pool content with the given upper limits.
func (s *MxcAPIBackend) TxPoolContent(
	maxTransactionsPerBlock uint64,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	minTxGasLimit uint64,
	locals []string,
) ([]types.Transactions, error) {
	pending := s.eth.TxPool().Pending(false)

	log.Debug(
		"Fetching L2 pending transactions finished",
		"length", core.PoolContent(pending).Len(),
		"maxTransactionsPerBlock", maxTransactionsPerBlock,
		"blockMaxGasLimit", blockMaxGasLimit,
		"maxBytesPerTxList", maxBytesPerTxList,
		"minTxGasLimit", minTxGasLimit,
		"locals", locals,
	)

	contentSplitter, err := core.NewPoolContentSplitter(
		s.eth.BlockChain().Config().ChainID,
		maxTransactionsPerBlock,
		blockMaxGasLimit,
		maxBytesPerTxList,
		minTxGasLimit,
		locals,
	)
	if err != nil {
		return nil, err
	}

	var (
		txsCount = 0
		txLists  []types.Transactions
	)
	for _, splittedTxs := range contentSplitter.Split(filterTxs(pending, s.eth.blockchain.CurrentHeader().BaseFee)) {
		if txsCount+splittedTxs.Len() < int(maxTransactionsPerBlock) {
			txLists = append(txLists, splittedTxs)
			txsCount += splittedTxs.Len()
			continue
		}

		txLists = append(txLists, splittedTxs[0:(int(maxTransactionsPerBlock)-txsCount)])
		break
	}

	return txLists, nil
}

func filterTxs(pendings map[common.Address]types.Transactions, baseFee *big.Int) map[common.Address]types.Transactions {
	executableTxs := make(map[common.Address]types.Transactions)
	gasPriceLowerLimit := big.NewInt(0).Div(big.NewInt(0).Mul(baseFee, big.NewInt(95)), big.NewInt(100))

	for addr, txs := range pendings {
		pendingTxs := make(types.Transactions, 0)
		for _, tx := range txs {
			// Check baseFee, should not be zero
			if tx.GasFeeCap().Uint64() == 0 || tx.GasPrice().Cmp(gasPriceLowerLimit) < 0 {
				log.Debug("Ignore max fee per gas less than block base fee",
					"gas price", baseFee.Uint64(),
					"tx gas price", tx.GasPrice().Uint64(),
					"lower limit", gasPriceLowerLimit.Uint64(),
				)
				break
			}

			pendingTxs = append(pendingTxs, tx)
		}

		if len(pendingTxs) > 0 {
			executableTxs[addr] = pendingTxs
		}
	}

	return executableTxs
}
