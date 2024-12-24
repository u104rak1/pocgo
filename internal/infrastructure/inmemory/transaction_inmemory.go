package inmemory

import (
	"context"
	"sort"
	"sync"

	transactionDomain "github.com/u104rak1/pocgo/internal/domain/transaction"
)

type transactionInMemoryRepository struct {
	mu           sync.RWMutex
	transactions []*transactionDomain.Transaction
}

func NewTransactionInMemoryRepository() transactionDomain.ITransactionRepository {
	return &transactionInMemoryRepository{
		transactions: []*transactionDomain.Transaction{},
	}
}

func (r *transactionInMemoryRepository) Save(ctx context.Context, transaction *transactionDomain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.transactions = append(r.transactions, transaction)
	return nil
}

func (r *transactionInMemoryRepository) ListWithTotalByAccountID(ctx context.Context, params transactionDomain.ListTransactionsParams) (transactions []*transactionDomain.Transaction, total int, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filteredTransactions []*transactionDomain.Transaction
	for _, t := range r.transactions {
		if t.AccountIDString() == params.AccountID.String() {
			if params.From != nil && t.TransactionAt().Before(*params.From) {
				continue
			}
			if params.To != nil && t.TransactionAt().After(*params.To) {
				continue
			}
			if len(params.OperationTypes) > 0 {
				match := false
				for _, opType := range params.OperationTypes {
					if t.OperationType() == opType {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			filteredTransactions = append(filteredTransactions, t)
		}
	}

	total = len(filteredTransactions)

	// Sort and paginate
	if params.Sort != nil && *params.Sort == "ASC" {
		sort.Slice(filteredTransactions, func(i, j int) bool {
			return filteredTransactions[i].TransactionAt().Before(filteredTransactions[j].TransactionAt())
		})
	} else {
		sort.Slice(filteredTransactions, func(i, j int) bool {
			return filteredTransactions[i].TransactionAt().After(filteredTransactions[j].TransactionAt())
		})
	}

	if params.Limit != nil && params.Page != nil {
		start := (*params.Page - 1) * *params.Limit
		end := start + *params.Limit
		if start < total {
			if end > total {
				end = total
			}
			transactions = filteredTransactions[start:end]
		}
	} else {
		transactions = filteredTransactions
	}

	return transactions, total, nil
}
