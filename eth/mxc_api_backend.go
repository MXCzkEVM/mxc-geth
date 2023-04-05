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

// MXCAPIBackend handles l2 node related RPC calls.
type MXCAPIBackend struct {
	eth *Ethereum
}

// NewMXCAPIBackend creates a new MXCAPIBackend instance.
func NewMXCAPIBackend(eth *Ethereum) *MXCAPIBackend {
	return &MXCAPIBackend{
		eth: eth,
	}
}

// HeadL1Origin returns the latest L2 block's corresponding L1 origin.
func (s *MXCAPIBackend) HeadL1Origin() (*rawdb.L1Origin, error) {
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
func (s *MXCAPIBackend) L1OriginByID(blockID *math.HexOrDecimal256) (*rawdb.L1Origin, error) {
	l1Origin, err := rawdb.ReadL1Origin(s.eth.ChainDb(), (*big.Int)(blockID))
	if err != nil {
		return nil, err
	}

	if l1Origin == nil {
		return nil, ethereum.NotFound
	}

	return l1Origin, nil
}

// GetThrowawayTransactionReceipts returns the throwaway block's receipts
// without checking whether the block is in the canonical chain.
func (s *MXCAPIBackend) GetThrowawayTransactionReceipts(hash common.Hash) (types.Receipts, error) {
	receipts := s.eth.blockchain.GetReceiptsByHash(hash)
	if receipts == nil {
		return nil, ethereum.NotFound
	}

	return receipts, nil
}

// TxPoolContent retrieves the transaction pool content with the given upper limits.
func (s *MXCAPIBackend) TxPoolContent(
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
	for _, splittedTxs := range contentSplitter.Split(pending) {
		if txsCount+splittedTxs.Len() < int(maxTransactionsPerBlock) {
			for i, tx := range splittedTxs {
				if tx.GasPrice().Cmp(s.eth.gasPrice) < 0 {
					splittedTxs = append(splittedTxs[:i], splittedTxs[i+1:]...)
					log.Info("Dropping tx with low gas price", "hash", tx.Hash(), "gasPrice", tx.GasPrice(), "minGasPrice", s.eth.gasPrice)
				}
			}
			txLists = append(txLists, splittedTxs)
			txsCount += splittedTxs.Len()
			continue
		}

		txLists = append(txLists, splittedTxs[0:(int(maxTransactionsPerBlock)-txsCount)])
		break
	}

	return txLists, nil
}
