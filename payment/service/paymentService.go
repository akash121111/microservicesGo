package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Wallet struct {
	UserID  uuid.UUID
	Balance float64
}

type PaymentService struct {
	wallets map[uuid.UUID]*Wallet
	mu      sync.Mutex
}

func NewPaymentService() *PaymentService {

	return &PaymentService{
		wallets: make(map[uuid.UUID]*Wallet),
	}
}

func (s *PaymentService) ProcessPayment(
	userID uuid.UUID,
	amount float64,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	wallet, exists := s.wallets[userID]

	if !exists {
		return fmt.Errorf("wallet not found")
	}

	if wallet.Balance < amount {
		return fmt.Errorf("insufficient wallet balance")
	}

	wallet.Balance -= amount

	return nil
}
func (s *PaymentService) CreateWallet(
	userID uuid.UUID,
	balance float64,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wallets[userID] = &Wallet{
		UserID:  userID,
		Balance: balance,
	}
}
