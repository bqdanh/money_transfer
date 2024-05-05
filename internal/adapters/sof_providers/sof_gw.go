package sof_providers

import (
	"context"
	"fmt"

	"github.com/bqdanh/money_transfer/internal/entities/account"
	"github.com/bqdanh/money_transfer/internal/entities/transaction"
)

type SofProviderRegistrar interface {
	GetSourceOfFundType() account.SourceOfFundType
	GetSourceOfFundCode() account.SourceOfFundCode

	MakeDepositTransaction(ctx context.Context, trans transaction.Transaction) (transaction.Data, error)
}

type SofGateway struct {
	sofsRegistry map[account.SourceOfFundType]map[account.SourceOfFundCode]SofProviderRegistrar
}

func NewSofGateway(sofProviders ...SofProviderRegistrar) (*SofGateway, error) {
	sofsRegistry := make(map[account.SourceOfFundType]map[account.SourceOfFundCode]SofProviderRegistrar)
	for _, sof := range sofProviders {
		if sofsRegistry[sof.GetSourceOfFundType()] == nil {
			sofsRegistry[sof.GetSourceOfFundType()] = make(map[account.SourceOfFundCode]SofProviderRegistrar)
		}
		if sofsRegistry[sof.GetSourceOfFundType()][sof.GetSourceOfFundCode()] != nil {
			return nil, fmt.Errorf("source of fund type(%s) code(%s) is duplicated", sof.GetSourceOfFundType(), sof.GetSourceOfFundCode())
		}
		sofsRegistry[sof.GetSourceOfFundType()][sof.GetSourceOfFundCode()] = sof
	}
	return &SofGateway{
		sofsRegistry: sofsRegistry,
	}, nil
}

func (s *SofGateway) MakeDepositTransaction(ctx context.Context, trans transaction.Transaction) (transaction.Data, error) {
	sofsTypeRegistry, ok := s.sofsRegistry[trans.Account.SourceOfFundData.GetSourceOfFundType()]
	if !ok {
		return transaction.Data{}, fmt.Errorf("source of fund type(%s) is not supported", trans.Account.SourceOfFundData.GetSourceOfFundType())
	}
	sofProvider, ok := sofsTypeRegistry[trans.Account.SourceOfFundData.GetSourceOfFundCode()]
	if !ok {
		return transaction.Data{}, fmt.Errorf("source of fund code(%s) is not supported", trans.Account.SourceOfFundData.GetSourceOfFundCode())
	}
	transData, err := sofProvider.MakeDepositTransaction(ctx, trans)
	if err != nil {
		return transaction.Data{}, fmt.Errorf("sof provider failed to make deposit transaction: %w", err)
	}
	return transData, nil
}
