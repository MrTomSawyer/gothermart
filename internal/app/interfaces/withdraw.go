package interfaces

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
)

type WithdrawalRepository interface {
	Withdraw(withdraw entity.Withdrawal, userID int) error
	GetWithdrawals(userID int) ([]models.Withdraw, error)
}

type WithdrawalService interface {
	Withdraw(withdraw models.Withdraw, userID int) error
	GetWithdrawals(userID int) ([]models.Withdraw, error)
}
