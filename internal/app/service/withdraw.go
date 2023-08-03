package service

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/ordererr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/withdrawerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/interfaces"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"time"
)

type WithdrawalService struct {
	WithdrawalRepository interfaces.WithdrawalRepository
}

func NewWithdrawalService(withdrawalRepository interfaces.WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{
		WithdrawalRepository: withdrawalRepository,
	}
}

func (w *WithdrawalService) Withdraw(withdraw models.Withdraw, userID int) error {
	isValid := validateOrderID(withdraw.OrderID)
	if !isValid {
		logger.Log.Errorf("%s has invalid format", withdraw.OrderID)
		return ordererr.ErrWrongOrderID
	}

	withdrawalEntity := entity.Withdrawal{
		UserID:      userID,
		OrderID:     withdraw.OrderID,
		Sum:         withdraw.Sum,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	return w.WithdrawalRepository.Withdraw(withdrawalEntity, userID)
}

func (w *WithdrawalService) GetWithdrawals(userID int) ([]models.Withdraw, error) {
	withdrawals, err := w.WithdrawalRepository.GetWithdrawals(userID)
	if err != nil {
		return nil, err
	}
	if len(withdrawals) == 0 {
		return nil, withdrawerr.ErrNoWithdrawals
	}

	return withdrawals, nil
}
